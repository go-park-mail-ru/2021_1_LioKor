package usecase

import (
	"errors"
	"fmt"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/mail"
	"log"
	"net"
	"net/smtp"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

type MailUseCase struct {
	Repository mail.MailRepository
	Config     common.Config
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

func (uc *MailUseCase) GetDialogues(username string, amount int, find string, folderId int) ([]mail.Dialogue, error) {
	username += "@" + uc.Config.MailDomain
	dialogues, err := uc.Repository.GetDialoguesForUser(username, amount, find, folderId, ("@" + uc.Config.MailDomain))
	if err != nil {
		return nil, err
	}
	return dialogues, nil
}

func (uc *MailUseCase) GetEmails(username string, email string, last int, amount int) ([]mail.DialogueEmail, error) {
	username += "@" + uc.Config.MailDomain
	emails, err := uc.Repository.GetMailsForUser(username, email, amount, last)
	if err != nil {
		return nil, err
	}
	err = uc.Repository.ReadMail(username, email)
	if err != nil {
		return nil, err
	}
	err = uc.Repository.ReadDialogue(username, email)
	if err != nil {
		return nil, err
	}
	return emails, nil
}

func (uc *MailUseCase) SendEmail(email mail.Mail) (mail.Mail, error) {
	email.Sender += "@" + uc.Config.MailDomain
	isInternal := strings.HasSuffix(email.Recipient, uc.Config.MailDomain)

	if !(uc.Config.Debug || isInternal) {
		lastMailsCount, err := uc.Repository.CountMailsFromUser(email.Sender, 3*time.Minute)
		if err != nil {
			return email, err
		}
		if lastMailsCount > 5 {
			return email, mail.InvalidEmailError{"too many mails, wait some time"}
		}
	}

	pStrict := bluemonday.StrictPolicy()
	email.Subject = pStrict.Sanitize(email.Subject)

	pUGC := bluemonday.UGCPolicy()
	email.Body = pUGC.Sanitize(email.Body)

	if len(email.Subject) == 0 || len(email.Body) == 0 {
		return email, errors.New("Empty subject or body after sanitizing!")
	}

	mailId, err := uc.Repository.AddMail(email)
	if err != nil {
		return email, err
	}

	if !isInternal {
		err = uc.SMTPSendMail(email.Sender, email.Recipient, email.Subject, email.Body)
		if err != nil {
			err = uc.Repository.UpdateMailStatus(mailId, 0)
			return email, err
		}
	}

	return email, nil
}

func (uc *MailUseCase) GetFolders(owner int)([]mail.Folder, error) {
	folders, err := uc.Repository.GetFolders(owner)
	if err != nil {
		return nil, err
	}
	return folders, nil
}

func (uc *MailUseCase) CreateFolder(owner int, folderName string) (mail.Folder, error){
	folder, err := uc.Repository.CreateFolder(owner, folderName)
	if err != nil {
		return mail.Folder{}, err
	}
	return folder, nil
}

func (uc *MailUseCase) UpdateFolder(owner string, folderId int, dialogueId int) error{
	owner += "@" + uc.Config.MailDomain

	err := uc.Repository.AddDialogueToFolder(owner, folderId, dialogueId)
	if err != nil {
		return err
	}
	return nil

}
