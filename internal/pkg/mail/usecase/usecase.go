package usecase

import "liokor_mail/internal/pkg/mail"

type MailUseCase struct {
	Repository mail.MailRepository
}

func (uc *MailUseCase)GetDialogues(username string, last int, amount int) ([]mail.Dialogue, error) {
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
	//TODO: call function to send email

	return nil
}
