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
		Config:     config,
	}

	dialogues := []mail.Dialogue{
		{
			Id:            1,
			Email:         "lio@liokor.ru",
			AvatarURL:     common.NullString{sql.NullString{String: "/media/test", Valid: true}},
			Body:          "Test",
			Received_date: time.Now(),
			Unread:        0,
		},
		{
			Id:            2,
			Email:         "ser@liokor.ru",
			AvatarURL:     common.NullString{sql.NullString{String: "", Valid: false}},
			Body:          "Test",
			Received_date: time.Now(),
			Unread:        1,
		},
	}

	mockRep.
		EXPECT().
		GetDialoguesForUser("alt@liokor.ru", 10, "", 0, "@liokor.ru").
		Return(dialogues, nil).
		Times(1)
	_, err := mailUC.GetDialogues("alt", 10, "", 0)
	if err != nil {
		t.Errorf("Didn't pass valid data: %v\n", err)
	}

	mockRep.
		EXPECT().
		GetDialoguesForUser("alt@liokor.ru", 10, "", 0, "@liokor.ru").
		Return(nil, mail.InvalidEmailError{
			"Error",
		}).
		Times(1)
	_, err = mailUC.GetDialogues("alt", 10, "", 0)
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
		Config:     config,
	}

	emails := []mail.DialogueEmail{
		{
			Id:            1,
			Sender:        "alt@liokor.ru",
			Subject:       "Test",
			Received_date: time.Now(),
			Body:          "Test",
			Unread:        false,
			Status:        1,
		},
		{
			Id:            2,
			Sender:        "lio@liokor.ru",
			Subject:       "Test",
			Received_date: time.Now(),
			Body:          "Test",
			Unread:        true,
			Status:        1,
		},
	}

	gomock.InOrder(
		mockRep.
			EXPECT().
			GetMailsForUser("alt@liokor.ru", "lio@liokor.ru", 10, 0).
			Return(emails, nil).
			Times(1),
		mockRep.
			EXPECT().
			ReadMail("alt@liokor.ru", "lio@liokor.ru").
			Return(nil).
			Times(1),
		mockRep.
			EXPECT().
			ReadDialogue("alt@liokor.ru", "lio@liokor.ru").
			Return(nil).
			Times(1),
	)
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
		Config:     config,
	}

	email := mail.Mail{
		Sender:    "alt",
		Recipient: "altana@liokor.ru",
		Body:      "Testing",
		Subject:   "Test",
	}
	emailSent := mail.Mail{
		Sender:    "alt@liokor.ru",
		Recipient: "altana@liokor.ru",
		Body:      "Testing",
		Subject:   "Test",
	}
	gomock.InOrder(
		mockRep.EXPECT().CountMailsFromUser("alt@liokor.ru", 3*time.Minute).Return(0, nil).Times(1),
		mockRep.EXPECT().AddMail(emailSent).Return(1, nil).Times(1),
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
		mockRep.EXPECT().AddMail(emailSent).Return(0, mail.InvalidEmailError{"Error"}).Times(1),
	)
	err = mailUC.SendEmail(email)
	switch err.(type) {
	case mail.InvalidEmailError:
		break
	default:
		t.Errorf("Didn't pass invalid data: %v\n", err)
	}
}

func TestGetFolders(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockMailRepository(mockCtrl)
	mailUC := MailUseCase{
		Repository: mockRep,
		Config:     config,
	}

	folders := []mail.Folder{
		{
			Id:         1,
			FolderName: "NewFolder",
			Owner:      1,
		},
		{
			Id:         2,
			FolderName: "AnotherFolder",
			Owner:      1,
		},
	}
	mockRep.EXPECT().GetFolders(1).Return(folders, nil).Times(1)
	_, err := mailUC.GetFolders(1)
	if err != nil {
		t.Errorf("Didn't get valid folders: %v\n", err)
	}
}

func TestCreateFolder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockMailRepository(mockCtrl)
	mailUC := MailUseCase{
		Repository: mockRep,
		Config:     config,
	}

	folder := mail.Folder{
		Id:         1,
		FolderName: "NewFolder",
		Owner:      1,
	}
	mockRep.EXPECT().CreateFolder(folder.Owner, folder.FolderName).Return(folder, nil).Times(1)
	_, err := mailUC.CreateFolder(folder.Owner, folder.FolderName)
	if err != nil {
		t.Errorf("Didn't create valid folders: %v\n", err)
	}

	mockRep.EXPECT().CreateFolder(folder.Owner, folder.FolderName).Return(mail.Folder{}, common.InvalidUserError{"User doesn't exist"}).Times(1)
	_, err = mailUC.CreateFolder(folder.Owner, folder.FolderName)
	switch err.(type) {
	case common.InvalidUserError:
		break
	default:
		t.Errorf("Didn't pass invalid data: %v\n", err)
	}
}

func TestUpdateFolder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockMailRepository(mockCtrl)
	mailUC := MailUseCase{
		Repository: mockRep,
		Config:     config,
	}

	mockRep.EXPECT().AddDialogueToFolder("alt@liokor.ru", 1, 1).Return(nil).Times(1)
	err := mailUC.UpdateFolder("alt", 1, 1)
	if err != nil {
		t.Errorf("Didn't update valid folders: %v\n", err)
	}

	mockRep.EXPECT().AddDialogueToFolder("alt@liokor.ru", 1, 1).Return(mail.InvalidEmailError{"Folder doesn't exist"}).Times(1)
	err = mailUC.UpdateFolder("alt", 1, 1)
	switch err.(type) {
	case mail.InvalidEmailError:
		break
	default:
		t.Errorf("Didn't pass invalid data: %v\n", err)
	}
}
