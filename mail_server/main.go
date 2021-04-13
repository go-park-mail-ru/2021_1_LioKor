package main

import (
	"bytes"
	"io"
	"log"
	"net/mail"
	"os"
	"time"
    "context"

    "liokor_mail/mail_server/utils"
    "liokor_mail/internal/pkg/common"
	"github.com/emersion/go-smtp"
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
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)

	data, _ := io.ReadAll(tee)
	_ = os.WriteFile("mail_latest.txt", data, 0744)

	message, err := mail.ReadMessage(&buf)
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

	return nil
}

func (s *Session) Reset() {
    log.Printf("Received mail from %s to %v\n", s.From, s.Recipients)

	subject := s.Header.Get("Subject")
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
	s.Domain = "liokor.ru"
	s.ReadTimeout = 30 * time.Second
	s.WriteTimeout = 30 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AuthDisabled = true

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
