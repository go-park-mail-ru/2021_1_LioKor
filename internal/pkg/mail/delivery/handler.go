package delivery

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"liokor_mail/internal/pkg/mail"
	"liokor_mail/internal/pkg/user"
	"net/http"
	"strconv"
)

type MailHandler struct {
	MailUsecase mail.MailUseCase
}

func (h *MailHandler) GetDialogues(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	last, err := strconv.Atoi(c.QueryParam("last"))
	if err != nil {
		last = 0
	}
	amount, err := strconv.Atoi(c.QueryParam("amount"))
	if err != nil || amount > 50 {
		amount = 50
	}
	find := c.QueryParam("find")

	dialogues, err := h.MailUsecase.GetDialogues(sessionUser.Username, last, amount, find)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, dialogues)
}

func (h *MailHandler) GetEmails(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	email := c.QueryParam("with")
	if email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid email"))
	}

	last, err := strconv.Atoi(c.QueryParam("last"))
	if err != nil {
		last = 0
	}
	amount, err := strconv.Atoi(c.QueryParam("amount"))
	if err != nil || amount > 50 {
		amount = 50
	}
	emails, err := h.MailUsecase.GetEmails(sessionUser.Username, email, last, amount)
	if err != nil {
		switch err.(type) {
		case mail.InvalidEmailError:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, emails)
}

func (h *MailHandler) SendEmail(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	newMail := mail.Mail{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&newMail)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	newMail.Sender = sessionUser.Username

	err = h.MailUsecase.SendEmail(newMail)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "Email sent")
}
