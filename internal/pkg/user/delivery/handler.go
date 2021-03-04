package delivery

import (
	"encoding/json"
	"lioKor_mail/internal/pkg/user"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type UserHandler struct {
	UserUsecase user.UseCase
}

func (h *UserHandler) Auth(c echo.Context) error{
	creds := user.Credentials{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&creds)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = h.UserUsecase.Login(creds)
	if err != nil {
		switch err.(type) {
		case user.InvalidUserError:
			return c.String(http.StatusUnauthorized, err.Error())
		default:
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	session, err := h.UserUsecase.CreateSession(creds.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name: "session_token",
		Value: session.Value,
		Expires: session.Expiration,
		HttpOnly: true,
	})
	return c.Redirect(http.StatusOK, "/user")
}

func (h *UserHandler) Logout(c echo.Context) error{
	_, err := h.isAuthenticated(c)
	if err != nil {
		return err
	}
	sessionToken, err := c.Cookie("session_token")
	err = h.UserUsecase.Logout(sessionToken.Value)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	c.SetCookie(&http.Cookie{
		Name: "session_token",
		Value: sessionToken.Value,
		Expires: time.Now().AddDate(0, 0, -1),
		HttpOnly: true,
	})
	return c.String(http.StatusOK, "Successfuly logged out")
}

func (h *UserHandler) Profile(c echo.Context) error{
	sessionUser, err := h.isAuthenticated(c)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sessionUser)
}

func (h *UserHandler) ProfileByUsername(c echo.Context) error{
	_, err := h.isAuthenticated(c)
	if err != nil {
		return err
	}

	username := c.Param("username")

	requestedUser, err := h.UserUsecase.GetUserByUsername(username)
	if err != nil {
		switch err.(type) {
		case user.InvalidUserError:
			return c.String(http.StatusNotFound, err.Error())
		default:
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, requestedUser)

}
func (h *UserHandler) SignUp(c echo.Context) error {
	newUser := user.UserSignUp {}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&newUser)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = h.UserUsecase.SignUp(newUser)
	if err != nil {
		switch err.(type) {
		case user.InvalidUserError:
			return c.String(http.StatusConflict, err.Error())
		default:
			return c.String(http.StatusInternalServerError, err.Error())
		}
	}

	return c.String(http.StatusOK,"Signed up successfuly")
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	sessionUser, err := h.isAuthenticated(c)
	if err != nil {
		return err
	}

	username := c.Param("username")
	if username != sessionUser.Username {
		return c.String(http.StatusUnauthorized, "Access denied")
	}

	newData := user.User{}

	defer c.Request().Body.Close()

	err = json.NewDecoder(c.Request().Body).Decode(&newData)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	sessionUser, err = h.UserUsecase.UpdateUser(sessionUser.Username, newData)
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, sessionUser)
}

func (h *UserHandler) ChangePassword(c echo.Context) error {
	sessionUser, err := h.isAuthenticated(c)
	if err != nil {
		return err
	}

	username := c.Param("username")
	if username != sessionUser.Username {
		return c.String(http.StatusUnauthorized, "Access denied")
	}

	changePassword := user.ChangePassword{}

	defer c.Request().Body.Close()

	err = json.NewDecoder(c.Request().Body).Decode(&changePassword)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = h.UserUsecase.ChangePassword(sessionUser, changePassword)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "")
}

func (h *UserHandler) isAuthenticated(c echo.Context) (user.User, error) {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			return user.User{}, c.String(http.StatusUnauthorized, err.Error())
		}
		return user.User{}, c.String(http.StatusBadRequest, err.Error())
	}

	sessionUser, err := h.UserUsecase.GetUserBySessionToken(sessionToken.Value)
	if err != nil {
		switch err.(type) {
		case user.InvalidSessionError, user.InvalidUserError:
			return user.User{}, c.String(http.StatusUnauthorized, err.Error())
		default:
			return user.User{}, c.String(http.StatusInternalServerError, err.Error())
		}
	}
	return sessionUser, nil
}