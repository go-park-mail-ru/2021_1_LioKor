package common

import (
	"github.com/labstack/echo/v4"
	"liokor_mail/internal/pkg/user"
	"net/http"
	"time"
)

func IsAuthenticated(c *echo.Context, userUsecase user.UseCase) (user.User, error) {
	sessionToken, err := (*c).Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return user.User{}, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return user.User{}, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sessionUser, err := userUsecase.GetUserBySessionToken(sessionToken.Value)
	if err != nil {
		switch err.(type) {
		case user.InvalidSessionError, user.InvalidUserError:
			DeleteSessionCookie(c)
			return user.User{}, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		default:
			return user.User{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return sessionUser, nil
}

func DeleteSessionCookie(c *echo.Context) {

	// SameSite to prevent warnings in js console
	(*c).SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().AddDate(0, 0, -1),
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
	})
}
