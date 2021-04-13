package repository

import (
	"context"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/mail"
	"time"
)

type PostgresMailRepository struct {
	DBInstance common.PostgresDataBase
}

func (mr *PostgresMailRepository) GetDialoguesForUser(username string, limit int, last int) ([]mail.Dialogue, error) {
	rows, err := mr.DBInstance.DBConn.Query(
		context.Background(),
		"SELECT d.id, "+
		    "CASE WHEN d.user_1=$1 THEN d.user_2 WHEN d.user_2=$1 THEN d.user_1 END AS email, "+
			"u.avatar_url, m.body, m.received_date FROM dialogues d JOIN mails m ON d.last_mail_id=m.id "+
			"LEFT JOIN users u ON "+
			"CASE WHEN d.user_1=$1 THEN SPLIT_PART(d.user_2,'@liokor.ru', 1)=u.username WHEN d.user_2=$1 THEN SPLIT_PART(d.user_1,'@liokor.ru', 1)=u.username END "+
			"WHERE (d.user_1=$1 OR d.user_2=$1) AND d.id > $3 "+
			"ORDER BY d.received_date DESC LIMIT $2;",
		username,
		limit,
		last,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	dialogues := make([]mail.Dialogue, 0, 0)
	for rows.Next() {
		dialogue := mail.Dialogue{}
		err = rows.Scan(
			&dialogue.Id,
			&dialogue.Email,
			&dialogue.AvatarURLDB,
			&dialogue.Body,
			&dialogue.Received_date,
		)
		if err != nil {
			return nil, err
		}
		//это костыль, который я потом исправлю
		if dialogue.AvatarURLDB.Valid {
			dialogue.AvatarURL = dialogue.AvatarURLDB.String
		} else {
			dialogue.AvatarURL = ""
		}
		dialogues = append(dialogues, dialogue)
	}

	return dialogues, nil
}
func (mr *PostgresMailRepository) GetMailsForUser(username string, email string, limit int, last int) ([]mail.DialogueEmail, error) {
	rows, err := mr.DBInstance.DBConn.Query(
		context.Background(),
		"SELECT id, sender, subject, received_date, body FROM mails "+
			"WHERE "+
			"((sender=$1 AND recipient=$2) OR (sender=$2 AND recipient=$1)) "+
			"AND id > $4 "+
			"ORDER BY id ASC LIMIT $3;",
		username,
		email,
		limit,
	 	last,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	mails := make([]mail.DialogueEmail, 0, 0)
	for rows.Next() {
		mail := mail.DialogueEmail{}
		err = rows.Scan(
			&mail.Id,
			&mail.Sender,
			&mail.Subject,
			&mail.Received_date,
			&mail.Body,
		)
		if err != nil {
			return nil, err
		}
		mails = append(mails, mail)
	}
	return mails, nil
}

func (mr *PostgresMailRepository) AddMail(mail mail.Mail) error {
	_, err := mr.DBInstance.DBConn.Exec(
		context.Background(),
		"INSERT INTO mails(sender, recipient, subject, body) VALUES($1, $2, $3, $4);",
		mail.Sender,
		mail.Recipient,
		mail.Subject,
		mail.Body,
		)
	if err != nil {
		return err
	}
	return nil
}

func (mr *PostgresMailRepository) CountMailsFromUser(username string, interval time.Duration) (int, error) {
	time := time.Now().Add(-interval)
	var num int
	err := mr.DBInstance.DBConn.QueryRow(
		context.Background(),
		"SELECT COUNT(*) FROM mails WHERE sender=$1 AND received_date>$2;",
		username,
		time,
		).Scan(
			&num,
			)

	if err != nil {
		return 0, err
	}
	return num, nil
}
