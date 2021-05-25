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

	amount, err := strconv.Atoi(c.QueryParam("amount"))
	if err != nil || amount > 50 {
		amount = 50
	}

	find := c.QueryParam("find")

	folder, err := strconv.Atoi(c.QueryParam("folder"))
	if err != nil {
		folder = 0
	}

	since := c.QueryParam("since")

	dialogues, err := h.MailUsecase.GetDialogues(sessionUser.Username, amount, find, folder, since)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, dialogues)
}

func (h *MailHandler) CreateDialogue(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var dialogueWith struct {
		With string `json:"username"`
	}
	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&dialogueWith)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	dialogue, err := h.MailUsecase.CreateDialogue(sessionUser.Username, dialogueWith.With)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, dialogue)
}

func (h *MailHandler) DeleteDialogue(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var deleteDialogue struct {
		DialogueId int `json:"id"`
	}
	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&deleteDialogue)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.MailUsecase.DeleteDialogue(sessionUser.Username, deleteDialogue.DialogueId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, mail.MessageResponse{Message: "Dialogue deleted"})
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

	last, err := strconv.Atoi(c.QueryParam("since"))
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

func (h *MailHandler) DeleteEmails(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	defer c.Request().Body.Close()
	var idsToDelete struct {
		Ids []int `json:"ids"`
	}
	err := json.NewDecoder(c.Request().Body).Decode(&idsToDelete)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.MailUsecase.DeleteEmails(sessionUser.Username, idsToDelete.Ids)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "ok")
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

	email, err := h.MailUsecase.SendEmail(newMail)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, email)
}

func (h *MailHandler) GetFolders(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	folders, err := h.MailUsecase.GetFolders(sessionUser.Id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, folders)
}

func (h *MailHandler) CreateFolder(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var folderName struct {
		FolderName string `json:"name"`
	}
	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&folderName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	folder, err := h.MailUsecase.CreateFolder(sessionUser.Id, folderName.FolderName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, folder)
}

func (h *MailHandler) UpdateFolder(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var updateFolder struct {
		FolderId   int     `json:"folderId"`
		DialogueId *int    `json:"dialogueId"`
		FolderName *string `json:"name"`
	}
	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&updateFolder)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if updateFolder.FolderName != nil {
		folder, err := h.MailUsecase.UpdateFolderName(sessionUser.Id, updateFolder.FolderId, *updateFolder.FolderName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, folder)
	}

	err = h.MailUsecase.UpdateFolderPutDialogue(sessionUser.Username, updateFolder.FolderId, *updateFolder.DialogueId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, mail.MessageResponse{Message: "Dialogue was added to folder"})
}

func (h *MailHandler) DeleteFolder(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var deleteFolder struct {
		FolderId int `json:"id"`
	}
	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&deleteFolder)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = h.MailUsecase.DeleteFolder(sessionUser.Username, sessionUser.Id, deleteFolder.FolderId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, mail.MessageResponse{Message: "Folder deleted"})
}
