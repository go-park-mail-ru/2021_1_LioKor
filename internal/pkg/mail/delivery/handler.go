package delivery

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"liokor_mail/internal/pkg/mail"
	"liokor_mail/internal/pkg/user"
	"net/http"
	"strconv"
	"time"
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
	var sinceTime time.Time
	if since == "" {
		sinceTime = time.Now()
	} else {
		sinceTime, err = time.Parse(time.RFC3339, since)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "неверный формат времени параметра since")
		}
	}
	dialogues, err := h.MailUsecase.GetDialogues(sessionUser.Username, amount, find, folder, sinceTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "не удалось получить диалоги")
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
		return echo.NewHTTPError(http.StatusBadRequest, "неверный формат json")
	}

	dialogue, err := h.MailUsecase.CreateDialogue(sessionUser.Username, dialogueWith.With)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "не получилось создать диалог сорре")
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
		return echo.NewHTTPError(http.StatusBadRequest, "неверный формат json")
	}

	err = h.MailUsecase.DeleteDialogue(sessionUser.Username, deleteDialogue.DialogueId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "не удалось удалить диалог")
	}

	return c.JSON(http.StatusOK, mail.MessageResponse{Message: "Dialogue deleted"})
}

func (h *MailHandler) DeleteMail(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	var idsToDelete struct {
		Ids []int `json:"ids"`
	}
	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&idsToDelete)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "неверный формат json")
	}

	err = h.MailUsecase.DeleteMails(sessionUser.Username, idsToDelete.Ids)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "не получилось удалить сообщение")
	}

	return c.JSON(http.StatusOK, mail.MessageResponse{Message: "Сообщение удаелено"})
}

func (h *MailHandler) GetEmails(c echo.Context) error {
	sUser := c.Get("sessionUser")
	sessionUser, ok := sUser.(user.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	email := c.QueryParam("with")
	if email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "неверный email собеседника")
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
			return echo.NewHTTPError(http.StatusBadRequest, "неверный запрос")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "не удалось получить письма")
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
		return echo.NewHTTPError(http.StatusBadRequest, "неверный формат json")
	}
	newMail.Sender = sessionUser.Username

	email, err := h.MailUsecase.SendEmail(newMail)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "ошибка при отправке")
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
		return echo.NewHTTPError(http.StatusInternalServerError, "не удалось получить папки")
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
		return echo.NewHTTPError(http.StatusBadRequest, "неверный формат json")
	}

	folder, err := h.MailUsecase.CreateFolder(sessionUser.Id, folderName.FolderName)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "не удалось создать папку")
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
		return echo.NewHTTPError(http.StatusBadRequest, "неверный формат json")
	}

	if updateFolder.FolderName != nil {
		folder, err := h.MailUsecase.UpdateFolderName(sessionUser.Id, updateFolder.FolderId, *updateFolder.FolderName)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "не получилось поменять имя папки")
		}

		return c.JSON(http.StatusOK, folder)
	}

	err = h.MailUsecase.UpdateFolderPutDialogue(sessionUser.Username, updateFolder.FolderId, *updateFolder.DialogueId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "не получилось добавить диалог в папку")
	}

	return c.JSON(http.StatusOK, mail.MessageResponse{Message: "Диалог добавлен в папку"})
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
		return echo.NewHTTPError(http.StatusBadRequest, "неверный формат json")
	}

	err = h.MailUsecase.DeleteFolder(sessionUser.Username, sessionUser.Id, deleteFolder.FolderId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "не получилось избавиться от папки")
	}

	return c.JSON(http.StatusOK, mail.MessageResponse{Message: "Folder deleted"})
}
