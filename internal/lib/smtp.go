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
	config *config.SmtpConfig
	user   *dto.User
}

func NewSmtp(cfg *config.SmtpConfig) *Smtp {
	return &Smtp{
		config: cfg,
	}
}

func (s *Smtp) Auth(user *dto.User) error {
	dialer := gomail.NewDialer(
		s.config.Host,
		s.config.Port,
		user.User,
		user.Pass,
	)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	sender, err := dialer.Dial()
	if err != nil {
		return err
	}
	defer sender.Close()

	s.user = user

	return nil
}

func (s *Smtp) Send(message *dto.Message) error {

	if s.user == nil {
		return ErrSmtpForbidden
	}

	dialer := gomail.NewDialer(
		s.config.Host,
		s.config.Port,
		s.user.User,
		s.user.Pass,
	)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", s.user.User)
	m.SetHeader("To", strings.Join(message.To, ","))
	m.SetHeader("Subject", message.Subject)
	m.AddAlternative("text/plain", message.Body)

	if err := dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
