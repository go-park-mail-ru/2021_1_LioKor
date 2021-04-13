package usecase

import (
	"liokor_mail/internal/pkg/mail"
	"net/smtp"
	"net"
	"strings"
	"fmt"
)

type MailUseCase struct {
	Repository mail.MailRepository
}

func (uc *MailUseCase) SMTPSendMail(from string, to string, subject string, data string) error {
    addr := strings.Split(to, "@")[1]
    mxrecords, err := net.LookupMX(addr)
    if err != nil {
        return err
    }

    host := mxrecords[0].Host
    host = host[:len(host) - 1]

	mail := fmt.Sprintf("From: <%s>\r\nTo: %s\r\nContent-Type: text/plain\r\nSubject: %s\r\n\r\n%s\r\n", from, to, subject, data);
    err = smtp.SendMail(host + ":25", nil, from, []string{to}, []byte(mail))
    if err != nil {
        return err
    }
    return nil
}

func (uc *MailUseCase) GetDialogues(username string, last int, amount int) ([]mail.Dialogue, error) {
	username += "@liokor.ru"
	dialogues, err := uc.Repository.GetDialoguesForUser(username, amount, last)
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

func (uc *MailUseCase) SendEmail(mail mail.Mail) error {
	mail.Sender += "@liokor.ru"

	err := uc.Repository.AddMail(mail)
	if err != nil {
		return err
	}

	//TODO: отправлять письмо до тех пор, пока не выйдет, или уже выдавать пользователю ошибку
	//например, в mail письмо сохраняется, а на ошибку отправляет письмо
	err = uc.SMTPSendMail(mail.Sender, mail.Recipient, mail.Subject, mail.Body)
	if err != nil {
		return err
	}

	return nil
}
