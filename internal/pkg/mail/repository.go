package mail

import (
	"time"
)

type MailRepository interface {
	GetDialoguesForUser(username string, limit int, last int, find string) ([]Dialogue, error)
	GetMailsForUser(username string, email string, limit int, last int) ([]DialogueEmail, error)
	AddMail(mail Mail) error
	CountMailsFromUser(username string, interval time.Duration) (int, error)
}
