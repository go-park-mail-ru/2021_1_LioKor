package middlewareHelpers

import (
	"context"
	"github.com/labstack/echo/v4"
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
		sessionUserId, err := m.isAuthenticated(&c)
		if err != nil {
			return err
		}
		c.Set("sessionUserId", sessionUserId)
		return next(c)
	}
}

func (m *AuthMiddleware) isAuthenticated(c *echo.Context) (int, error) {
	sessionToken, err := (*c).Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return -1, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return -1, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	s, err := m.SessionManager.Get(
		context.Background(),
		&session.SessionToken{SessionToken: sessionToken.Value},
	)
	if err != nil {
		m.deleteSessionCookie(c)
		return -1, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}
	return int(s.UserId), nil
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
