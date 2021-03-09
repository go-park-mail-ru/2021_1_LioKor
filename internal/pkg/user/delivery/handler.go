package delivery

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"liokor_mail/internal/pkg/user"
	"liokor_mail/internal/pkg/user/validators"
	"net/http"
	"time"
)

type UserHandler struct {
	UserUsecase user.UseCase
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

	session, err := h.UserUsecase.CreateSession(creds.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    session.Value,
		Expires:  session.Expiration,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		HttpOnly: true,
	})
	return c.String(http.StatusOK, "ok")
}

func (h *UserHandler) Logout(c echo.Context) error {
	_, httpErr := h.isAuthenticated(c)
	if httpErr != nil {
		return httpErr
	}

	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	err = h.UserUsecase.Logout(sessionToken.Value)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().AddDate(0, 0, -1),
		HttpOnly: true,
	})
	return c.String(http.StatusOK, "Successfuly logged out")
}

func (h *UserHandler) Profile(c echo.Context) error {
	sessionUser, httpErr := h.isAuthenticated(c)
	if httpErr != nil {
		return httpErr
	}

	return c.JSON(http.StatusOK, sessionUser)
}

func (h *UserHandler) ProfileByUsername(c echo.Context) error {
	_, httpErr := h.isAuthenticated(c)
	if httpErr != nil {
		return httpErr
	}

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

	if !validators.ValidateUsername(newUser.Username) {
		return echo.NewHTTPError(http.StatusBadRequest, "username");
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

	return c.String(http.StatusOK, "Signed up successfuly")
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	sessionUser, httpErr := h.isAuthenticated(c)
	if httpErr != nil {
		return httpErr
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
	sessionUser, httpErr := h.isAuthenticated(c)
	if httpErr != nil {
		return httpErr
	}

	username := c.Param("username")
	if username != sessionUser.Username {
		return c.String(http.StatusUnauthorized, "Access denied")
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

func (h *UserHandler) isAuthenticated(c echo.Context) (user.User, error) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return user.User{}, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return user.User{}, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sessionUser, err := h.UserUsecase.GetUserBySessionToken(sessionToken.Value)
	if err != nil {
		switch err.(type) {
		case user.InvalidSessionError, user.InvalidUserError:
			return user.User{}, echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		default:
			return user.User{}, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return sessionUser, nil
}
