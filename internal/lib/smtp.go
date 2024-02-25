package lib

import (
	"crypto/tls"
	"errors"
	"gopkg.in/gomail.v2"
	"mail-client/internal/config"
	"mail-client/internal/dto"
	"net/smtp"
	"strings"
)

type loginAuth struct {
	username, password string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

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

	if err := s.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
