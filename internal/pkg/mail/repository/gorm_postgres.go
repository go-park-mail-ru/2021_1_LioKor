package repository

import (
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/mail"
	"time"
)

type GormPostgresMailRepository struct {
	DBInstance common.GormPostgresDataBase
}


func (gmr *GormPostgresMailRepository) AddMail(mail mail.Mail) (int, error) {
	result := gmr.DBInstance.DB.
		Table("mails").
		Select("sender", "recipient", "subject", "body").
		Create(&mail)
	if err := result.Error; err != nil {
		return 0, err
	}
	return mail.Id, nil
}

func (gmr *GormPostgresMailRepository) GetMailsForUser(username string, email string, limit int, last int) ([]mail.DialogueEmail, error) {
	mails := make([]mail.DialogueEmail, 0, 0)
	gmr.DBInstance.DB.
		Table("mails").
		Select("id, sender, subject, received_date, body, unread, status").
		Limit(limit).
		Order("id desc").
		Where(
		gmr.DBInstance.DB.Where(
		"sender=? AND recipient=? AND deleted_by_sender=FALSE",
			username,
			email,
			).Or(
			"sender=? AND recipient=? AND deleted_by_recipient=FALSE",
			email,
			username,
		)).
		Where(
			"id > ?",
			last,
			).
		Scan(&mails)
	if err := gmr.DBInstance.DB.Error; err != nil {
		return nil, err
	}
	return mails, nil
}

func (gmr *GormPostgresMailRepository) ReadMail(owner, other string) error {
	result := gmr.DBInstance.DB.
		Table("mails").
		Where(
			"recipient=? AND sender=?",
			owner,
			other,
			).
		Update("unread", false)
	if err:= result.Error; err != nil {
		return err
	}
	return nil
}

func (gmr *GormPostgresMailRepository) UpdateMailStatus(mailId, status int) error {
	result := gmr.DBInstance.DB.
		Table("mails").
		Where(
			"id=?",
			mailId,
		).
		Update("status", status)
	if err:= result.Error; err != nil {
		return err
	}
	return nil
}

func (gmr *GormPostgresMailRepository) DeleteMail(owner string, mailIds []int) error {
	tx := gmr.DBInstance.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	for _, id := range mailIds {
		var m struct {
			Sender    string `gorm:"sender"`
			Recipient string `gorm:"recipient"`
		}
		tx.Table("mails").
			Select("sender, recipient").
			Where("id=?", id).
			Take(&m)
		var deletedBy string
		if m.Sender == owner {
			deletedBy = "deleted_by_sender"
		} else if m.Recipient == owner {
			deletedBy = "deleted_by_recipient"
		} else {
			tx.Rollback()
			return mail.InvalidEmailError{
				"Access denied",
			}
		}
		err := tx.Table("mails").
			Where("id=?", id).
			Update(deletedBy, true).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func (gmr *GormPostgresMailRepository) CountMailsFromUser(username string, interval time.Duration) (int, error) {
	timeLimit := time.Now().Add(-interval)
	var count int64
	err := gmr.DBInstance.DB.
		Table("mails").
		Where(
			"sender=? AND received_date>?",
			username,
			timeLimit,
			).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (gmr *GormPostgresMailRepository) CreateDialogue(owner string, other string) (mail.Dialogue, error) {
	dialogue := mail.Dialogue {
		Owner: owner,
		Email: other,
	}
	result := gmr.DBInstance.DB.
		Table("dialogues").
		Select("owner","other").
		Create(&dialogue)
	if err := result.Error; err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "dialogues_owner_fkey" {
				return mail.Dialogue{}, common.InvalidUserError{"username doesn't exist"}
			} else if pgerr.ConstraintName == "dialogues_owner_other_key" {
				return mail.Dialogue{}, mail.InvalidEmailError{"dialogue already exists"}
			}
		}
		return mail.Dialogue{}, err
	}
	return dialogue, nil
}


func (gmr *GormPostgresMailRepository) UpdateDialogueWithMailId(owner string, lastMailId int) error {
	lastMail := mail.Mail {
		Id : lastMailId,
	}
	err := gmr.DBInstance.DB.Table("mails").Take(&lastMail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mail.InvalidEmailError{"Mail doesn't exist"}
		}
		return err
	}
	var other string
	switch owner {
	case lastMail.Sender:
		other = lastMail.Recipient
	case lastMail.Recipient:
		other = lastMail.Sender
	default:
		return common.InvalidUserError{"Access denied"}
	}
	result := gmr.DBInstance.DB.
		Table("dialogues").
		Where(
			"owner=? AND other=?",
			owner,
			other,
			).
		Updates(map[string]interface{}{
		"last_mail_id" : lastMailId,
		"received_date" : lastMail.Received_date,
		"unread": gorm.Expr("unread + 1"),
		"body" : lastMail.Body,
	})
	if err = result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mail.InvalidEmailError{"Dialogue doesn't exist"}
		}
		return err
	}
	return nil
}

func (gmr *GormPostgresMailRepository) UpdateDialogueLastMail(owner string, other string) error {
	var lastMail mail.DialogueEmail
	err := gmr.DBInstance.DB.
		Table("mails").
		Where(
			gmr.DBInstance.DB.Where(
				"sender=? AND recipient=? AND deleted_by_sender=FALSE",
				owner,
				other,
			).Or(
				"sender=? AND recipient=? AND deleted_by_recipient=FALSE",
				other,
				owner,
			)).
		Last(&lastMail).Error
	if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return mail.InvalidEmailError{"Mail doesn't exist"}
			}
			return err
	}
	updates := map[string]interface{}{
		"last_mail_id" : lastMail.Id,
		"received_date" : lastMail.Received_date,
		"body" : lastMail.Body,
	}
	if lastMail.Unread {
		updates["unread"] = gorm.Expr("unread + 1")
	} else {
		updates["unread"] = 0
	}
	result := gmr.DBInstance.DB.
		Table("dialogues").
		Where(
			"owner=? AND other=?",
			owner,
			other,
		).
		Updates(updates)
	if err = result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mail.InvalidEmailError{"Dialogue doesn't exist"}
		}
		return err
	}
	return nil
}

func (gmr *GormPostgresMailRepository) GetDialoguesInFolder(username string, limit int, folderId int, domain string, since string) ([]mail.Dialogue, error) {
	var sinceTime time.Time
	var err error
	if since == "" {
		sinceTime = time.Now().Add(time.Second)
	} else {
		sinceTime, err = time.Parse(time.RFC3339, since)
		if err != nil {
			return nil, mail.InvalidEmailError{"Invalid parameters"}
		}
	}
	dialogues := make([]mail.Dialogue, 0, 0)
	var folderCond string
	if folderId == 0 {
		folderCond = "dialogues.folder IS NULL"
	} else {
		folderCond = fmt.Sprintf("dialogues.folder=%d", folderId)
	}
	err = gmr.DBInstance.DB.
		Table("dialogues").
		Limit(limit).
		Order("dialogues.received_date desc").
		Where("dialogues.owner=?", username).
		Where(folderCond).
		Where("dialogues.received_date<?", sinceTime).
		Select(
			"dialogues.id",
			"dialogues.other",
				"users.avatar_url",
				"dialogues.body",
				"dialogues.received_date",
				"dialogues.unread",
			).
		Joins("LEFT JOIN users ON LOWER(SPLIT_PART(dialogues.other, ?, 1))=LOWER(users.username)", domain).
		Scan(&dialogues).Error

	if err != nil {
		return nil, err
	}
	return dialogues, nil
}

func (gmr *GormPostgresMailRepository) FindDialogues(username string, find string, limit int, domain string, since string) ([]mail.Dialogue, error) {
	var sinceTime time.Time
	var err error
	if since == "" {
		sinceTime = time.Now().Add(time.Second)
	} else {
		sinceTime, err = time.Parse(time.RFC3339, since)
		if err != nil {
			return nil, mail.InvalidEmailError{"Invalid parameters"}
		}
	}
	dialogues := make([]mail.Dialogue, 0, 0)
	err = gmr.DBInstance.DB.
		Table("dialogues").
		Limit(limit).
		Order("dialogues.received_date desc").
		Where("dialogues.owner=?", username).
		Where("dialogues.received_date<?", sinceTime).
		Where("dialogues.other LIKE ?", "%" + find + "%").
		Select(
			"dialogues.id",
			"dialogues.other",
			"users.avatar_url",
			"dialogues.body",
			"dialogues.received_date",
			"dialogues.unread",
		).
		Joins("LEFT JOIN users ON LOWER(SPLIT_PART(dialogues.other, ?, 1))=LOWER(users.username)", domain).
		Scan(&dialogues).Error

	if err != nil {
		return nil, err
	}
	return dialogues, nil
}

func (gmr *GormPostgresMailRepository) ReadDialogue(owner, other string) error {
		result := gmr.DBInstance.DB.
		Table("dialogues").
		Where(
		"owner=? AND other=?",
		owner,
		other,
	).
		Update("unread", 0)
		if err:= result.Error; err != nil {
		return err
	}
		return nil
}
func (gmr *GormPostgresMailRepository) DeleteDialogue(owner string, dialogueId int) error {
	gmr.DBInstance.DB.
		Table("dialogues").
		Where("id=? AND owner=?", dialogueId, owner).
		Delete(&mail.Dialogue{})
	if err := gmr.DBInstance.DB.Error; err != nil || gmr.DBInstance.DB.RowsAffected == 0{
		return err
	}
	return nil
}

func (gmr *GormPostgresMailRepository) CreateFolder(ownerId int, folderName string) (mail.Folder, error) {
	folder := mail.Folder {
		FolderName: folderName,
		Owner: ownerId,
	}
	result := gmr.DBInstance.DB.Table("folders").Create(&folder)
	if err := result.Error; err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "folders_owner_fkey" {
				return mail.Folder{}, common.InvalidUserError{"user doesn't exist"}
			} else if pgerr.ConstraintName == "folders_folder_name_owner_key" {
				return mail.Folder{}, mail.InvalidEmailError{"folder already exists"}
			}
		}
		return mail.Folder{}, err
	}
	return folder, nil
}

func (gmr *GormPostgresMailRepository) GetFolders(ownerId int) ([]mail.Folder, error) {
	folders := make([]mail.Folder, 0, 0)
	err := gmr.DBInstance.DB.Table("folders").
		Where("owner=?", ownerId).
		Scan(&folders).Error
	if err != nil {
		return nil, err
	}
	return folders, nil
}

func (gmr *GormPostgresMailRepository) AddDialogueToFolder(owner string, folderId, dialogueId int) error {
	updates := map[string]interface{}{"folder": folderId}
	if folderId == 0 {
		updates["folder"] = nil
	}
	result := gmr.DBInstance.DB.
		Table("dialogues").
		Where("id=? AND owner=?", dialogueId, owner).
		Select("folder").
		Updates(updates)
	if err := result.Error; err != nil {	if pgerr, ok := err.(*pgconn.PgError); ok {
		if pgerr.ConstraintName == "dialogues_folder_fkey" {
			return mail.InvalidEmailError{"Folder doesn't exists"}
		}
	}
		return err
	}
	return nil
}

func (gmr *GormPostgresMailRepository) UpdateFolderName(owner, folderId int, folderName string) (mail.Folder, error) {
	folder := mail.Folder {
		Id : folderId,
		FolderName: folderName,
		Owner: owner,
	}
	result := gmr.DBInstance.DB.Table("folders").
		Where("id=? AND owner=?", folderId, owner).
		Update("folder_name", folderName)
	if err := result.Error; err != nil {
		return mail.Folder{}, err
	}
	return folder, nil
}

func (gmr *GormPostgresMailRepository) ShiftToMainFolderDialogues(owner string, folderId int) error {
	result := gmr.DBInstance.DB.Table("dialogues").
		Where("owner=? AND folder=?", owner, folderId).
		Update("folder", nil)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (gmr *GormPostgresMailRepository) DeleteFolder(owner, folderId int) error {
	gmr.DBInstance.DB.
		Table("folders").
		Where("id=? AND owner=?", folderId, owner).
		Delete(&mail.Folder{})
	if err := gmr.DBInstance.DB.Error; err != nil || gmr.DBInstance.DB.RowsAffected == 0{
		return err
	}
	return nil
}