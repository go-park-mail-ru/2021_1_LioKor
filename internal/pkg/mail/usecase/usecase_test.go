package usecase

import (
	"database/sql"
	"github.com/golang/mock/gomock"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/mail"
	"liokor_mail/internal/pkg/mail/mocks"
	"testing"
	"time"
)

var config = common.Config{
	MailDomain: "liokor.ru",
}

func TestGetDialogues(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockMailRepository(mockCtrl)
	mailUC := MailUseCase{
		Repository: mockRep,
		Config: config,
	}

	dialogues := []mail.Dialogue{
		{
			Id: 1,
			Email: "lio@liokor.ru",
			AvatarURL: common.NullString{sql.NullString{String: "/media/test",Valid: true}},
			Body: "Test",
			Received_date: time.Now(),
		},
		{
			Id: 2,
			Email: "ser@liokor.ru",
			AvatarURL: common.NullString{sql.NullString{String: "",Valid: false}},
			Body: "Test",
			Received_date: time.Now(),
		},
	}

	mockRep.
		EXPECT().
		GetDialoguesForUser("alt@liokor.ru", 10, 0, "", "@liokor.ru").
		Return(dialogues, nil).
		Times(1)
	_, err := mailUC.GetDialogues("alt", 0, 10, "")
	if err != nil {
		t.Errorf("Didn't pass valid data: %v\n", err)
	}

	mockRep.
		EXPECT().
		GetDialoguesForUser("alt@liokor.ru", 10, 0, "", "@liokor.ru").
		Return(nil, mail.InvalidEmailError{
			"Error",
	}).
		Times(1)
	_, err = mailUC.GetDialogues("alt", 0, 10, "")
	switch err.(type) {
	case mail.InvalidEmailError:
		break
	default:
		t.Errorf("Didn't pass invalid data: %v\n", err)
	}
}

func TestGetEmails(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockMailRepository(mockCtrl)
	mailUC := MailUseCase{
		Repository: mockRep,
		Config: config,
	}

	emails := []mail.DialogueEmail{
		{
			Id : 1,
			Sender: "alt@liokor.ru",
			Subject: "Test",
			Received_date: time.Now(),
			Body: "Test",
		},
		{
			Id : 2,
			Sender: "lio@liokor.ru",
			Subject: "Test",
			Received_date: time.Now(),
			Body: "Test",
		},
	}

	mockRep.
		EXPECT().
		GetMailsForUser("alt@liokor.ru", "lio@liokor.ru", 10, 0).
		Return(emails, nil).
		Times(1)
	_, err := mailUC.GetEmails("alt", "lio@liokor.ru", 0, 10)
	if err != nil {
		t.Errorf("Didn't pass valid data: %v\n", err)
	}

	mockRep.
		EXPECT().
		GetMailsForUser("alt@liokor.ru", "lio@liokor.ru", 10, 0).
		Return(nil, mail.InvalidEmailError{
			"Error",
	}).
		Times(1)
	_, err = mailUC.GetEmails("alt", "lio@liokor.ru", 0, 10)
	switch err.(type) {
	case mail.InvalidEmailError:
		break
	default:
		t.Errorf("Didn't pass invalid data: %v\n", err)
	}
}

func TestSendEmail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockMailRepository(mockCtrl)
	mailUC := MailUseCase{
		Repository: mockRep,
		Config: config,
	}

	email := mail.Mail{
		Sender: "alt",
		Recipient: "altana@liokor.ru",
		Body: "Testing",
		Subject: "Test",
	}
	emailSent := mail.Mail{
		Sender: "alt@liokor.ru",
		Recipient: "altana@liokor.ru",
		Body: "Testing",
		Subject: "Test",
	}
	gomock.InOrder(
		mockRep.EXPECT().CountMailsFromUser("alt@liokor.ru", 3*time.Minute).Return(0, nil).Times(1),
		mockRep.EXPECT().AddMail(emailSent).Return(nil).Times(1),
		)
	err := mailUC.SendEmail(email)
	if err != nil {
		t.Errorf("Couldn't send email: %v\n", err)
	}

	mockRep.EXPECT().CountMailsFromUser("alt@liokor.ru", 3*time.Minute).Return(6, nil).Times(1)
	err = mailUC.SendEmail(email)
	switch err.(type) {
	case mail.InvalidEmailError:
		break
	default:
		t.Errorf("Didn't pass invalid data: %v\n", err)
	}

	gomock.InOrder(
		mockRep.EXPECT().CountMailsFromUser("alt@liokor.ru", 3*time.Minute).Return(0, nil).Times(1),
		mockRep.EXPECT().AddMail(emailSent).Return(mail.InvalidEmailError{"Error"}).Times(1),
	)
	err = mailUC.SendEmail(email)
	switch err.(type) {
	case mail.InvalidEmailError:
		break
	default:
		t.Errorf("Didn't pass invalid data: %v\n", err)
	}
}
