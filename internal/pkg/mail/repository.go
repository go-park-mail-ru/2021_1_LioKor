package mail

import (
	"time"
)

type MailRepository interface {
	GetDialoguesForUser(username string, limit int, find string, folderId int, domain string) ([]Dialogue, error)
	DeleteDialogue(owner string, dialogueId int) error
	GetMailsForUser(username string, email string, limit int, last int) ([]DialogueEmail, error)
	AddMail(mail Mail) (int, error)
	CountMailsFromUser(username string, interval time.Duration) (int, error)
	ReadDialogue(owner, other string) error
	ReadMail(owner, other string) error
	CreateFolder(ownerId int, folderName string) (Folder, error)
	GetFolders(ownerId int) ([]Folder, error)
	AddDialogueToFolder(owner string, folderId, dialogueId int) error
	UpdateFolderName(owner, folderId int, folderName string) (Folder, error)
	DeleteFolder(owner, folderId int) error
	UpdateMailStatus(mailId, status int) error
}
