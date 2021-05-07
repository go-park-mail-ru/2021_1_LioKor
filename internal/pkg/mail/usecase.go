package mail

type MailUseCase interface {
	GetDialogues(username string, amount int, find string, folderId int) ([]Dialogue, error)
	DeleteDialogue(owner string, dialogueId int) error
	GetEmails(username string, email string, last int, amount int) ([]DialogueEmail, error)
	SendEmail(mail Mail) (Mail, error)
	GetFolders(owner int) ([]Folder, error)
	CreateFolder(owner int, folderName string) (Folder, error)
	UpdateFolderPutDialogue(owner string, folderId int, dialogueId int) error
	UpdateFolderName(owner, folderId int, folderName string) (Folder, error)
	DeleteFolder(ownerName string, owner, folderId int) error
}
