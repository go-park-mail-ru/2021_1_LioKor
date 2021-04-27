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

	url := "/email/dialogues/?last=1&amount=5&find=a"
	req := httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
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
			Id: 1,
			Email: "lio@liokor.ru",
			AvatarURL: common.NullString{sql.NullString{Valid: true, String: "/media/test"}},
			Body: "Test",
			Received_date: time.Now(),
		},
		{
			Id: 2,
			Email: "ser@liokor.ru",
			AvatarURL: common.NullString{sql.NullString{Valid: false, String: ""}},
			Body: "Test",
			Received_date: time.Now(),
		},
	}
	echoContext.Set("sessionUser", sessionUser)

	mockMailUC.EXPECT().GetDialogues(sessionUser.Username, 1, 5, "a").Return(dialogues, nil).Times(1)
	err := mailHandler.GetDialogues(echoContext)
	if err != nil {
		t.Errorf("Didn't pass valid data: %v\n", err)
	}

	dEmails := make([]mail.Dialogue, 0, 0)
	err = json.Unmarshal(response.Body.Bytes(), &dEmails)
	if err == nil {
		t.Log(dEmails)
	}

	req = httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)

	err = mailHandler.GetDialogues(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusUnauthorized {
			t.Errorf("Didn't pass invalid session token: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass invalid session token: %v\n", err)
	}

	req = httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	echoContext.Set("sessionUser", sessionUser)

	mockMailUC.EXPECT().GetDialogues(sessionUser.Username, 1, 5, "a").Return(nil, mail.InvalidEmailError{"Error"}).Times(1)
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
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	sessionUser := user.User{
		Username:     "alt",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{ String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
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

	echoContext.Set("sessionUser", sessionUser)

	mockMailUC.EXPECT().GetEmails(sessionUser.Username, "lio@liokor.ru", 1, 5).Return(emails, nil).Times(1)
	err := mailHandler.GetEmails(echoContext)
	if err != nil {
		t.Errorf("Didn't pass valid data: %v\n", err)
	}

	url = "/email/emails/?with=lio@liokor.ru&amount=100"
	req = httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
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

	url = "/email/emails/?amount=100"
	req = httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	echoContext.Set("sessionUser", sessionUser)

	err = mailHandler.GetEmails(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusBadRequest {
			t.Errorf("Didn't pass invalid GET params: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass GET params: %v\n", err)
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
		Body: "Testing",
		Subject: "Test",
	}
	body, _ := json.Marshal(email)
	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
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
		Sender: sessionUser.Username,
		Recipient: "altana@liokor.ru",
		Body: "Testing",
		Subject: "Test",
	}
	echoContext.Set("sessionUser", sessionUser)

	mockMailUC.EXPECT().SendEmail(emailSent).Return(nil).Times(1)
	err := mailHandler.SendEmail(echoContext)
	if err != nil {
		t.Errorf("Didn't pass valid data: %v\n", err)
	}

	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
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