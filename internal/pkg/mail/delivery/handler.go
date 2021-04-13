package delivery

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/mail"
	"liokor_mail/internal/pkg/user"
	"net/http"
	"strconv"
)

type MailHandler struct {
	MailUsecase mail.MailUseCase
	UserUsecase user.UseCase
}

func (h *MailHandler) GetDialogues(c echo.Context) error {
	user, err := common.IsAuthenticated(&c, h.UserUsecase)
	if err != nil {
		return err
	}

	last, err := strconv.Atoi(c.QueryParam("last"))
	if err != nil {
		last = 0
	}
	amount, err := strconv.Atoi(c.QueryParam("amount"))
	if err != nil || amount > 50 {
		amount = 50
	}

	dialogues, err := h.MailUsecase.GetDialogues(user.Username, last, amount)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, dialogues)
}

func (h *MailHandler) GetEmails(c echo.Context) error {
	user, err := common.IsAuthenticated(&c, h.UserUsecase)
	if err != nil {
		return err
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
	emails, err := h.MailUsecase.GetEmails(user.Username, email, last, amount)
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
	user, err := common.IsAuthenticated(&c, h.UserUsecase)
	if err != nil {
		return err
	}

	newMail := mail.Mail{}

	defer c.Request().Body.Close()

	err = json.NewDecoder(c.Request().Body).Decode(&newMail)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	newMail.Sender = user.Username

	err = h.MailUsecase.SendEmail(newMail)
	if err != nil {
		switch err.(type) {
		case mail.InvalidEmailError:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	return c.String(http.StatusOK, "Email sent")
}
