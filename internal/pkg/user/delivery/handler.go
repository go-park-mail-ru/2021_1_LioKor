package delivery

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"liokor_mail/internal/pkg/common"
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
	s, err := h.UserUsecase.CreateSession(username)
	if err != nil {
		return err
	}

	(*c).SetCookie(&http.Cookie{
		Name:     "session_token",
		Value:    s.SessionToken,
		Path:     "/",
		Expires:  s.Expiration,
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
		return echo.NewHTTPError(http.StatusBadRequest, "неверный формат json")
	}

	err = h.UserUsecase.Login(creds)
	if err != nil {
		switch err.(type) {
		case common.InvalidUserError:
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	err = h.setSessionCookie(&c, creds.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "не удалось установить печеньки")
	}

	return c.String(http.StatusOK, "ok")
}

func (h *UserHandler) Logout(c echo.Context) error {
	sessionToken, err := c.Cookie("session_token")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	err = h.UserUsecase.Logout(sessionToken.Value)
	if err != nil {
		switch err.(type) {
		case common.InvalidSessionError:
			return echo.NewHTTPError(http.StatusUnauthorized)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "что-то пошло не так")
		}
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
		case common.InvalidUserError:
			return echo.NewHTTPError(http.StatusNotFound, "нет такого пользователя")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "что-то пошло не так")
		}
	}

	return c.JSON(http.StatusOK, requestedUser)
}

func (h *UserHandler) SignUp(c echo.Context) error {
	newUser := user.UserSignUp{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "неверный формат json")
	}

	err = h.UserUsecase.SignUp(newUser)
	if err != nil {
		switch err.(type) {
		case common.InvalidUserError:
			return echo.NewHTTPError(http.StatusConflict, "такое имя уже занято")
		case user.InvalidUsernameError:
			return echo.NewHTTPError(http.StatusBadRequest, "попробуйте другое имя")
		case user.WeakPasswordError:
			return echo.NewHTTPError(http.StatusBadRequest, "слишком слабый пароль")
		default:
			return echo.NewHTTPError(http.StatusBadRequest, "что-то пошло не так")
		}
	}

	err = h.setSessionCookie(&c, newUser.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "не удалось установить печеньки")
	}

	return c.String(http.StatusOK, "вы успешно зарегестрированы")
}

func (h *UserHandler) UpdateProfile(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	username := c.Param("username")
	if username != sessionUser.Username {
		return echo.NewHTTPError(http.StatusUnauthorized, "не лезь")
	}

	newData := user.User{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&newData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "неверный формат json")
	}

	sessionUser, err = h.UserUsecase.UpdateUser(sessionUser.Username, newData)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	return c.JSON(http.StatusOK, sessionUser)
}

func (h *UserHandler) UploadImage(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var uploadImage struct {
		DataUrl string `json:"dataUrl"`
	}

	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&uploadImage)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "неверный формат json")
	}

	imagePath, err := h.UserUsecase.UploadImage(sessionUser.Username, uploadImage.DataUrl)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var uploadedImage struct {
		Url string `json:"url"`
	}
	uploadedImage.Url = imagePath
	return c.JSON(http.StatusOK, uploadedImage)
}

func (h *UserHandler) UpdateAvatar(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	username := c.Param("username")
	if username != sessionUser.Username {
		return echo.NewHTTPError(http.StatusUnauthorized, "не твое")
	}

	var newAvatar struct {
		AvatarUrl string `json:"avatarUrl"`
	}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&newAvatar)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "неверный формат json")
	}

	sessionUser, err = h.UserUsecase.UpdateAvatar(sessionUser.Username, newAvatar.AvatarUrl)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
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
		return echo.NewHTTPError(http.StatusUnauthorized, "ну нельзя")
	}

	changePassword := user.ChangePassword{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&changePassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "неверный формат json")
	}

	err = h.UserUsecase.ChangePassword(sessionUser, changePassword)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "не получилось поменять пароль")
	}

	return c.String(http.StatusOK, "")
}
