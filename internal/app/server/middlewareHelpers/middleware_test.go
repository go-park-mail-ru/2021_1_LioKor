package middlewareHelpers

import (
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/timestamppb"
	"liokor_mail/internal/pkg/common"
	session "liokor_mail/internal/pkg/common/protobuf_sessions"
	sMocks "liokor_mail/internal/pkg/common/protobuf_sessions/mocks"
	"liokor_mail/internal/pkg/user"
	"liokor_mail/internal/pkg/user/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIsAuthMiddleware(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUC := mocks.NewMockUseCase(mockCtrl)
	mockSession := sMocks.NewMockIsAuthClient(mockCtrl)
	authMdlwr := AuthMiddleware{
		mockUC,
		mockSession,
	}

	s := session.Session{
		UserId:       int32(1),
		SessionToken: "sessionToken",
		Expiration:   timestamppb.New(time.Now().Add(10 * 24 * time.Hour)),
	}

	retUser := user.User{
		Id:           1,
		Username:     "test",
		HashPassword: "hash",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test test",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}

	gomock.InOrder(
		mockSession.EXPECT().
			Get(gomock.Any(), &session.SessionToken{SessionToken: "sessionToken"}).
			Return(&s, nil).
			Times(1),
		mockUC.EXPECT().GetUserById(1).Return(retUser, nil).Times(1),
	)
	e := echo.New()
	url := "/user"
	req := httptest.NewRequest("POST", url, nil)
	req.Header.Add("Cookie", "session_token=sessionToken; Expires=Wed, 03 Jun 2021 03:30:48 GMT; HttpOnly")
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

	_, err := authMdlwr.isAuthenticated(&echoContext)
	if err != nil {
		t.Errorf("Didn't pass valid data:%v\n", err)
	}

	mockSession.EXPECT().
		Get(gomock.Any(), &session.SessionToken{SessionToken: "sessionToken"}).
		Return(&session.Session{}, errors.New("Some error")).
		Times(1)

	_, err = authMdlwr.isAuthenticated(&echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusUnauthorized {
			t.Errorf("Didn't pass invalid credentails: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass invalid credentails: %v\n", err)
	}

	req = httptest.NewRequest("POST", url, nil)
	req.Header.Add("Cookie", "no_sessionToken")
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	_, err = authMdlwr.isAuthenticated(&echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusUnauthorized {
			t.Errorf("Didn't pass invalid credentails: %v\n", err)
		}
	} else {
		t.Errorf("Didn't pass invalid credentails: %v\n", err)
	}
}

