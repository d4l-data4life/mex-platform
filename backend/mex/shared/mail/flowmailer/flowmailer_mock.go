package flowmailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/d4l-data4life/mex/mex/shared/mail"
)

const (
	hashPrefix     = "mockmailer"
	mailExpiration = 10 * time.Minute
)

type mockmailer struct {
	redis *redis.Client
}

func NewMockMailer(redis *redis.Client) mail.Mailer {
	return &mockmailer{redis: redis}
}

func (mm *mockmailer) SendMails(ctx context.Context, order *mail.MailOrder, data mail.Data) ([]string, error) {
	messageIDs := []string{}

	msg, err := mailOrderToFlowmailerMessage(order, data, "noreply@data4life.care")
	if err != nil {
		return messageIDs, err
	}

	for _, recipient := range order.RecipientsTo {
		msg.RecipientAddress = recipient.Email
		messageID, err := mm.submitMessage(ctx, msg, order.OrderId, "to")
		if err != nil {
			return messageIDs, err
		}
		messageIDs = append(messageIDs, messageID)
	}

	for _, recipient := range order.RecipientsCc {
		msg.RecipientAddress = recipient.Email
		messageID, err := mm.submitMessage(ctx, msg, order.OrderId, "cc")
		if err != nil {
			return messageIDs, err
		}
		messageIDs = append(messageIDs, messageID)
	}

	for _, recipient := range order.RecipientsBcc {
		msg.RecipientAddress = recipient.Email
		messageID, err := mm.submitMessage(ctx, msg, order.OrderId, "bcc")
		if err != nil {
			return messageIDs, err
		}
		messageIDs = append(messageIDs, messageID)
	}

	return messageIDs, nil
}

func (mm *mockmailer) submitMessage(ctx context.Context, msg *SubmitMessage, orderID string, mailType string) (string, error) {
	rawBody := bytes.Buffer{}
	encoder := json.NewEncoder(&rawBody)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(*msg)
	if err != nil {
		return "", err
	}

	key := fmt.Sprintf("%s:%s:%s:%s", hashPrefix, orderID, mailType, msg.RecipientAddress)
	err = mm.redis.HSet(ctx, key,
		"message_id", key,
		"message", rawBody.String(),
		"mail_type", mailType,
	).Err()
	if err != nil {
		return "", err
	}
	_ = mm.redis.Expire(ctx, key, mailExpiration)

	return key, nil
}
