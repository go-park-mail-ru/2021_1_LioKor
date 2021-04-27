package delivery

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"liokor_mail/internal/pkg/user"
	"net/http"
	"time"
)

type UserHandler struct {
	UserUsecase user.UseCase
}

func DeleteSessionCookie(c *echo.Context) {
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

func (h *UserHandler) setSessionCookie(c *echo.Context, username string) error {
	session, err := h.UserUsecase.CreateSession(username)
	if err != nil {
		return err
	}

	(*c).SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    session.Value,
		Path:     "/",
		Expires:  session.Expiration,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
	})

	return nil
}

func (h *UserHandler) Auth(c echo.Context) error {
	creds := user.Credentials{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&creds)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.UserUsecase.Login(creds)
	if err != nil {
		switch err.(type) {
		case user.InvalidUserError:
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	err = h.setSessionCookie(&c, creds.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "ok")
}

func (h *UserHandler) Logout(c echo.Context) error {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	err = h.UserUsecase.Logout(sessionToken.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	DeleteSessionCookie(&c)
	return c.String(http.StatusOK, "Successfuly logged out")
}

func (h *UserHandler) Profile(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, sessionUser)
}

func (h *UserHandler) ProfileByUsername(c echo.Context) error {
	username := c.Param("username")

	requestedUser, err := h.UserUsecase.GetUserByUsername(username)
	if err != nil {
		switch err.(type) {
		case user.InvalidUserError:
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, requestedUser)
}

func (h *UserHandler) SignUp(c echo.Context) error {
	newUser := user.UserSignUp{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.UserUsecase.SignUp(newUser)
	if err != nil {
		switch err.(type) {
		case user.InvalidUserError:
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		default:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	err = h.setSessionCookie(&c, newUser.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Signed up successfuly")
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	username := c.Param("username")
	if username != sessionUser.Username {
		return echo.NewHTTPError(http.StatusUnauthorized, "Access denied")
	}

	newData := user.User{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&newData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sessionUser, err = h.UserUsecase.UpdateUser(sessionUser.Username, newData)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, sessionUser)
}

func (h *UserHandler) ChangePassword(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	username := c.Param("username")
	if username != sessionUser.Username {
		return echo.NewHTTPError(http.StatusUnauthorized, "Access denied")
	}

	changePassword := user.ChangePassword{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&changePassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.UserUsecase.ChangePassword(sessionUser, changePassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "")
}
