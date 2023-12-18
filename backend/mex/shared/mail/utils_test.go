package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DetermineSenderNameAndEmail(t *testing.T) {
	cases := []struct {
		desc         string
		mailTemplate *MailTemplate
		formData     FormData
		wantError    bool
		out          *StaticContact
	}{
		{
			desc:         "nothing given",
			mailTemplate: nil,
			wantError:    true,
		},
		{
			desc: "mail template broken",
			mailTemplate: &MailTemplate{
				Sender: nil,
			},
			wantError: true,
		},

		/////////////////////////////////////////////////////////////////////
		{
			mailTemplate: &MailTemplate{
				Sender: &Sender{
					ContactType: &Sender_Static{
						Static: &StaticContact{
							Name:  "Henry Jones",
							Email: "h.jones@barnett.edu",
						},
					},
				},
			},
			out: &StaticContact{
				Name:  "Henry Jones",
				Email: "h.jones@barnett.edu",
			},
		},
		{
			desc: "need name",
			mailTemplate: &MailTemplate{
				Sender: &Sender{
					ContactType: &Sender_Static{
						Static: &StaticContact{
							Name:  "",
							Email: "h.jones@barnett.edu",
						},
					},
				},
			},
			wantError: true,
		},
		{
			desc: "need email",
			mailTemplate: &MailTemplate{
				Sender: &Sender{
					ContactType: &Sender_Static{
						Static: &StaticContact{
							Name: "Henry Jones",
						},
					},
				},
			},
			wantError: true,
		},
		/////////////////////////////////////////////////////////////////////
		{
			mailTemplate: &MailTemplate{
				Sender: &Sender{
					ContactType: &Sender_FormData{
						FormData: &FieldNamesContact{
							NameField:  "name",
							EmailField: "email",
						},
					},
				},
			},
			formData: map[string]any{
				"name":  "Henry Jones",
				"email": "h.jones@barnett.edu",
			},
			out: &StaticContact{
				Name:  "Henry Jones",
				Email: "h.jones@barnett.edu",
			},
		},
		{
			desc: "multiple field values",
			mailTemplate: &MailTemplate{
				Sender: &Sender{
					ContactType: &Sender_FormData{
						FormData: &FieldNamesContact{
							NameField:  "name",
							EmailField: "email",
						},
					},
				},
			},
			formData: map[string]any{
				"name":  []string{"Henry Jones", "Guybrush Threepwood"},
				"email": "h.jones@barnett.edu",
			},
			out: &StaticContact{
				Name:  "Henry Jones",
				Email: "h.jones@barnett.edu",
			},
		},
		{
			desc: "name field name not existing",
			mailTemplate: &MailTemplate{
				Sender: &Sender{
					ContactType: &Sender_FormData{
						FormData: &FieldNamesContact{
							NameField: "fooo",
						},
					},
				},
			},
			formData: map[string]any{
				"name":  "Henry Jones",
				"email": "h.jones@barnett.edu",
			},
			wantError: true,
		},
		{
			desc: "email field name not specified",
			mailTemplate: &MailTemplate{
				Sender: &Sender{
					ContactType: &Sender_FormData{
						FormData: &FieldNamesContact{
							NameField: "name",
						},
					},
				},
			},
			formData: map[string]any{
				"name":  "Henry Jones",
				"email": "h.jones@barnett.edu",
			},
			wantError: true,
		},
	}

	for _, test := range cases {
		out, err := test.mailTemplate.DetermineSenderNameAndEmail(test.formData)
		if test.wantError {
			if err == nil {
				t.Error("expected error but got nil")
			}
		} else {
			require.Equal(t, test.out, out, test.desc)
		}
	}
}

func Test_DetermineRecipientsNameAndEmail(t *testing.T) {
	cases := []struct {
		desc          string
		mailTemplate  *MailTemplate
		formData      FormData
		recipientItem Item
		contextItem   Item
		wantError     bool
		out           RecipientsContacts
	}{
		{
			mailTemplate: nil,
			wantError:    true,
		},
		{
			mailTemplate: &MailTemplate{
				Recipients: nil,
			},
			wantError: true,
		},

		/////////////////////////////////////////////////////////////////////
		{
			mailTemplate: &MailTemplate{
				Recipients: []*Recipient{
					{
						Type: RecipientType_TO,
						ContactType: &Recipient_Static{
							Static: &StaticContact{
								Name:  "Henry Jones",
								Email: "h.jones@barnett.edu",
							},
						},
					},
				},
			},
			out: RecipientsContacts{
				RecipientType_TO: {
					{
						Name:  "Henry Jones",
						Email: "h.jones@barnett.edu",
					},
				},
				RecipientType_CC:  {},
				RecipientType_BCC: {},
			},
		},
		/////////////////////////////////////////////////////////////////////
		{
			mailTemplate: &MailTemplate{
				Recipients: []*Recipient{
					{
						Type: RecipientType_TO,
						ContactType: &Recipient_Static{
							Static: &StaticContact{
								Name:  "Henry Jones",
								Email: "h.jones@barnett.edu",
							},
						},
					},
					{
						Type: RecipientType_TO,
						ContactType: &Recipient_Static{
							Static: &StaticContact{
								Name:  "Guybrush Threepwood",
								Email: "guybrush@scummbar.mi",
							},
						},
					},
					{
						Type: RecipientType_CC,
						ContactType: &Recipient_Static{
							Static: &StaticContact{
								Name:  "LeChuck",
								Email: "lechuck@scummbar.mi",
							},
						},
					},
					{
						Type: RecipientType_BCC,
						ContactType: &Recipient_Static{
							Static: &StaticContact{
								Name:  "Zak McKracken",
								Email: "zak@natinq.com",
							},
						},
					},
				},
			},
			out: RecipientsContacts{
				RecipientType_TO: {
					{
						Name:  "Henry Jones",
						Email: "h.jones@barnett.edu",
					},
					{
						Name:  "Guybrush Threepwood",
						Email: "guybrush@scummbar.mi",
					},
				},
				RecipientType_CC: {
					{
						Name:  "LeChuck",
						Email: "lechuck@scummbar.mi",
					},
				},
				RecipientType_BCC: {
					{
						Name:  "Zak McKracken",
						Email: "zak@natinq.com",
					},
				},
			},
		},

		{
			mailTemplate: &MailTemplate{
				Recipients: []*Recipient{
					{
						Type: RecipientType_TO,
						ContactType: &Recipient_Static{
							Static: &StaticContact{
								Email: "h.jones@barnett.edu",
							},
						},
					},
				},
			},
			wantError: true,
		},
		{
			mailTemplate: &MailTemplate{
				Recipients: []*Recipient{
					{
						Type: RecipientType_TO,
						ContactType: &Recipient_Static{
							Static: &StaticContact{
								Name: "h.jones@barnett.edu",
							},
						},
					},
				},
			},
			wantError: true,
		},
		/////////////////////////////////////////////////////////////////////
		{
			mailTemplate: &MailTemplate{
				Recipients: []*Recipient{
					{
						Type: RecipientType_TO,
						ContactType: &Recipient_FormData{
							FormData: &FieldNamesContact{
								NameField:  "name",
								EmailField: "email",
							},
						},
					},
				},
			},
			formData: map[string]any{
				"name":  "Henry Jones",
				"email": "h.jones@barnett.edu",
			},
			out: RecipientsContacts{
				RecipientType_TO: {
					{
						Name:  "Henry Jones",
						Email: "h.jones@barnett.edu",
					},
				},
				RecipientType_CC:  {},
				RecipientType_BCC: {},
			},
		},
		{
			desc: "fail if name is empty",
			mailTemplate: &MailTemplate{
				Recipients: []*Recipient{
					{
						Type: RecipientType_TO,
						ContactType: &Recipient_FormData{
							FormData: &FieldNamesContact{
								NameField:  "name",
								EmailField: "email",
							},
						},
					},
				},
			},
			formData: map[string]any{
				"name":  "",
				"email": "h.jones@barnett.edu",
			},
			wantError: true,
		},
		{
			desc: "fail if there is no name",
			mailTemplate: &MailTemplate{
				Recipients: []*Recipient{
					{
						Type: RecipientType_TO,
						ContactType: &Recipient_FormData{
							FormData: &FieldNamesContact{
								NameField:  "name",
								EmailField: "email",
							},
						},
					},
				},
			},
			formData: map[string]any{
				"email": "h.jones@barnett.edu",
			},
			wantError: true,
		},
		/////////////////////////////////////////////////////////////////////
		{
			mailTemplate: &MailTemplate{
				Recipients: []*Recipient{
					{
						Type: RecipientType_TO,
						ContactType: &Recipient_ContactItem{
							ContactItem: &FieldNamesContact{
								NameField:  "name",
								EmailField: "email",
							},
						},
					},
				},
			},
			recipientItem: map[string][]string{
				"name":  {"Henry Jones"},
				"email": {"h.jones@barnett.edu"},
			},
			out: RecipientsContacts{
				RecipientType_TO: {
					{
						Name:  "Henry Jones",
						Email: "h.jones@barnett.edu",
					},
				},
				RecipientType_CC:  {},
				RecipientType_BCC: {},
			},
		},
		{
			mailTemplate: &MailTemplate{
				Recipients: []*Recipient{
					{
						Type: RecipientType_TO,
						ContactType: &Recipient_ContactItem{
							ContactItem: &FieldNamesContact{
								NameField:  "name",
								EmailField: "email",
							},
						},
					},
				},
			},
			recipientItem: map[string][]string{
				"name":  {},
				"email": {"h.jones@barnett.edu"},
			},
			wantError: true,
		},
	}

	for _, test := range cases {
		out, err := test.mailTemplate.DetermineRecipientsNameAndEmail(test.formData, test.recipientItem, test.contextItem)
		if test.wantError {
			if err == nil {
				t.Error("expected error but got nil")
			}
		} else {
			require.Equal(t, test.out, out, test.desc)
		}
	}
}
