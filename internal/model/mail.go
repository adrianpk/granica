package model

import (
	uuid "github.com/satori/go.uuid"
	m "gitlab.com/mikrowezel/backend/model"
)

// Email struct
type Email struct {
	m.Identification
	Name    string
	From    string
	To      string
	CC      string
	BCC     string
	Subject string
	Body    string
	Charset string
}

// MakeEmail
func MakeEmail(name, from, to, cc, bcc, subject, body string) Email {
	return Email{
		Identification: m.Identification{
			ID: uuid.NewV4(),
		},
		Name:    name,
		From:    from,
		To:      to,
		CC:      cc,
		BCC:     bcc,
		Subject: subject,
		Body:    body,
		Charset: "UTF-8",
	}
}
