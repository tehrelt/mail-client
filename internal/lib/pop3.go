package lib

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/knadh/go-pop3"
	"io"
	"mail-client/internal/config"
	"strings"
	"time"
)

//
//func IsPop3Error(message string) bool {
//	return strings.Contains(message, pop3.ERR)
//}

var (
	ErrPop3Disconnected = errors.New("disconnected")
)

type Pop3 struct {
	client     *pop3.Client
	connection *pop3.Conn
}

type Mail struct {
	From    string         `json:"from"`
	To      string         `json:"to"`
	Subject string         `json:"subject"`
	Body    string         `json:"body"`
	Date    time.Time      `json:"date"`
	Meta    pop3.MessageID `json:"meta"`
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

func (p *Pop3) connectionAlive() error {
	if p.connection == nil {
		return ErrPop3Disconnected
	}
	return nil
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
	if err := p.connectionAlive(); err != nil {
		return nil, err
	}

	msgs, err := p.connection.List(0)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (p *Pop3) Stat() (int, int, error) {
	if err := p.connectionAlive(); err != nil {
		return 0, 0, err
	}
	return p.connection.Stat()
}

func (p *Pop3) Retrieve(id int) (*Mail, error) {
	if err := p.connectionAlive(); err != nil {
		return nil, err
	}

	message, err := p.connection.Retr(id)
	if err != nil {
		return nil, err
	}

	reader := message.Body
	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	if err != nil {
		log.Fatal(err)
	}

	f := message.Header.Get("from")
	from, err := base64.StdEncoding.DecodeString(strings.Split(strings.Split(f, "=?UTF-8?B?")[1], "?=")[0])
	if err != nil {
		log.Fatal("error:", err)
	}

	sender := fmt.Sprintf("%s %s", string(from), strings.Split(f, " ")[1])

	subj := message.Header.Get("subject")

	subjj := strings.Split(subj, "UTF-8?B?")

	var subject []string

	for i := 1; i < len(subjj); i++ {
		su, err := base64.StdEncoding.DecodeString(strings.Split(subjj[i], "?=")[0])
		if err != nil {
			log.Fatal("error:", err)
		}

		subject = append(subject, string(su))
	}

	d := message.Header.Get("Date")
	t, err := time.Parse("02 Jan 2006 15:04:05 -0700", strings.Split(d, ", ")[1])
	if err != nil {
		return nil, err
	}

	elems := strings.Split(buf.String(), "\r\n\r\n----ALT")[0]
	encodedParts := strings.Split(elems, "\r\n")[5:]

	encoded := strings.Join(encodedParts, "")

	content, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	return &Mail{
		From:    sender,
		To:      message.Header.Get("to"),
		Subject: strings.Join(subject, ""),
		Body:    string(content),
		Date:    t,
	}, nil
}
