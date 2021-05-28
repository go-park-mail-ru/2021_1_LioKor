package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/mail"
	"os"
	"os/signal"
	"syscall"
	"time"
	"strings"

	"github.com/emersion/go-smtp"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/utils"
	"liokor_mail/internal/pkg/mail/repository"
	liokorMail "liokor_mail/internal/pkg/mail"

	"github.com/microcosm-cc/bluemonday"
)

const CONFIG_PATH = "config.json"

var dbConn repository.GormPostgresMailRepository

type Backend struct{
	Config common.Config
}

func (bkd *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	return &Session{ Config: bkd.Config }, nil
}
func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return &Session{ Config: bkd.Config }, nil
}

type Session struct {
	From       string
	Recipients []string
	Header     mail.Header
	Body       string
	Config     common.Config
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

	pStrict := bluemonday.StrictPolicy()
	subject = pStrict.Sanitize(subject)

	pUGC := bluemonday.UGCPolicy()
	body = pUGC.Sanitize(body)

	for _, recipient := range s.Recipients {
		if strings.HasSuffix(recipient, "@" + s.Config.MailDomain) {
			newMail := liokorMail.Mail{
				Sender : s.From,
				Recipient: recipient,
				Subject: subject,
				Body: body,
			}
			_, err := dbConn.AddMail(newMail, s.Config.MailDomain)
			if err != nil {
				switch err.(type) {
				case common.InvalidUserError:
					break
				default:
					log.Printf("WARN: Mail wasn't added to database: %v\n", err)
				}
			}
		} else {
			log.Printf("WARN: Mail to %s was not saved!", recipient)
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
	err := config.ReadFromFile(CONFIG_PATH)
	if err != nil {
		log.Fatal(err)
	}

	db, err := common.NewGormPostgresDataBase(config)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()
	dbConn.DBInstance = db

	b := &Backend{ Config: config }
	s := smtp.NewServer(b)

	s.Addr = fmt.Sprintf("%s:%d", config.SmtpHost, config.SmtpPort)
	s.Domain = config.MailDomain
	s.ReadTimeout = 30 * time.Second
	s.WriteTimeout = 30 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AuthDisabled = true

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		log.Printf("Starting SMTP server at %s for @%s", s.Addr, config.MailDomain)
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal("Error occured while trying to start server: " + err.Error())
		}
		log.Println("Server was shut down with no errors!")
	}()
	<-quit

	log.Println("Interrupt signal received. Shutting down server...")
	if err := s.Close(); err != nil {
		log.Fatal("Server closed with and error: " + err.Error())
	}
}
