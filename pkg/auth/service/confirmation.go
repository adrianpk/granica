package service

import (
	"fmt"

	"gitlab.com/mikrowezel/backend/granica/internal/model"
)

func (s *Service) MakeConfirmationEmail(u *model.User) model.Email {
	cfg := s.Cfg()

	name := cfg.ValOrDef("mailer.agent.name", "mailer")
	from := cfg.ValOrDef("mailer.agent.mail", "dontreply@localhost")
	to := u.Email.String
	subject := fmt.Sprintf("%s, please confirm your account!", u.Username.String)

	site := cfg.ValOrDef("site.url", "localhost")
	path := cfg.ValOrDef("user.confirmation.path", "users/%s/verify/%s")
	confPath := fmt.Sprintf(path, u.Slug.String, u.ConfirmationToken.String)
	link := fmt.Sprintf("htts://%s/%s", site, confPath)

	body := "<p>Hi %s, follow this link to confirm your account: <br/><br/>"
	body = body + "<a href=\"%s\">%s</a><br/<br/>"
	body = body + "Thanks!"
	body = fmt.Sprintf(body, u.Username.String, link, link)

	m := model.MakeEmail(name, from, to, "", "", subject, body)

	s.Log().Info("User account confirmation", "mail-body", fmt.Sprintf("%+v", m))

	return m
}
