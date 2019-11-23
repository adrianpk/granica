package web

import "github.com/nicksnyder/go-i18n/v2/i18n"

var (
	// NOTE:  Reference i18n message.
	messages = &i18n.Message{
		ID:          "Emails",
		Description: "The number of unread emails a user has",
		One:         "{{.Name}} has {{.Count}} email.",
		Other:       "{{.Name}} has {{.Count}} emails.",
	}
)
