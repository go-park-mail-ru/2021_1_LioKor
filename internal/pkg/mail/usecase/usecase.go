package usecase

import (
	"errors"
	"fmt"
	"liokor_mail/internal/pkg/mail"
	"log"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type MailUseCase struct {
	Repository mail.MailRepository
}

func (uc *MailUseCase) SMTPSendMail(from string, to string, subject string, data string) error {
	recipientSplitted := strings.Split(to, "@")
	if len(recipientSplitted) != 2 {
		return errors.New("invalid recipient address!")
	}
	hostAddr := recipientSplitted[1]
	mxrecords, err := net.LookupMX(hostAddr)
	if err != nil {
		log.Println(err)
		return err
	}

	host := mxrecords[0].Host
	host = host[:len(host)-1]

	mail := fmt.Sprintf("From: <%s>\r\nTo: %s\r\nContent-Type: text/plain\r\nSubject: %s\r\n\r\n%s\r\n", from, to, subject, data)
	err = smtp.SendMail(host+":25", nil, from, []string{to}, []byte(mail))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (uc *MailUseCase) GetDialogues(username string, last int, amount int, find string) ([]mail.Dialogue, error) {
	username += "@liokor.ru"
	dialogues, err := uc.Repository.GetDialoguesForUser(username, amount, last, find)
	if err != nil {
		return nil, err
	}
	return dialogues, nil
}

func (uc *MailUseCase) GetEmails(username string, email string, last int, amount int) ([]mail.DialogueEmail, error) {
	username += "@liokor.ru"
	emails, err := uc.Repository.GetMailsForUser(username, email, amount, last)
	if err != nil {
		return nil, err
	}
	return emails, nil
}

func (uc *MailUseCase) SendEmail(email mail.Mail) error {
	email.Sender += "@liokor.ru"

	lastMailsCount, err := uc.Repository.CountMailsFromUser(email.Sender, 3*time.Minute)
	if err != nil {
		return err
	}
	if lastMailsCount > 5 {
		return mail.InvalidEmailError{"too many mails, wait some time"}
	}

	err = uc.SMTPSendMail(email.Sender, email.Recipient, email.Subject, email.Body)
	if err != nil {
		return err
	}

	err = uc.Repository.AddMail(email)
	if err != nil {
		return err
	}

	return nil
}
