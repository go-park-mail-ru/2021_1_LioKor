package middlewareHelpers

import (
	"context"
	"github.com/labstack/echo/v4"
	"liokor_mail/internal/pkg/common"
	session "liokor_mail/internal/pkg/common/protobuf_sessions"
	"liokor_mail/internal/pkg/user"
	"net/http"
	"time"
)

type AuthMiddleware struct {
	UserUsecase    user.UseCase
	SessionManager session.IsAuthClient
}

func (m *AuthMiddleware) IsAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sessionUser, err := m.isAuthenticated(&c)
		if err != nil {
			return err
		}
		c.Set("sessionUser", sessionUser)
		return next(c)
	}
}

func (m *AuthMiddleware) isAuthenticated(c *echo.Context) (user.User, error) {
	sessionToken, err := (*c).Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return user.User{}, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return user.User{}, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	s, err := m.SessionManager.Get(
		context.Background(),
		&session.SessionToken{SessionToken: sessionToken.Value},
	)
	if err != nil {
		m.deleteSessionCookie(c)
		return user.User{}, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	sessionUser, err := m.UserUsecase.GetUserById(int(s.UserId))

	if err != nil {
		switch err.(type) {
		case common.InvalidUserError:
			return user.User{}, echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return user.User{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return sessionUser, nil
}

func (m *AuthMiddleware) deleteSessionCookie(c *echo.Context) {
	(*c).SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Time{},
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
	})
}
