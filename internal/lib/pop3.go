package lib

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"github.com/knadh/go-pop3"
	"io"
	"mail-client/internal/config"
	"mail-client/internal/dto"
	"strings"
	"time"
)

var (
	ErrPop3Disconnected = errors.New("disconnected")
	ErrSmtpForbidden    = errors.New("forbidden")
)

type Pop3 struct {
	*pop3.Conn
}

type part struct {
	ContentType string `json:"contentType"`
	Charset     string `json:"charset"`
	Body        string `json:"body"`
}

type Mail struct {
	From    string         `json:"from"`
	To      string         `json:"to"`
	Subject string         `json:"subject"`
	Body    []part         `json:"body"`
	Date    time.Time      `json:"date"`
	Meta    pop3.MessageID `json:"meta"`
}

func NewPop(connection *pop3.Conn) *Pop3 {
	return &Pop3{connection}
}

func Pop3Auth(config *config.Pop3Config, user *dto.User) (*Pop3, error) {
	client := pop3.New(pop3.Opt{
		Host:          config.Host,
		Port:          config.Port,
		TLSEnabled:    false,
		TLSSkipVerify: false,
	})

	conn, err := client.NewConn()
	if err != nil {
		return nil, err
	}

	if err := conn.Auth(user.User, user.Pass); err != nil {
		conn.Quit()
		return nil, err
	}

	return NewPop(conn), nil
}

func (p *Pop3) ListAll() ([]pop3.MessageID, error) {

	msgs, err := p.List(0)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (p *Pop3) Retrieve(id int) (*Mail, error) {
	message, err := p.Retr(id)
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

	log.Debug(strings.Join(subject, ""))

	var parts []part
	if strings.Contains(buf.String()[:10], "ALT") {
		key := buf.String()[:55]

		elems := strings.Split(buf.String(), key)
		for i, e := range elems {
			if i == 0 {
				continue
			}
			//log.Debugf("%d:%s\n", i, e)

			var p part
			p.ContentType = strings.Split(strings.Split(e, "Content-Type: ")[1], ";")[0]
			p.Charset = strings.Split(strings.Split(e, "charset=")[1], "\r\n")[0]
			p.Body = strings.Join(
				strings.Split(
					strings.Split(
						strings.Split(e, "\r\n\r\n")[1], "\r\n\r\n")[0],
					"\r\n"),
				"")

			content, err := base64.StdEncoding.DecodeString(p.Body)
			if err != nil {
				return nil, err
			}

			p.Body = string(content)

			parts = append(parts, p)
		}
	} else {

	}

	//log.Debug(buf.String())

	//elems := strings.Split(buf.String(), "\r\n\r\n----ALT--")
	//for _, elem := range elems {
	//	ctype := strings.Split(strings.Split(elem, "Content-Type: ")[1], ";")[0]
	//	ctype = strings.Replace(ctype, "/", "_", -1)
	//	if err := os.WriteFile(fmt.Sprintf("%s_%s.txt", strings.Join(subject, ""), ctype), []byte(elem), 777); err != nil {
	//		return nil, err
	//	}
	//}

	//elems := strings.Split(buf.String(), "\r\n\r\n----ALT")[0]
	//encodedParts := strings.Split(elems, "\r\n")[5:]
	//
	//encoded := strings.Join(encodedParts, "")
	//
	//content, err := base64.StdEncoding.DecodeString(encoded)
	//if err != nil {
	//	return nil, err
	//}

	return &Mail{
		From:    sender,
		To:      message.Header.Get("to"),
		Subject: strings.Join(subject, ""),
		Body:    parts,
		Date:    t,
	}, nil
}
