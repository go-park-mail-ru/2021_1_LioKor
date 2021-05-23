package repository

import (
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/mail"
	"time"
)

type GormPostgresMailRepository struct {
	DBInstance common.GormPostgresDataBase
}


func (gmr *GormPostgresMailRepository) AddMail(mail mail.Mail) (int, error) {
	result := gmr.DBInstance.DB.Table("mails").Create(&mail)
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
		"sender=? AND recipient=? AND deleted_sender=FALSE",
			username,
			email,
			).Or(
			"sender=? AND recipient=? AND deleted_recipient=FALSE",
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
		var deleted_by string
		if m.Sender == owner {
			deleted_by = "deleted_by_sender"
		} else if m.Recipient == owner {
			deleted_by = "deleted_by_recipient"
		} else {
			tx.Rollback()
			return mail.InvalidEmailError{
				"Access denied",
			}
		}
		err := tx.Table("mails").
			Where("id=?", id).
			Update(deleted_by, true).Error
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

