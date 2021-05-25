package mail

import "time"

type MailUseCase interface {
	GetDialogues(username string, amount int, find string, folderId int, since time.Time) ([]Dialogue, error)
	CreateDialogue(owner, with string) (Dialogue, error)
	DeleteDialogue(owner string, dialogueId int) error
	GetEmails(username string, email string, last int, amount int) ([]DialogueEmail, error)
	SendEmail(mail Mail) (Mail, error)
	DeleteMails(owner string, mailIds []int) error
	GetFolders(owner int) ([]Folder, error)
	CreateFolder(owner int, folderName string) (Folder, error)
	UpdateFolderPutDialogue(owner string, folderId int, dialogueId int) error
	UpdateFolderName(owner, folderId int, folderName string) (Folder, error)
	DeleteFolder(ownerName string, owner, folderId int) error
}
