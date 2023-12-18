package flowmailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/d4l-data4life/mex/mex/shared/mail"
	"github.com/d4l-data4life/mex/mex/shared/utils"
)

// Simple Flowmailer:
//   - uses the Default flow
//   - does not do any template interpolation in Flowmailer
//     (texts are interpolated via Golang templates before sending)
type flowmailerSimple struct {
	originOAuth string
	originAPI   string
	clientID    string
	accountID   string

	clientSecret string

	client      http.Client
	accessToken string

	noreplyEmail string
}

type Header struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SubmitMessage struct {
	MessageType string `json:"messageType,omitempty"`

	HeaderFromName    string `json:"headerFromName,omitempty"`
	HeaderFromAddress string `json:"headerFromAddress,omitempty"`

	Headers []Header `json:"headers,omitempty"` // will include the To header

	RecipientAddress string `json:"recipientAddress,omitempty"`
	SenderAddress    string `json:"senderAddress,omitempty"`

	Subject string `json:"subject,omitempty"`

	HTML string `json:"html,omitempty"`
	Text string `json:"text,omitempty"`
}

type Params struct {
	OriginOAuth         string
	OriginAPI           string
	ClientID            string
	ClientSecret        string
	AccountID           string
	NoReplyEmailAddress string
	Timeout             time.Duration
}

func NewSimpleFlowmailer(params Params) (mail.Mailer, error) {
	if params.NoReplyEmailAddress == "" {
		return nil, fmt.Errorf("no-reply email address is not specified")
	}

	return &flowmailerSimple{
		originOAuth:  params.OriginOAuth,
		originAPI:    params.OriginAPI,
		clientID:     params.ClientID,
		clientSecret: params.ClientSecret,
		accountID:    params.AccountID,
		noreplyEmail: params.NoReplyEmailAddress,

		client: http.Client{
			Timeout: params.Timeout,
		},
	}, nil
}

// This function takes a mail order object and sends out emails for each recipient (to, cc, bcc).
// Departing from the Go standard, this function *always* returns a string slice, even in the error case.
// This slice contains the message IDs of successfully submitted emails, since we cannot undo the submission
// of earlier mails in case of an error occurring submitting later mails (of the same mail order).
func (fm *flowmailerSimple) SendMails(ctx context.Context, order *mail.MailOrder, data mail.Data) ([]string, error) {
	messageIDs := []string{}

	msg, err := mailOrderToFlowmailerMessage(order, data, fm.noreplyEmail)
	if err != nil {
		return messageIDs, err
	}

	// We get a fresh access token so that it holds for the entire sequence of email submissions below.
	err = fm.authenticate(ctx)
	if err != nil {
		return messageIDs, err
	}

	for _, recipient := range order.RecipientsTo {
		msg.RecipientAddress = recipient.Email
		messageID, err := fm.submitMessage(ctx, msg)
		if err != nil {
			return messageIDs, err
		}
		messageIDs = append(messageIDs, messageID)
	}

	for _, recipient := range order.RecipientsCc {
		msg.RecipientAddress = recipient.Email
		messageID, err := fm.submitMessage(ctx, msg)
		if err != nil {
			return messageIDs, err
		}
		messageIDs = append(messageIDs, messageID)
	}

	for _, recipient := range order.RecipientsBcc {
		msg.RecipientAddress = recipient.Email
		messageID, err := fm.submitMessage(ctx, msg)
		if err != nil {
			return messageIDs, err
		}
		messageIDs = append(messageIDs, messageID)
	}

	return messageIDs, nil
}

func mailOrderToFlowmailerMessage(order *mail.MailOrder, data mail.Data, noreplyEmail string) (*SubmitMessage, error) {
	subject, err := mail.InterpolateGoPlain(order.Subject, data)
	if err != nil {
		return nil, err
	}

	textBody, err := mail.InterpolateGoPlain(order.TextBody, data)
	if err != nil {
		return nil, err
	}

	htmlBody, err := mail.InterpolateGoHTML(order.HtmlBody, data)
	if err != nil {
		return nil, err
	}

	// If both SenderAddress and HeaderFromAddress are given, HeaderFromAddress must be from an
	// authenticated domain while SenderAddress may be anything.
	// If only SenderAddress is given, it must be from an authenticated domain.
	msg := SubmitMessage{
		MessageType: "EMAIL",

		Subject: subject,
		Text:    textBody,
		HTML:    htmlBody,

		SenderAddress: noreplyEmail,

		// HeaderFromAddress not used
		HeaderFromName: order.Sender.Name,

		Headers: []Header{{
			Name:  "To",
			Value: concatHeaderNames(order.RecipientsTo),
		}, {
			Name:  "Reply-To",
			Value: order.Sender.Email,
		}, {
			Name:  "Order-Id",
			Value: order.OrderId,
		}},
	}

	if len(order.RecipientsCc) > 0 {
		msg.Headers = append(msg.Headers, Header{
			Name:  "Cc",
			Value: concatHeaderNames(order.RecipientsCc),
		})
	}
	return &msg, nil
}

func (fm *flowmailerSimple) authenticate(ctx context.Context) error {
	fm.accessToken = ""

	creds := url.Values{}
	creds.Set("client_id", fm.clientID)
	creds.Set("client_secret", fm.clientSecret)
	creds.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/oauth/token", fm.originOAuth), strings.NewReader(creds.Encode()))
	if err != nil {
		return fmt.Errorf("error creating new HTTP request: %w", err)
	}

	req.Close = true
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := fm.client.Do(req)
	if err != nil {
		return fmt.Errorf("error during HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unsuccessful OAuth request: %d, %w", resp.StatusCode, err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	var token struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		Scope       string `json:"scope"`
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return fmt.Errorf("body parsing failed: %w", err)
	}

	fm.accessToken = token.AccessToken
	return nil
}

func (fm *flowmailerSimple) submitMessage(ctx context.Context, msg *SubmitMessage) (string, error) {
	rawBody := bytes.Buffer{}
	encoder := json.NewEncoder(&rawBody)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(*msg)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s/messages/submit", fm.originAPI, fm.accountID), &rawBody)
	if err != nil {
		return "", fmt.Errorf("error creating new HTTP request: %w", err)
	}

	req.Close = true
	req.Header.Add("Content-Type", "application/vnd.flowmailer.v1.12+json;charset=UTF-8")
	req.Header.Add("Authorization", "Bearer "+fm.accessToken)

	resp, err := fm.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error during HTTP request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return "", fmt.Errorf("failed to authenticate to Flowmailer API: %d: %s", resp.StatusCode, string(respBody))
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("unsuccessful submission request: %d: %s", resp.StatusCode, string(respBody))
	}

	messageID := strings.Split(resp.Header.Get("Location"), "messages/")[1]
	return messageID, nil
}

func concatHeaderNames(contacts []*mail.StaticContact) string {
	if contacts == nil {
		return "<nil>"
	}
	return strings.Join(utils.Map(contacts, func(c *mail.StaticContact) string { return c.HeaderString() }), ", ")
}
