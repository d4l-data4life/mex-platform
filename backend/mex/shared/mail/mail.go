package mail

import "context"

type Mailer interface {
	SendMails(ctx context.Context, order *MailOrder, data Data) ([]string, error)
}
