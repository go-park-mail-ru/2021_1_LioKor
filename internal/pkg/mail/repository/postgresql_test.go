package repository

import (
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/mail"
	"testing"
	"time"
)

var dbConfig = "host=localhost user=postgres password=12 dbname=liokor_mail_test sslmode=disable"

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
		Sender:    "alt@liokor.ru",
		Recipient: "lio@liokor.ru",
		Subject:   "Test mail",
		Body:      "Test adding mails",
	}

	err = mailRep.AddMail(newMail)
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

	mails, err := mailRep.GetMailsForUser("alt@liokor.ru", "lio@liokor.ru", 10, 0)
	if err != nil {
		t.Errorf("Error: %v\n", err)
	}
	for _, mail := range mails {
		t.Log(mail)
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

	dialogues, err := mailRep.GetDialoguesForUser("alt@liokor.ru", 10, 0)
	if err != nil {
		t.Errorf("Error: %v\n", err)
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

	count, err := mailRep.CountMailsFromUser("alt@liokor.ru", time.Now().Add(time.Minute*(-10)))
	if err != nil {
		t.Errorf("Error: %v\n", err)
	}

	t.Log(count)
}
