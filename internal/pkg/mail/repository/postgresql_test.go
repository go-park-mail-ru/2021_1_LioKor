package repository

import (
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/mail"
	"liokor_mail/internal/pkg/user"
	userRep "liokor_mail/internal/pkg/user/repository"
	"testing"
	"time"
)

var dbConfig = common.Config{
	DBHost:           "127.0.0.1",
	DBPort:           5432,
	DBUser:           "liokor",
	DBPassword:       "Qwerty123",
	DBDatabase:       "liokor_mail",
	DBConnectTimeout: 10,
}

var owner user.User
var other user.User
var folder mail.Folder
var dialogues []mail.Dialogue


func TestCreateUsers(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Fatalf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	ur := userRep.PostgresUserRepository{
		dbInstance,
	}

	randomString := common.GenerateRandomString()

	err = ur.CreateUser(user.User{
			Username:     randomString[:8],
			HashPassword: "hashPassword",
	})
	if err != nil {
		t.Fatalf("Database error: %v\n", err)
	}
	err = ur.CreateUser(user.User{
		Username:     randomString[8:16],
		HashPassword: "hashPassword",
	})
	if err != nil {
		t.Fatalf("Database error: %v\n", err)
	}
	owner, err = ur.GetUserByUsername(randomString[:8])
	if err != nil {
		t.Fatalf("Database error: %v\n", err)
	}
	other, err = ur.GetUserByUsername(randomString[8:16])
	if err != nil {
		t.Fatalf("Database error: %v\n", err)
	}

}

func TestAddMail(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	mailRep := PostgresMailRepository{
		dbInstance,
	}

	newMail := mail.Mail{
		Sender:    owner.Username,
		Recipient: other.Username,
		Subject:   "Test mail",
		Body:      "Test adding mails",
	}

	_, err = mailRep.AddMail(newMail)
	if err != nil {
		t.Errorf("Didn't pass adding mail: %v\n", err)
	}
}

func TestGetMailsForUser(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	mailRep := PostgresMailRepository{
		dbInstance,
	}

	mails, err := mailRep.GetMailsForUser(owner.Username, other.Username, 10, 0)
	if err != nil {
		t.Errorf("Didn't pass valid get emails: %v\n", err)
	}
	t.Log(mails)
}

func TestReadMail(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	mailRep := PostgresMailRepository{
		dbInstance,
	}

	err = mailRep.ReadMail(other.Username, owner.Username)
	if err != nil {
		t.Errorf("Didn't read emails: %v\n", err)
	}
}

func TestGetDialoguesForUser(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	mailRep := PostgresMailRepository{
		dbInstance,
	}

	dialogues, err = mailRep.GetDialoguesForUser(owner.Username, 10, "", 0, "@liokor.ru")
	if err != nil {
		t.Errorf("Didn't pass valid get dialogues: %v\n", err)
	}

	t.Log(dialogues)
}

func TestCountMailsFromUser(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	mailRep := PostgresMailRepository{
		dbInstance,
	}

	_, err = mailRep.CountMailsFromUser(owner.Username, time.Minute*(-10))
	if err != nil {
		t.Errorf("Error: %v\n", err)
	}
}

func TestReadDialogue(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	mailRep := PostgresMailRepository{
		dbInstance,
	}

	err = mailRep.ReadDialogue(other.Username, owner.Username)
	if err != nil {
		t.Errorf("Error: %v\n", err)
	}
}

func TestCreateFolder(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	mailRep := PostgresMailRepository{
		dbInstance,
	}

	folder, err = mailRep.CreateFolder(owner.Id, common.GenerateRandomString()[:8])

	if err != nil {
		t.Errorf("Didn't pass valid creating folder: %v\n", err)
	}

	_, err = mailRep.CreateFolder(0, "InvalidFolder")
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Created folder for non existing user: %v\n", err)
	}
}

func TestGetFolders(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	mailRep := PostgresMailRepository{
		dbInstance,
	}

	folders, err := mailRep.GetFolders(owner.Id)
	if err != nil {
		t.Errorf("Error: %v\n", err)
	}
	t.Log(folders)
}

func TestAddDialogueToFolder(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	mailRep := PostgresMailRepository{
		dbInstance,
	}

	err = mailRep.AddDialogueToFolder(owner.Username, folder.Id, dialogues[0].Id)
	if err != nil {
		t.Errorf("Didn't pass valid adding dialogues to folder: %v\n", err)
	}

	err = mailRep.AddDialogueToFolder(owner.Username, 0, dialogues[0].Id)
	if err != nil {
		t.Errorf("Didn't pass valid adding dialogues to global folder: %v\n", err)
	}

	err = mailRep.AddDialogueToFolder(owner.Username, -1, dialogues[0].Id)
	switch err.(type) {
	case mail.InvalidEmailError:
		break
	default:
		t.Errorf("Added dialogue to non existing folder: %v\n", err)
	}
}

func TestRemoveUsers(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Fatalf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	ur := userRep.PostgresUserRepository{
		dbInstance,
	}
	err = ur.RemoveUser(owner.Username)
	if err != nil {
		t.Fatalf("Database error: %v\n", err)
	}
	err = ur.RemoveUser(other.Username)
	if err != nil {
		t.Fatalf("Database error: %v\n", err)
	}
}
