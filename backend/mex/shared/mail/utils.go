package mail

import (
	"fmt"
)

type FormData map[string]any

type Data = FormData

type Item map[string][]string

func (mt *MailTemplate) DetermineSenderNameAndEmail(formData FormData) (*StaticContact, error) {
	if mt == nil {
		return nil, fmt.Errorf("mail template is nil")
	}

	if mt.Sender == nil {
		return nil, fmt.Errorf("mail template sender is nil")
	}

	name := ""
	email := ""
	ok := false

	switch t := mt.Sender.ContactType.(type) {
	case *Sender_Static:
		email = t.Static.Email
		if email == "" {
			return nil, fmt.Errorf("no email in static contact")
		}
		name = t.Static.Name
		if name == "" {
			return nil, fmt.Errorf("no name in static contact")
		}
	case *Sender_FormData:
		name, ok = formData.GetFirstValue(t.FormData.NameField)
		if !ok {
			return nil, fmt.Errorf("sender: form data: no value for field %q", t.FormData.NameField)
		}
		email, ok = formData.GetFirstValue(t.FormData.EmailField)
		if !ok {
			return nil, fmt.Errorf("sender: form data: no value for field %q", t.FormData.EmailField)
		}
	default:
		return nil, fmt.Errorf("unsupported sender contact type: %T", t)
	}

	return &StaticContact{
		Name:  name,
		Email: email,
	}, nil
}

type RecipientsContacts map[RecipientType][]*StaticContact

func (mt *MailTemplate) DetermineRecipientsNameAndEmail(formData FormData, recipientItem Item, contextItem Item) (RecipientsContacts, error) {
	if mt == nil {
		return nil, fmt.Errorf("mail template is nil")
	}

	if mt.Recipients == nil {
		return nil, fmt.Errorf("mail template recipients is nil")
	}

	ret := map[RecipientType][]*StaticContact{
		RecipientType_TO:  {},
		RecipientType_CC:  {},
		RecipientType_BCC: {},
	}

	for _, recipient := range mt.Recipients {
		name := ""
		email := ""
		ok := false

		switch t := recipient.ContactType.(type) {
		case *Recipient_Static:
			email = t.Static.Email
			if email == "" {
				return nil, fmt.Errorf("no email in static contact")
			}
			name = t.Static.Name
			if name == "" {
				return nil, fmt.Errorf("no name in static contact")
			}
		case *Recipient_FormData:
			name, ok = formData.GetFirstValue(t.FormData.NameField)
			if !ok {
				return nil, fmt.Errorf("recipient: form data: no value for name field %q", t.FormData.NameField)
			}
			email, ok = formData.GetFirstValue(t.FormData.EmailField)
			if !ok {
				return nil, fmt.Errorf("recipient: form data: no value for email field %q", t.FormData.EmailField)
			}

		case *Recipient_ContactItem:
			name, ok = recipientItem.GetFirstValue(t.ContactItem.NameField)
			if !ok {
				return nil, fmt.Errorf("recipient: contact item: no value for name field %q", t.ContactItem.NameField)
			}
			email, ok = recipientItem.GetFirstValue(t.ContactItem.EmailField)
			if !ok {
				return nil, fmt.Errorf("recipient: contact item: no value for email field %q", t.ContactItem.EmailField)
			}

		case *Recipient_ContextItem:
			name, ok = contextItem.GetFirstValue(t.ContextItem.NameField)
			if !ok {
				return nil, fmt.Errorf("recipient: context item: no value for name field %q", t.ContextItem.NameField)
			}
			email, ok = contextItem.GetFirstValue(t.ContextItem.EmailField)
			if !ok {
				return nil, fmt.Errorf("recipient: context item: no value for email field %q", t.ContextItem.EmailField)
			}

		default:
			return nil, fmt.Errorf("unsupported recipient contact type: %T", t)
		}

		ret[recipient.Type] = append(ret[recipient.Type], &StaticContact{
			Name:  name,
			Email: email,
		})
	}

	return ret, nil
}

func (item Item) GetFirstValue(fieldName string) (string, bool) {
	if fieldName == "" {
		return "", false
	}

	values := item[fieldName]
	if values == nil {
		return "", false
	}

	if len(values) == 0 {
		return "", false
	}

	return (values)[0], (values)[0] != ""
}

// FormData fields may be strings or string slices.
// This function treats strings as one-element slices.
func (formData FormData) GetFirstValue(fieldName string) (string, bool) {
	if fieldName == "" {
		return "", false
	}

	values := formData[fieldName]
	if values == nil {
		return "", false
	}

	switch t := values.(type) {
	case string:
		return t, t != ""
	case []string:
		if len(t) == 0 {
			return "", false
		}
		return t[0], t[0] != ""
	default:
		return "", false
	}
}

// Turn a StaticContact to a string formatted as "Henry Jones" <hjones@atlantis.com>
// if both name and email are given.
func (c *StaticContact) HeaderString() string {
	if c.Name != "" && c.Email != "" {
		return fmt.Sprintf(`"%s" <%s>`, c.Name, c.Email)
	}

	if c.Name != "" {
		return c.Name
	}

	if c.Email != "" {
		return c.Email
	}

	return "n/a"
}
