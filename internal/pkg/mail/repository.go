package mail

type MailRepository interface {
	GetDialoguesForUser(username string, limit int, offset int) ([]Dialogue, error)
	GetMailsForUser(username string, email string, limit int, offset int) ([]DialogueEmail, error)
}
