package main

import (
	"context"
	"io"
	"log"
	"net/mail"
	"time"
	"errors"

	"github.com/emersion/go-smtp"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/mail_server/utils"
)

var db common.PostgresDataBase

type Backend struct{}

func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return &Session{}, nil
}
func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return &Session{}, nil
}

type Session struct {
	From       string
	Recipients []string
	Header     mail.Header
	Body       string
}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	s.From = from
	return nil
}

func (s *Session) Rcpt(recipient string) error {
	s.Recipients = append(s.Recipients, recipient)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	message, err := mail.ReadMessage(r)
	if err != nil {
		log.Println(err)
		return err
	}

	body, err := utils.ParseBodyText(message)
	if err != nil {
		log.Println(err)
		return err
	}

	s.Header = message.Header
	s.Body = body

	return s.HandleMail()
}

func (s *Session) HandleMail() error {
	if len(s.From) == 0 || len(s.Recipients) == 0 || len(s.Body) == 0 {
		log.Println("Invalid mail received!")
		return errors.New("Invalid mail received!")
	}

	log.Printf("Received mail from %s to %v\n", s.From, s.Recipients)

	subject := utils.ParseSubject(s.Header.Get("Subject"))
	body := s.Body

	for _, recipient := range s.Recipients {
		_, err := db.DBConn.Exec(
			context.Background(),
			"INSERT INTO mails(sender, recipient, subject, body) VALUES($1, $2, $3, $4);",
			s.From,
			recipient,
			subject,
			body,
		)
		if err != nil {
			log.Println(err)
		}
	}
	return nil
}

func (s *Session) Reset() {
	s.From = ""
	s.Recipients = nil
	s.Body = ""
}

func (s *Session) Logout() error {
	return nil
}

func main() {
	config := common.Config{}
	err := config.ReadFromFile("config.json")
	if err != nil {
		log.Fatal(err)
	}

	db, err = common.NewPostgresDataBase(config.DbString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	b := &Backend{}
	s := smtp.NewServer(b)

	s.Addr = ":25"
	s.Domain = config.MailDomain
	s.ReadTimeout = 30 * time.Second
	s.WriteTimeout = 30 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AuthDisabled = true

	log.Printf("Starting SMTP server at %s for @%s", s.Addr, config.MailDomain)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
