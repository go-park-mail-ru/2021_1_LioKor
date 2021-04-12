package repository

import (
	"liokor_mail/internal/pkg/common"
	"testing"
)

var dbConfig = "host=localhost user=postgres password=12 dbname=liokor_mail_test sslmode=disable"

func TestGetMailsForUser(t *testing.T) {
	dbInstance, err := common.NewPostgresDataBase(dbConfig)
	if err != nil {
		t.Errorf("Database error: %v\n", err)
	}
	defer dbInstance.Close()

	mailRep := PostgresMailRepository{
		dbInstance,
	}

	mails, err := mailRep.GetMailsForUser("alt@alt.alt", "lio@kor.ru", 10, 0)
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

	dialogues, err := mailRep.GetDialoguesForUser("ser@liokor.ru", 10, 0)
	if err != nil {
		t.Errorf("Error: %v\n", err)
	}

	t.Log(dialogues)
}
