package mail

import "time"

type MailRepository interface {
	AddMail(mail Mail, domain string) (int, error)
	GetMailsForUser(username string, email string, limit int, last int) ([]DialogueEmail, error)
	ReadMail(owner, other string) error
	CountMailsFromUser(username string, interval time.Duration) (int, error)
	UpdateMailStatus(mailId, status int) error
	DeleteMail(owner string, mailIds []int, domain string) error

	CreateDialogue(owner string, other string, domain string) (Dialogue, error)
	UpdateDialogueLastMail(owner string, other string, domain string) error
	GetDialoguesInFolder(username string, limit int, folderId int, domain string, since time.Time) ([]Dialogue, error)
	FindDialogues(username string, find string, limit int, domain string, since time.Time) ([]Dialogue, error)
	ReadDialogue(owner, other string) error
	DeleteDialogue(owner string, dialogueId int, domain string) error

	CreateFolder(ownerId int, folderName string) (Folder, error)
	GetFolders(ownerId int) ([]Folder, error)
	AddDialogueToFolder(owner string, folderId, dialogueId int) error
	UpdateFolderName(owner, folderId int, folderName string) (Folder, error)
	ShiftToMainFolderDialogues(owner string, folderId int) error
	DeleteFolder(owner, folderId int) error
}
