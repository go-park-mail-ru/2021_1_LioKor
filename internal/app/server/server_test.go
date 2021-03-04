package server

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"lioKor_mail/internal/pkg/user"
	"lioKor_mail/internal/pkg/user/delivery"
	"lioKor_mail/internal/pkg/user/repository"
	"lioKor_mail/internal/pkg/user/usecase"

	"net/http"
	"net/http/httptest"
	"testing"
)

var h = delivery.UserHandler{
	&usecase.UserUseCase{
		&repository.UserRepository{
				map[string]user.User{},
				map[string]user.Session{},
		},
	},
}

func TestSignUp(t *testing.T) {
	e := echo.New()
	 testUser:= user.UserSignUp{
		"test",
		"pswd",
		"some url",
		"test Testing",
		"someemail@mail.ru",
	}
	body, _ := json.Marshal(testUser)
	url := "/user"
	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)

	h.SignUp(c)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", w.Code, http.StatusOK)
	}

	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	w = httptest.NewRecorder()

	c = e.NewContext(req, w)
	h.SignUp(c)
	if w.Code != http.StatusConflict {
		t.Errorf("Wrong status code: %d, expected: %d", w.Code, http.StatusConflict)
	}

	testUser2:= user.UserSignUp{
		"test2",
		"pswd",
		"some url",
		"test Testing",
		"someemail@mail.ru",
	}

	body, _ = json.Marshal(testUser2)
	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	w = httptest.NewRecorder()

	c = e.NewContext(req, w)
	h.SignUp(c)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", w.Code, http.StatusOK)
	}
}


func TestAuthenticate(t *testing.T) {
	e := echo.New()
	creds := user.Credentials{
		"test",
		"pswd",
	}

	body, _ := json.Marshal(creds)
	url := "/user/auth"
	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)

	h.Auth(c)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", w.Code, http.StatusOK)
	}

	cookies := w.Result().Cookies()
	assert.Equal(t, "test", cookies[0].Value)


	creds2 := user.Credentials{
		"te",
		"pswd",
	}

	body, _ = json.Marshal(creds2)
	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	w = httptest.NewRecorder()
	c = e.NewContext(req, w)

	h.Auth(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Wrong status code: %d, expected: %d", w.Code, http.StatusUnauthorized)
	}

}

func TestCookie(t *testing.T) {
	e := echo.New()
	url := "/user"
	req := httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=test; Expires=Wed, 03 Mar 2021 03:30:48 GMT; HttpOnly")
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)

	h.Profile(c)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", w.Code, http.StatusOK)
	}

	sessionUser := user.User{}
	b := w.Result().Body
	err := json.NewDecoder(b).Decode(&sessionUser)
	if err != nil {
		t.Errorf("Json error")
	}

	expectedUser := user.User{
		"test",
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
	updUser:= user.User{
		"test",
		"",
		"new url",
		"test2 test2",
		"someemail@mail.ru",
		"",
		false,
	}

	body, _ := json.Marshal(updUser)
	url := "/user/test"
	req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=test; Expires=Wed, 03 Mar 2021 03:30:48 GMT; HttpOnly")
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)
	c.SetPath("/:username")
	c.SetParamNames("username")
	c.SetParamValues("test")

	h.UpdateProfile(c)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", w.Code, http.StatusOK)
	}
	sessionUser := user.User{}
	b := w.Result().Body
	err := json.NewDecoder(b).Decode(&sessionUser)
	if err != nil {
		t.Errorf("Json error")
	}

	expectedUser := user.User{
		"test",
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
		"pswd",
		"newPSWD",
	}
	body, _ := json.Marshal(chPSWD)
	url := "/user/test/password"
	req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=test; Expires=Wed, 03 Mar 2021 03:30:48 GMT; HttpOnly")
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)
	c.SetPath("/:username/password")
	c.SetParamNames("username")
	c.SetParamValues("test")
	h.ChangePassword(c)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", w.Code, http.StatusOK)
	}

}

func TestLogout(t *testing.T) {
	e := echo.New()
	url := "/user/logout"
	req := httptest.NewRequest("POST", url, nil)
	req.Header.Add("Cookie", "session_token=test; Expires=Wed, 03 Mar 2021 03:30:48 GMT; HttpOnly")
	w := httptest.NewRecorder()
	c := e.NewContext(req, w)
	h.Logout(c)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", w.Code, http.StatusOK)
	}


	creds := user.Credentials{
		"test",
		"newPSWD",
	}

	body, _ := json.Marshal(creds)
	url = "/user/auth"
	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	w = httptest.NewRecorder()
	c = e.NewContext(req, w)

	h.Auth(c)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", w.Code, http.StatusOK)
	}
}