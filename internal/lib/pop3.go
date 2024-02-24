package lib

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/knadh/go-pop3"
	"mail-client/internal/config"
	"strings"
)

//
//func IsPop3Error(message string) bool {
//	return strings.Contains(message, pop3.ERR)
//}

type Pop3 struct {
	client     *pop3.Client
	connection *pop3.Conn
}

type Mail struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewPop(cfg *config.Pop3Config) *Pop3 {
	return &Pop3{
		client: pop3.New(pop3.Opt{
			Host:       cfg.Host,
			Port:       cfg.Port,
			TLSEnabled: false,
		}),
	}
}

func (p *Pop3) Auth(user, pass string) error {
	conn, err := p.client.NewConn()
	if err != nil {
		return err
	}

	if err := conn.Auth(user, pass); err != nil {
		conn.Quit()
		return err
	}

	p.connection = conn

	return nil
}

func (p *Pop3) ListAll() ([]pop3.MessageID, error) {
	if p.connection == nil {
		return nil, errors.New("disconnected")
	}

	msgs, err := p.connection.List(0)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (p *Pop3) Retrieve(messageInfo pop3.MessageID) (*Mail, error) {
	if p.connection == nil {
		return nil, errors.New("disconnected")
	}

	message, err := p.connection.Retr(messageInfo.ID)
	if err != nil {
		return nil, err
	}

	body := make([]byte, messageInfo.Size)
	_, err = message.Body.Read(body)
	if err != nil {
		return nil, err
	}

	f := message.Header.Get("from")
	from, err := base64.StdEncoding.DecodeString(strings.Split(strings.Split(f, "=?UTF-8?B?")[1], "?=")[0])
	if err != nil {
		log.Fatal("error:", err)
	}

	sender := fmt.Sprintf("%s %s", string(from), strings.Split(f, " ")[1])

	//content, err := base64.StdEncoding.DecodeString(string(body))
	//if err != nil {
	//	log.Fatal("error:", err)
	//}
	subj := message.Header.Get("subject")
	subject, err := base64.StdEncoding.DecodeString(strings.Split(strings.Split(subj, "UTF-8?B?")[1], "?=")[0])
	if err != nil {
		log.Fatal("error:", err)
	}

	encoded := strings.Split(string(body), "\r\n")
	for i, m := range encoded {
		log.Debugf("%d\t-\t%s", i, m)
	}

	m := encoded[5]

	content, err := base64.StdEncoding.DecodeString(m)
	if err != nil {
		log.Fatal("error:", err)
	}
	return &Mail{
		From:    sender,
		To:      message.Header.Get("to"),
		Subject: string(subject),
		Body:    string(content),
	}, nil
}
