package delivery

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/mail"
	mailMocks "liokor_mail/internal/pkg/mail/mocks"
	"liokor_mail/internal/pkg/user"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetDialogues(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMailUC := mailMocks.NewMockMailUseCase(mockCtrl)
	mailHandler := MailHandler{
		mockMailUC,
	}

	e := echo.New()

	url := "/email/dialogues/?amount=5&find=a&folder=1"
	req := httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 Jun 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	sessionUser := user.User{
		Username:     "sessionTest",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{Valid: true, String: "/media/test"}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	dialogues := []mail.Dialogue{
		{
			Id:            1,
			Email:         "alt@liokor.ru",
			AvatarURL:     common.NullString{sql.NullString{String: "/media/test", Valid: true}},
			Body:          "Test",
			Received_date: time.Now(),
			Unread: 0,
		},
		{
			Id:            2,
			Email:         "aser@liokor.ru",
			AvatarURL:     common.NullString{sql.NullString{String: "", Valid: false}},
			Body:          "Test",
			Received_date: time.Now(),
			Unread: 1,
		},
	}
	echoContext.Set("sessionUser", sessionUser)

	mockMailUC.EXPECT().GetDialogues(sessionUser.Username, 5, "a", 1).Return(dialogues, nil).Times(1)
	err := mailHandler.GetDialogues(echoContext)
	if err != nil {
		t.Errorf("Didn't pass valid data: %v\n", err)
	}

	dEmails := make([]mail.Dialogue, 0, 0)
	err = json.Unmarshal(response.Body.Bytes(), &dEmails)
	if err != nil {
		t.Errorf("Json error: %v\n", err.Error())
	}

	req = httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 Jun 2021 03:30:48 GMT; HttpOnly")
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	echoContext.Set("sessionUser", sessionUser)

	mockMailUC.EXPECT().GetDialogues(sessionUser.Username, 5, "a", 1).Return(nil, mail.InvalidEmailError{"Error"}).Times(1)
	err = mailHandler.GetDialogues(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusInternalServerError {
			t.Errorf("Didn't pass invalid session token: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass invalid session token: %v\n", err)
	}
}

func TestGetEmails(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMailUC := mailMocks.NewMockMailUseCase(mockCtrl)

	mailHandler := MailHandler{
		mockMailUC,
	}

	e := echo.New()

	url := "/email/emails/?with=lio@liokor.ru&last=1&amount=5"
	req := httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 Jun 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	sessionUser := user.User{
		Username:     "alt",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	emails := []mail.DialogueEmail{
		{
			Id:            1,
			Sender:        "lio@liokor.ru",
			Subject:       "Test",
			Received_date: time.Now(),
			Body:          "Test",
			Unread: false,
			Status: 1,
		},
		{
			Id:            2,
			Sender:        "lio@liokor.ru",
			Subject:       "Test",
			Received_date: time.Now(),
			Body:          "Test",
			Unread: true,
			Status: 1,
		},
	}
	echoContext.Set("sessionUser", sessionUser)

	mockMailUC.EXPECT().GetEmails(sessionUser.Username, "lio@liokor.ru", 1, 5).Return(emails, nil).Times(1)
	err := mailHandler.GetEmails(echoContext)
	if err != nil {
		t.Errorf("Didn't pass valid data: %v\n", err)
	}

	url = "/email/emails/?with=lio@liokor.ru&amount=100"
	req = httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 Jun 2021 03:30:48 GMT; HttpOnly")
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	echoContext.Set("sessionUser", sessionUser)

	mockMailUC.EXPECT().GetEmails(sessionUser.Username, "lio@liokor.ru", 0, 50).Return(nil, mail.InvalidEmailError{"error"}).Times(1)
	err = mailHandler.GetEmails(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusBadRequest {
			t.Errorf("Didn't pass invalid data: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass invalid data: %v\n", err)
	}
}

func TestSendEmail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMailUC := mailMocks.NewMockMailUseCase(mockCtrl)

	mailHandler := MailHandler{
		mockMailUC,
	}

	e := echo.New()

	url := "/email"

	email := mail.Mail{
		Recipient: "altana@liokor.ru",
		Body:      "Testing",
		Subject:   "Test",
	}
	body, _ := json.Marshal(email)
	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 Jun 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	sessionUser := user.User{
		Username:     "alt",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	emailSent := mail.Mail{
		Sender:    sessionUser.Username,
		Recipient: "altana@liokor.ru",
		Body:      "Testing",
		Subject:   "Test",
	}
	echoContext.Set("sessionUser", sessionUser)

	mockMailUC.EXPECT().SendEmail(emailSent).Return(nil).Times(1)
	err := mailHandler.SendEmail(echoContext)
	if err != nil {
		t.Errorf("Didn't pass valid data: %v\n", err)
	}

	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 Jun 2021 03:30:48 GMT; HttpOnly")
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	echoContext.Set("sessionUser", sessionUser)

	mockMailUC.EXPECT().SendEmail(emailSent).Return(mail.InvalidEmailError{"error"}).Times(1)
	err = mailHandler.SendEmail(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusInternalServerError {
			t.Errorf("Didn't pass invalid data: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass invalid data: %v\n", err)
	}
}

func TestGetFolders(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMailUC := mailMocks.NewMockMailUseCase(mockCtrl)

	mailHandler := MailHandler{
		mockMailUC,
	}

	e := echo.New()

	url := "/email/folders"
	req := httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 Jun 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	sessionUser := user.User{
		Id : 1,
		Username:     "alt",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	folders := []mail.Folder{
		{
			Id : 1,
			FolderName: "NewFolder",
			Owner: 1,
		},
		{
			Id : 2,
			FolderName: "AnotherFolder",
			Owner: 1,
		},
	}
	echoContext.Set("sessionUser", sessionUser)

	mockMailUC.EXPECT().GetFolders(sessionUser.Id).Return(folders, nil).Times(1)
	err := mailHandler.GetFolders(echoContext)
	if err != nil {
		t.Errorf("Didn't get valid folders: %v\n", err.Error())
	}
}

func TestCreateFolder(t *testing.T){
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMailUC := mailMocks.NewMockMailUseCase(mockCtrl)

	mailHandler := MailHandler{
		mockMailUC,
	}

	e := echo.New()
	folderName := struct{
		FolderName string `json:"folderName"` } {
		FolderName: "NewFolderName",
	}
	body, _ := json.Marshal(folderName)
	url := "/email/folders"
	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 Jun 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	sessionUser := user.User{
		Id : 1,
		Username:     "alt",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	folder := mail.Folder{
		Id : 1,
		FolderName: folderName.FolderName,
		Owner: 1,
	}
	echoContext.Set("sessionUser", sessionUser)

	mockMailUC.EXPECT().CreateFolder(sessionUser.Id, folderName.FolderName).Return(folder, nil).Times(1)
	err := mailHandler.CreateFolder(echoContext)
	if err != nil {
		t.Errorf("Didn't get valid folders: %v\n", err.Error())
	}

	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 Jun 2021 03:30:48 GMT; HttpOnly")
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	echoContext.Set("sessionUser", sessionUser)
	mockMailUC.EXPECT().CreateFolder(sessionUser.Id, folderName.FolderName).Return(mail.Folder{}, common.InvalidUserError{"User doesn't exist"}).Times(1)
	err = mailHandler.CreateFolder(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusInternalServerError {
			t.Errorf("Didn't pass invalid data: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass invalid data: %v\n", err)
	}
}

func TestUpdateFolder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockMailUC := mailMocks.NewMockMailUseCase(mockCtrl)

	mailHandler := MailHandler{
		mockMailUC,
	}

	e := echo.New()
	updateFolder := struct {
		FolderId int `json:"folderId"`
		DialogueId int `json:"dialogueId"`
	} {
		FolderId: 1,
		DialogueId: 1,
	}
	body, _ := json.Marshal(updateFolder)
	url := "/email/folders"
	req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 Jun 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	sessionUser := user.User{
		Id : 1,
		Username:     "alt",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}
	echoContext.Set("sessionUser", sessionUser)

	mockMailUC.EXPECT().UpdateFolder(sessionUser.Username, updateFolder.FolderId, updateFolder.DialogueId).Return(nil).Times(1)
	err := mailHandler.UpdateFolder(echoContext)
	if err != nil {
		t.Errorf("Didn't add valid dialogue to folder: %v\n", err.Error())
	}

	req = httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 Jun 2021 03:30:48 GMT; HttpOnly")
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	echoContext.Set("sessionUser", sessionUser)
	mockMailUC.EXPECT().UpdateFolder(sessionUser.Username, updateFolder.FolderId, updateFolder.DialogueId).Return(mail.InvalidEmailError{"Folder doesn't exist"}).Times(1)
	err = mailHandler.UpdateFolder(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusInternalServerError {
			t.Errorf("Didn't pass invalid data: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass invalid data: %v\n", err)
	}
}