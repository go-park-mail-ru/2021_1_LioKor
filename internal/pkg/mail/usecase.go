package mail

type MailUseCase interface {
	GetDialogues(username string, amount int, find string, folderId int) ([]Dialogue, error)
	GetEmails(username string, email string, last int, amount int) ([]DialogueEmail, error)
	SendEmail(mail Mail) error
	GetFolders(owner int)([]Folder, error)
	CreateFolder(owner int, folderName string) (Folder, error)
	UpdateFolder(owner string, folderId int, dialogueId int) error
}
