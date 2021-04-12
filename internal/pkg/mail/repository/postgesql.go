package repository

import (
	"context"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/mail"
)

type PostgresMailRepository struct {
	DBInstance common.PostgresDataBase
}

func (mr *PostgresMailRepository) GetDialoguesForUser(username string, limit int, offset int) ([]mail.Dialogue, error) {
	rows, err := mr.DBInstance.DBConn.Query(
		context.Background(),
		"SELECT CASE WHEN d.user_1=$1 THEN d.user_2 WHEN d.user_2=$1 THEN d.user_1 END AS email, "+
			"u.avatar_url, m.body, m.received_date FROM dialogues d JOIN mails m ON d.mail_id=m.id "+
			"LEFT  JOIN users u ON "+
			"CASE WHEN d.user_1=$1 THEN SPLIT_PART(d.user_2,'@liokor.ru', 1)=u.username  WHEN d.user_2=$1 THEN SPLIT_PART(d.user_1,'@liokor.ru', 1)=u.username END "+
			"WHERE d.user_1=$1 OR d.user_2=$1 "+
			"ORDER BY d.received_date DESC LIMIT $2 OFFSET $3;",
		username,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	dialogues := make([]mail.Dialogue, 0, 0)
	for rows.Next() {
		dialogue := mail.Dialogue{}
		err = rows.Scan(
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
func (mr *PostgresMailRepository) GetMailsForUser(username string, email string, limit int, offset int) ([]mail.DialogueEmail, error) {
	rows, err := mr.DBInstance.DBConn.Query(
		context.Background(),
		"SELECT sender, subject, received_date, body FROM mails "+
			"WHERE (sender=$1 AND recipient=$2) OR "+
			"(sender=$2 AND recipient=$1) "+
			"ORDER BY received_date DESC LIMIT $3 OFFSET $4;",
		username,
		email,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	mails := make([]mail.DialogueEmail, 0, 0)
	for rows.Next() {
		mail := mail.DialogueEmail{}
		err = rows.Scan(
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
