package server

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"liokor_mail/internal/pkg/user"
	"liokor_mail/internal/pkg/user/delivery"
	"liokor_mail/internal/pkg/user/repository"
	"liokor_mail/internal/pkg/user/usecase"
	"sync"

	"net/http"
	"net/http/httptest"
	"testing"
)

var userHandler = delivery.UserHandler{
	&usecase.UserUseCase{
		&repository.UserRepository{
			repository.UserStruct{map[string]user.User{}, sync.Mutex{}},
			repository.SessionStruct{map[string]user.Session{}, sync.Mutex{}},
		},
	},
}

func TestSignUp(t *testing.T) {
	e := echo.New()
	testUser := user.UserSignUp{
		"TEST",
		"pswdPSWD12",
		"some url",
		"test Testing",
		"someemail@mail.ru",
	}
	body, _ := json.Marshal(testUser)
	url := "/user"
	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	userHandler.SignUp(echoContext)

	if response.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", response.Code, http.StatusOK)
	}

	// sending same user again
	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	response = httptest.NewRecorder()

	echoContext = e.NewContext(req, response)
	err := userHandler.SignUp(echoContext)
	if err == nil {
		t.Error("Expected error, but it's nil")
	}

	testUser2 := user.UserSignUp{
		"test2",             // username
		"pswdPSWD12",              // password
		"http://wolf.wolf",  // avatar url
		"test Testing",      // fullname
		"someemail@mail.ru", // email
	}

	body, _ = json.Marshal(testUser2)
	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	response = httptest.NewRecorder()

	echoContext = e.NewContext(req, response)
	userHandler.SignUp(echoContext)

	if response.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", response.Code, http.StatusOK)
	}

	testUser3 := user.UserSignUp{
		"test",              // username
		"pswdPSWD12",              // password
		"http://wolf.wolf",  // avatar url
		"test Testing",      // fullname
		"someemail@mail.ru", // email
	}

	body, _ = json.Marshal(testUser3)
	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	response = httptest.NewRecorder()

	echoContext = e.NewContext(req, response)
	err = userHandler.SignUp(echoContext)
	if err == nil {
		t.Error("Was able to create two users with the same username, but different letter case!")
	}
}

func TestAuthenticate(t *testing.T) {
	e := echo.New()
	creds := user.Credentials{
		"test",
		"pswdPSWD12",
	}

	body, _ := json.Marshal(creds)
	url := "/user/auth"
	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	userHandler.Auth(echoContext)

	if response.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", response.Code, http.StatusOK)
	}

	cookies := response.Result().Cookies()
	if len(cookies) == 0 {
		t.Error("No cookie found after authentication!")
	}

	creds2 := user.Credentials{
		"te",
		"pswdPSWD12",
	}

	body, _ = json.Marshal(creds2)
	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)

	err := userHandler.Auth(echoContext)
	if err == nil {
		t.Error("Expected error, but it's nil")
	}
}

func TestCookie(t *testing.T) {
	e := echo.New()
	url := "/user"
	req := httptest.NewRequest("GET", url, nil)

	session, err := userHandler.UserUsecase.CreateSession("test")
	if err != nil {
		t.Errorf("Unable to create session!")
		return
	}
	req.Header.Add("Cookie", "session_token="+session.Value+"; Expires=Wed, 03 Mar 2021 03:30:48 GMT; HttpOnly")

	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	userHandler.Profile(echoContext)

	if response.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", response.Code, http.StatusOK)
	}

	sessionUser := user.User{}
	b := response.Result().Body
	err = json.NewDecoder(b).Decode(&sessionUser)
	if err != nil {
		t.Errorf("Json error")
		return
	}

	expectedUser := user.User{
		"TEST",
		"",
		"some url",
		"test Testing",
		"someemail@mail.ru",
		"",
		false,
	}
	assert.Equal(t, expectedUser.Username, sessionUser.Username)
	assert.Equal(t, expectedUser.HashPassword, sessionUser.HashPassword)
	assert.Equal(t, expectedUser.AvatarURL, sessionUser.AvatarURL)
	assert.Equal(t, expectedUser.FullName, sessionUser.FullName)
	assert.Equal(t, expectedUser.ReserveEmail, sessionUser.ReserveEmail)
}

func TestUpdate(t *testing.T) {
	e := echo.New()
	updUser := user.User{
		"TEST",
		"",
		"new url",
		"test2 test2",
		"someemail@mail.ru",
		"",
		false,
	}

	body, _ := json.Marshal(updUser)
	url := "/user/TEST"
	req := httptest.NewRequest("PUT", url, bytes.NewReader(body))

	session, err := userHandler.UserUsecase.CreateSession("test")
	if err != nil {
		t.Errorf("Unable to create session!")
		return
	}
	req.Header.Add("Cookie", "session_token="+session.Value+"; Expires=Wed, 03 Mar 2021 03:30:48 GMT; HttpOnly")

	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)
	echoContext.SetPath("/:username")
	echoContext.SetParamNames("username")
	echoContext.SetParamValues("TEST")

	userHandler.UpdateProfile(echoContext)

	if response.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", response.Code, http.StatusOK)
	}
	sessionUser := user.User{}
	b := response.Result().Body
	err = json.NewDecoder(b).Decode(&sessionUser)
	if err != nil {
		t.Errorf("Json error")
		return
	}

	expectedUser := user.User{
		"TEST",
		"",
		"new url",
		"test2 test2",
		"someemail@mail.ru",
		"",
		false,
	}
	assert.Equal(t, expectedUser.Username, sessionUser.Username)
	assert.Equal(t, expectedUser.AvatarURL, sessionUser.AvatarURL)
	assert.Equal(t, expectedUser.FullName, sessionUser.FullName)
	assert.Equal(t, expectedUser.ReserveEmail, sessionUser.ReserveEmail)
}

func TestChangePassword(t *testing.T) {
	e := echo.New()
	chPSWD := user.ChangePassword{
		"pswdPSWD12",
		"newPSWD12",
	}
	body, _ := json.Marshal(chPSWD)
	url := "/user/test/password"
	req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=test; Expires=Wed, 03 Mar 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)
	echoContext.SetPath("/:username/password")
	echoContext.SetParamNames("username")
	echoContext.SetParamValues("test")
	userHandler.ChangePassword(echoContext)

	if response.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", response.Code, http.StatusOK)
	}

}

func TestLogout(t *testing.T) {
	e := echo.New()
	url := "/user/logout"
	req := httptest.NewRequest("POST", url, nil)
	req.Header.Add("Cookie", "session_token=test; Expires=Wed, 03 Mar 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)
	userHandler.Logout(echoContext)

	if response.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", response.Code, http.StatusOK)
	}

	creds := user.Credentials{
		"test",
		"newPSWD12",
	}

	body, _ := json.Marshal(creds)
	url = "/user/auth"
	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)

	userHandler.Auth(echoContext)

	if response.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", response.Code, http.StatusOK)
	}
}
