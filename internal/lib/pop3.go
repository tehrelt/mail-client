package lib

import (
	"encoding/base64"
	"errors"
	"fmt"
	_ "github.com/emersion/go-message"
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

type Attachment struct {
	FileName  string `json:"fileName"`
	MediaType string `json:"mediaType"`
	Data      string `json:"data"`
}

type Mail struct {
	From        string         `json:"from"`
	To          string         `json:"to"`
	Subject     string         `json:"subject"`
	Body        string         `json:"body"`
	Date        time.Time      `json:"date"`
	Meta        pop3.MessageID `json:"meta"`
	Attachments []Attachment   `json:"attachments,omitempty"`
}

func NewPop(connection *pop3.Conn) *Pop3 {
	return &Pop3{connection}
}

func parseB64UTF8(s string) ([]byte, error) {
	b64 := strings.Split(strings.Split(s, "UTF-8?B?")[1], "?=")[0]
	r, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}
	return r, nil
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

func (p *Pop3) RetrWithAttachments(id int) (*Mail, error) {
	msg, err := p.Retr(id)
	if err != nil {
		return nil, err
	}

	f := msg.Header.Get("From")
	from, err := parseB64UTF8(f)
	if err != nil {
		return nil, err
	}

	subj := msg.Header.Get("Subject")
	subjj := strings.Split(subj, "UTF-8?B?")

	var subject []string

	for i := 1; i < len(subjj); i++ {
		su, err := base64.StdEncoding.DecodeString(strings.Split(subjj[i], "?=")[0])
		if err != nil {
			return nil, err
		}

		subject = append(subject, string(su))
	}

	d := msg.Header.Get("Date")
	t, err := time.Parse("02 Jan 2006 15:04:05 -0700", strings.Split(d, ", ")[1])
	if err != nil {
		return nil, err
	}

	var mail Mail

	mail.From = fmt.Sprintf("%s %s", string(from), strings.Split(f, " ")[1])
	mail.Subject = strings.Join(subject, "")
	mail.Date = t
	mail.To = msg.Header.Get("to")

	mr := msg.MultipartReader()
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		mediaType := p.Header.Get("Content-Type")

		slurp, err := io.ReadAll(p.Body)
		if err != nil {
			return nil, err
		}

		cType := strings.Split(mediaType, ";")[0]

		if strings.Compare(cType, "text/html") == 0 {
			mail.Body = string(slurp)
		} else if strings.HasPrefix(cType, "text/") {
			continue
		} else if strings.HasPrefix(cType, "multipart/") {
			//fmt.Printf("\t\t%s\n", mediaType)
		} else {
			fmt.Println(mediaType)

			var name string

			if strings.Contains(mediaType, "UTF-8") {
				b, err := parseB64UTF8(mediaType)
				if err != nil {
					return nil, err
				}
				name = string(b)
			} else {
				name = strings.Split(mediaType, "=\"")[1]
				name = strings.Split(name, "\"")[0]
			}

			//fmt.Printf("Has part %q: %q\n", cType, name)

			data := base64.StdEncoding.EncodeToString(slurp)
			//fmt.Printf("%q ---> %q\n\n\n", slurp, data)
			mail.Attachments = append(mail.Attachments, Attachment{
				FileName:  name,
				MediaType: cType,
				Data:      data,
			})
		}

		if err != nil {
			return nil, err
		}

	}

	return &mail, nil
}

func (p *Pop3) Retrieve(id int) (*Mail, error) {
	msg, err := p.Retr(id)
	if err != nil {
		return nil, err
	}

	f := msg.Header.Get("From")
	from, err := parseB64UTF8(f)
	if err != nil {
		return nil, err
	}

	subj := msg.Header.Get("Subject")
	subjj := strings.Split(subj, "UTF-8?B?")

	var subject []string

	for i := 1; i < len(subjj); i++ {
		su, err := base64.StdEncoding.DecodeString(strings.Split(subjj[i], "?=")[0])
		if err != nil {
			return nil, err
		}

		subject = append(subject, string(su))
	}

	d := msg.Header.Get("Date")
	t, err := time.Parse("02 Jan 2006 15:04:05 -0700", strings.Split(d, ", ")[1])
	if err != nil {
		return nil, err
	}

	var mail Mail

	mail.From = fmt.Sprintf("%s %s", string(from), strings.Split(f, " ")[1])
	mail.Subject = strings.Join(subject, "")
	mail.Date = t
	mail.To = msg.Header.Get("to")

	mr := msg.MultipartReader()
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}

		mediaType := p.Header.Get("Content-Type")

		slurp, err := io.ReadAll(p.Body)
		if err != nil {
			return nil, err
		}

		cType := strings.Split(mediaType, ";")[0]

		if strings.HasPrefix(cType, "text/") {
			if strings.Compare(cType, "text/plain") == 0 {
				mail.Body = string(slurp)
			}
		}

		if err != nil {
			return nil, err
		}

	}

	return &mail, nil
}
