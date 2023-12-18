package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/encoding/protojson"

	E "github.com/d4l-data4life/mex/mex/shared/errstat"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/mail"
	"github.com/d4l-data4life/mex/mex/shared/uuid"

	dbNotify "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/notify/db"
	pbNotify "github.com/d4l-data4life/mex/mex/services/metadata/endpoints/notify/pb"
)

type Service struct {
	Log L.Logger

	DB    *pgxpool.Pool
	Redis *redis.Client

	Mailer mail.Mailer

	ConfigServiceOrigin string

	pbNotify.UnimplementedNotifyServer
}

func (svc *Service) SendNotification(ctx context.Context, request *pbNotify.SendNotificationRequest) (*pbNotify.SendNotificationResponse, error) {
	if request == nil {
		return nil, E.MakeGRPCStatus(codes.InvalidArgument, "item notification is nil").Err()
	}

	if request.TemplateInfo == nil {
		return nil, E.MakeGRPCStatus(codes.InvalidArgument, "item notification template info is nil").Err()
	}

	if request.TemplateInfo.TemplateName == "" {
		return nil, E.MakeGRPCStatus(codes.InvalidArgument, "no template name given").Err()
	}

	// Gather all information

	formDataMap := mail.FormData{}
	err := json.Unmarshal([]byte(request.FormData), &formDataMap)
	if err != nil {
		return nil, err
	}

	contextItemMap, err := svc.getItem(ctx, request.TemplateInfo.ContextItemId)
	if err != nil {
		return nil, err
	}

	recipientItemMap, err := svc.getItem(ctx, request.TemplateInfo.RecipientItemId)
	if err != nil {
		return nil, err
	}

	mailTemplate, err := svc.retrieveTemplate(request.TemplateInfo.TemplateName)
	if err != nil {
		return nil, err
	}

	// Determine senders's name and email address.
	sender, err := mailTemplate.DetermineSenderNameAndEmail(formDataMap)
	if err != nil {
		return nil, err
	}

	// Determine recipients' names and email addresses from item.
	recipients, err := mailTemplate.DetermineRecipientsNameAndEmail(formDataMap, recipientItemMap, contextItemMap)
	if err != nil {
		return nil, err
	}

	orderID := uuid.MustNewV4()

	order := mail.MailOrder{
		OrderId:        orderID,
		TemplateEngine: mailTemplate.TemplateEngine,

		Subject:  mailTemplate.Subject,
		TextBody: mailTemplate.TextBody,
		HtmlBody: mailTemplate.HtmlBody,

		Sender: sender,

		RecipientsTo:  recipients[mail.RecipientType_TO],
		RecipientsCc:  recipients[mail.RecipientType_CC],
		RecipientsBcc: recipients[mail.RecipientType_BCC],
	}

	messageIDs, err := svc.Mailer.SendMails(ctx, &order, map[string]any{
		"context": contextItemMap,
		"form":    formDataMap,
		"orderId": orderID,
	})
	if err != nil {
		if len(messageIDs) == 0 {
			return nil, fmt.Errorf("no mails could be sent: %s", err.Error())
		}
		svc.Log.Warn(ctx, L.Messagef("not mails could be sent: %s", err.Error()))
	}

	response := &pbNotify.SendNotificationResponse{
		MessageIds: messageIDs,
		OrderId:    orderID,
	}

	svc.Log.BIEvent(ctx, L.BIActivity("mail-notify"), L.BIData(response))
	return response, nil
}

func (svc *Service) getItem(ctx context.Context, itemID string) (mail.Item, error) {
	queries := dbNotify.New(svc.DB)

	fieldValues, err := queries.DbGetItemValues(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving item: %s", err.Error())
	}

	if len(fieldValues) == 0 {
		return nil, fmt.Errorf("could not find item or item value for item ID %q", itemID)
	}

	item := mail.Item{}
	for _, v := range fieldValues {
		f := item[v.FieldName]
		if f == nil {
			f = []string{}
		}

		f = append(f, v.FieldValue)
		item[v.FieldName] = f
	}

	return item, nil
}

func (svc *Service) retrieveTemplate(templateName string) (*mail.MailTemplate, error) {
	client := &http.Client{Timeout: time.Second}
	resp, err := client.Get(fmt.Sprintf("%s/api/v0/config/files/mail_templates/%s", svc.ConfigServiceOrigin, templateName))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not fetch template %q: status code %d", templateName, resp.StatusCode)
	}

	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	mailTemplate := mail.MailTemplate{}
	err = protojson.Unmarshal(buf, &mailTemplate)
	if err != nil {
		return nil, err
	}

	return &mailTemplate, nil
}
