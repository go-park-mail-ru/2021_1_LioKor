package main

import (
	"bytes"
	"encoding/json"
	"github.com/go-park-mail-ru/2021_1_LioKor/internal/pkg/user"
	"github.com/go-park-mail-ru/2021_1_LioKor/internal/pkg/user/delivery"
	"github.com/go-park-mail-ru/2021_1_LioKor/internal/pkg/user/repository"
	"github.com/go-park-mail-ru/2021_1_LioKor/internal/pkg/user/usecase"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var h = delivery.UserHandler{
	usecase.UserUseCase{
repository.UserRepository{
		map[string]user.User{},
		map[string]user.Session{},
		},
	},
}

func TestSignUp(t *testing.T) {
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

	h.UserPage(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", w.Code, http.StatusOK)
	}

	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	w = httptest.NewRecorder()

	h.UserPage(w, req)
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

	h.UserPage(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Wrong status code: %d, expected: %d", w.Code, http.StatusOK)
	}
}

func TestAuthenticate(t *testing.T) {
	creds := user.Credentials{
		"test",
		"pswd",
	}

	body, _ := json.Marshal(creds)
	url := "/user/auth"
	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	w := httptest.NewRecorder()

	h.Authenticate(w, req)

	if w.Code != http.StatusFound {
		t.Errorf("Wrong status code: %d, expected: %d", w.Code, http.StatusFound)
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

	h.Authenticate(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Wrong status code: %d, expected: %d", w.Code, http.StatusUnauthorized)
	}

}

func TestCookie(t *testing.T) {
	url := "/user"
	req := httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=test; Expires=Wed, 03 Mar 2021 03:30:48 GMT; HttpOnly")
	w := httptest.NewRecorder()

	h.UserPage(w, req)

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
	url := "/user"
	req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=test; Expires=Wed, 03 Mar 2021 03:30:48 GMT; HttpOnly")
	w := httptest.NewRecorder()

	h.UserPage(w, req)
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