package mail

type MailUseCase interface {
	GetDialogues(username string, last int, amount int) ([]Dialogue, error)
	GetEmails(username string, email string, last int, amount int) ([]DialogueEmail, error)
	SendEmail(mail Mail) error
}