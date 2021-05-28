package usecase

import (
	"errors"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/mail"
	"liokor_mail/internal/utils"
	"log"
	"strings"
	"time"

	"crypto/rsa"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"

	"html"
)

type MailUseCase struct {
	Repository mail.MailRepository
	Config     common.Config
	PrivateKey *rsa.PrivateKey
}

func (uc *MailUseCase) GetDialogues(username string, amount int, find string, folderId int, since time.Time) ([]mail.Dialogue, error) {
	var dialogues []mail.Dialogue
	var err error
	if find == "" {
		dialogues, err = uc.Repository.GetDialoguesInFolder(username, amount, folderId, ("@" + uc.Config.MailDomain), since)
	} else {
		dialogues, err = uc.Repository.FindDialogues(username, find, amount, ("@" + uc.Config.MailDomain), since)
	}
	if err != nil {
		return nil, err
	}
	return dialogues, nil
}

func (uc *MailUseCase) CreateDialogue(owner, with string) (mail.Dialogue, error) {
	dialogue, err := uc.Repository.CreateDialogue(owner, with)
	if err != nil {
		return mail.Dialogue{}, err
	}
	return dialogue, nil
}

func (uc *MailUseCase) DeleteDialogue(owner string, dialogueId int) error {
	err := uc.Repository.DeleteDialogue(owner, dialogueId, uc.Config.MailDomain)
	if err != nil {
		return err
	}
	return nil
}

func (uc *MailUseCase) GetEmails(username string, email string, last int, amount int) ([]mail.DialogueEmail, error) {
	emails, err := uc.Repository.GetMailsForUser(username + "@" + uc.Config.MailDomain, email, amount, last)
	if err != nil {
		return nil, err
	}
	err = uc.Repository.ReadMail(username + "@" + uc.Config.MailDomain, email)
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
	email.Body = html.EscapeString(email.Body)

	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	email.Body = strings.ReplaceAll(email.Body, "&gt;", ">") // for quotes to work
	email.Body = strings.ReplaceAll(email.Body, "\n", "\n\n") // for newlines
	md := []byte(email.Body)
	email.Body = string(markdown.ToHTML(md, parser, nil))

	email.Body = strings.ReplaceAll(email.Body, "&amp;gt;", "&gt;") // for unescape to work
	email.Body = html.UnescapeString(email.Body)

	// 2-nd layer of sec - just in case smth above breaks
	pUGC := bluemonday.UGCPolicy()
	email.Body = pUGC.Sanitize(email.Body)

	if len(email.Subject) == 0 || len(email.Body) == 0 {
		return email, errors.New("Empty subject or body after sanitizing!")
	}

	mailId, err := uc.Repository.AddMail(email, uc.Config.MailDomain)
	if err != nil {
		return email, err
	}
	email.Id = mailId
	email.Status = 1

	if !isInternal {
		err = utils.SMTPSendMail(email.Sender, email.Recipient, email.Subject, email.Body, uc.PrivateKey)
		if err != nil {
			log.Printf("WARN: Unable to send email to %s\n", email.Recipient)
			email.Status = 0
			errDb := uc.Repository.UpdateMailStatus(mailId, 0)
			if errDb != nil {
				log.Printf("ERROR: Unable to change mail status!\n")
				return email, err
			}
		}
	}

	return email, nil
}

func (uc *MailUseCase) DeleteMails(owner string, mailIds []int) error{
	err := uc.Repository.DeleteMail(owner, mailIds, uc.Config.MailDomain)
	if err != nil {
		return err
	}
	return nil
}

func (uc *MailUseCase) GetFolders(owner int) ([]mail.Folder, error) {
	folders, err := uc.Repository.GetFolders(owner)
	if err != nil {
		return nil, err
	}
	return folders, nil
}

func (uc *MailUseCase) CreateFolder(owner int, folderName string) (mail.Folder, error) {
	folder, err := uc.Repository.CreateFolder(owner, folderName)
	if err != nil {
		return mail.Folder{}, err
	}
	return folder, nil
}

func (uc *MailUseCase) UpdateFolderPutDialogue(owner string, folderId int, dialogueId int) error {
	err := uc.Repository.AddDialogueToFolder(owner, folderId, dialogueId)
	if err != nil {
		return err
	}
	return nil
}

func (uc *MailUseCase) UpdateFolderName(owner, folderId int, folderName string) (mail.Folder, error) {
	folder, err := uc.Repository.UpdateFolderName(owner, folderId, folderName)
	if err != nil {
		return mail.Folder{}, err
	}
	return folder, nil
}

func (uc *MailUseCase) DeleteFolder(ownerName string, owner, folderId int) error {
	err := uc.Repository.ShiftToMainFolderDialogues(ownerName, folderId)
	if err != nil {
		return err
	}
	err = uc.Repository.DeleteFolder(owner, folderId)
	if err != nil {
		return err
	}
	return nil
}
