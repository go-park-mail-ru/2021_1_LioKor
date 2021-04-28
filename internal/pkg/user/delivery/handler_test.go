package delivery

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/user"
	"liokor_mail/internal/pkg/user/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUC := mocks.NewMockUseCase(mockCtrl)
	userHandler := UserHandler{
		UserUsecase: mockUC,
	}

	e := echo.New()
	creds := user.Credentials{
		"test",
		"StrongPassword1",
	}

	body, _ := json.Marshal(creds)
	url := "/user/auth"
	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	retSession := user.SessionToken{
		Value:      "session token",
		Expiration: time.Now().Add(10 * 24 * time.Hour),
	}
	gomock.InOrder(
		mockUC.EXPECT().Login(creds).Return(nil).Times(1),
		mockUC.EXPECT().CreateSession(creds.Username).Return(retSession, nil).Times(1),
	)
	err := userHandler.Auth(echoContext)
	if err != nil {
		t.Errorf("Didn't pass valid credentails: %v\n", err)
	}

	cookies := response.Result().Cookies()
	if len(cookies) == 0 {
		t.Error("No cookie found after authentication!")
	}

	wrongCreds := user.Credentials{
		"te",
		"Strong",
	}

	body, _ = json.Marshal(wrongCreds)
	url = "/user/auth"
	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	mockUC.EXPECT().Login(wrongCreds).Return(user.InvalidUserError{"Invalid credentials"}).Times(1)
	err = userHandler.Auth(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusUnauthorized {
			t.Errorf("Didn't pass invalid credentails: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass invalid credentails: %v\n", err)
	}
}

func TestLogout(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUC := mocks.NewMockUseCase(mockCtrl)
	userHandler := UserHandler{
		UserUsecase: mockUC,
	}

	e := echo.New()
	url := "/user/logout"
	req := httptest.NewRequest("DELETE", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	mockUC.EXPECT().Logout("sessionToken").Return(nil).Times(1)
	err := userHandler.Logout(echoContext)
	if err != nil {
		t.Errorf("Didn't pass valid session token: %v\n", err)
	}

	req = httptest.NewRequest("DELETE", url, nil)
	req.Header.Add("Cookie", "")
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	err = userHandler.Logout(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusUnauthorized {
			t.Errorf("Didn't pass no cookie: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass no cookie: %v\n", err)
	}

	req = httptest.NewRequest("DELETE", url, nil)
	req.Header.Add("Cookie", "not_session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	err = userHandler.Logout(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusUnauthorized {
			t.Errorf("Didn't pass invalid cookie: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass invalid cookie: %v\n", err)
	}
}

func TestProfile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUC := mocks.NewMockUseCase(mockCtrl)
	userHandler := UserHandler{
		UserUsecase: mockUC,
	}

	e := echo.New()
	url := "/user"
	req := httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	retUser := user.User{
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}
	echoContext.Set("sessionUser", retUser)

	err := userHandler.Profile(echoContext)
	if err != nil {
		t.Errorf("Didn't pass valid session token: %v\n", err)
	}
}

func TestProfileByUsername(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUC := mocks.NewMockUseCase(mockCtrl)
	userHandler := UserHandler{
		UserUsecase: mockUC,
	}

	e := echo.New()
	url := "/user/test"
	req := httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)
	echoContext.SetPath("/:username")
	echoContext.SetParamNames("username")
	echoContext.SetParamValues("test")

	retUser := user.User{
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	mockUC.EXPECT().GetUserByUsername("test").Return(retUser, nil).Times(1)
	err := userHandler.ProfileByUsername(echoContext)
	if err != nil {
		t.Errorf("Didn't pass valid data: %v\n", err)
	}

	url = "/user/test"
	req = httptest.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	echoContext.SetPath("/:username")
	echoContext.SetParamNames("username")
	echoContext.SetParamValues("test")
	mockUC.EXPECT().GetUserByUsername("test").Return(user.User{}, user.InvalidUserError{"user doesn't exist"}).Times(1)
	err = userHandler.ProfileByUsername(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusNotFound {
			t.Errorf("Didn't pass invalid username: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass invalid username: %v\n", err)
	}
}

func TestSignUp(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUC := mocks.NewMockUseCase(mockCtrl)
	userHandler := UserHandler{
		UserUsecase: mockUC,
	}

	e := echo.New()
	newUser := user.UserSignUp{
		Username:     "test",
		Password:     "StrongPassword1",
		AvatarURL:    "/media/test",
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
	}
	body, _ := json.Marshal(newUser)
	url := "/user"
	req := httptest.NewRequest("POST", url, bytes.NewReader(body))
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	retSessionToken := user.SessionToken{
		Value:      "sessionToken",
		Expiration: time.Now().Add(10 * 24 * time.Hour),
	}

	gomock.InOrder(
		mockUC.EXPECT().SignUp(newUser).Return(nil).Times(1),
		mockUC.EXPECT().CreateSession(newUser.Username).Return(retSessionToken, nil).Times(1),
	)
	err := userHandler.SignUp(echoContext)
	if err != nil {
		t.Errorf("Didn't pass valid user: %v\n", err)
	}

	req = httptest.NewRequest("POST", url, bytes.NewReader(body))
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	mockUC.EXPECT().SignUp(newUser).Return(user.InvalidUserError{"username exists"}).Times(1)
	err = userHandler.SignUp(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusConflict {
			t.Errorf("Didn't pass invalid user: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass invalid user: %v\n", err)
	}
}

func TestUpdateProfile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUC := mocks.NewMockUseCase(mockCtrl)
	userHandler := UserHandler{
		UserUsecase: mockUC,
	}

	e := echo.New()
	newData := struct {
		AvatarURL    string `json:"avatarUrl"`
		FullName     string `json:"fullname"`
		ReserveEmail string `json:"reserveEmail"`
	}{
		AvatarURL:    "",
		FullName:     "New Full Name",
		ReserveEmail: "newtest@test.test",
	}
	body, _ := json.Marshal(newData)
	url := "/user/test"
	req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)
	echoContext.SetPath("/:username")
	echoContext.SetParamNames("username")
	echoContext.SetParamValues("test")

	retUser := user.User{
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}
	updUser := user.User{
		Username:     "",
		HashPassword: "",
		AvatarURL:    common.NullString{sql.NullString{String: "", Valid: true}},
		FullName:     "New Full Name",
		ReserveEmail: "newtest@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	echoContext.Set("sessionUser", retUser)

	mockUC.EXPECT().UpdateUser("test", updUser).Return(updUser, nil).Times(1)
	err := userHandler.UpdateProfile(echoContext)
	if err != nil {
		t.Errorf("Didn't pass valid upd data: %v\n", err)
	}

	req = httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	echoContext.SetPath("/:username")
	echoContext.SetParamNames("username")
	echoContext.SetParamValues("test2")
	echoContext.Set("sessionUser", retUser)

	err = userHandler.UpdateProfile(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusUnauthorized {
			t.Errorf("Didn't pass username not equal session user: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass username not equal session user: %v\n", err)
	}
}

func TestChangePassword(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUC := mocks.NewMockUseCase(mockCtrl)
	userHandler := UserHandler{
		UserUsecase: mockUC,
	}

	e := echo.New()
	newPSWD := user.ChangePassword{
		OldPassword: "StrongPassword1",
		NewPassword: "NewStrongPassword2",
	}
	body, _ := json.Marshal(newPSWD)
	url := "/user/test/password"
	req := httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)
	echoContext.SetPath("/:username/password")
	echoContext.SetParamNames("username")
	echoContext.SetParamValues("test")

	retUser := user.User{
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	echoContext.Set("sessionUser", retUser)

	mockUC.EXPECT().ChangePassword(retUser, newPSWD).Return(nil).Times(1)
	err := userHandler.ChangePassword(echoContext)
	if err != nil {
		t.Errorf("Didn't pass valid change password: %v\n", err)
	}

	req = httptest.NewRequest("PUT", url, bytes.NewReader(body))
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 May 2021 03:30:48 GMT; HttpOnly")
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	echoContext.SetPath("/:username/password")
	echoContext.SetParamNames("username")
	echoContext.SetParamValues("test2")
	echoContext.Set("sessionUser", retUser)

	err = userHandler.ChangePassword(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusUnauthorized {
			t.Errorf("Didn't pass username not equal session user: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass username not equal session user: %v\n", err)
	}
}
