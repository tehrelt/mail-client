package lib

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"mail-client/internal/config"
	"mail-client/internal/dto"
	"strings"
)

type Smtp struct {
	*gomail.Dialer
}

func NewSmtp(dialer *gomail.Dialer) *Smtp {
	return &Smtp{dialer}
}

func SmtpAuth(cfg *config.SmtpConfig, user *dto.User) (*Smtp, error) {
	dialer := gomail.NewDialer(
		cfg.Host,
		cfg.Port,
		user.User,
		user.Pass,
	)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	sender, err := dialer.Dial()
	if err != nil {
		return nil, err
	}
	sender.Close()

	return NewSmtp(dialer), nil
}

func (s *Smtp) SendMessage(message *dto.Message) error {

	m := gomail.NewMessage()
	m.SetHeader("From", message.From)
	m.SetHeader("To", strings.Join(message.To, ","))
	m.SetHeader("Subject", message.Subject)
	m.AddAlternative("text/plain", message.Body)
	for _, file := range message.Attachments {
		m.Attach(file)
	}

	if err := s.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
