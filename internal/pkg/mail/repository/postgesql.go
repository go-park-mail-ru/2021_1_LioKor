package repository

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/mail"
	"time"
)

type PostgresMailRepository struct {
	DBInstance common.PostgresDataBase
}

func (mr *PostgresMailRepository) GetDialoguesForUser(username string, limit int, find string, folderId int, domain string, since string) ([]mail.Dialogue, error) {
	find = "%" + find + "%"

	query := "SELECT d.id, " +
		"d.other AS email, " +
		"u.avatar_url, m.body, m.received_date, d.unread FROM dialogues d " +
		"JOIN mails m ON d.last_mail_id=m.id " +
		"LEFT JOIN users u ON " +
		"LOWER(SPLIT_PART(d.other, $4, 1))=LOWER(u.username) " +
		"WHERE d.owner=$1 AND " +
		"d.other LIKE $3"

	if since != "" {
		query += " AND d.received_date < $5"
	}

	if find == "" {
		query += " AND d.folder"
		if folderId == 0 {
			query += " IS NULL"
		} else {
			if since != "" {
				query += "=$6"
			} else {
				query += "=$5"
			}
		}
	}

	query += " ORDER BY d.received_date DESC LIMIT $2;"

	var rows pgx.Rows
	var err error
	if find == "" {
		if folderId == 0 {
			if since != "" {
				rows, err = mr.DBInstance.DBConn.Query(
					context.Background(),
					query,
					username,
					limit,
					find,
					domain,
					since,
				)
			} else {
				rows, err = mr.DBInstance.DBConn.Query(
					context.Background(),
					query,
					username,
					limit,
					find,
					domain,
				)
			}
		} else {
			if since != "" {
				rows, err = mr.DBInstance.DBConn.Query(
					context.Background(),
					query,
					username,
					limit,
					find,
					domain,
					since,
					folderId,
				)
			} else {
				rows, err = mr.DBInstance.DBConn.Query(
					context.Background(),
					query,
					username,
					limit,
					find,
					domain,
					folderId,
				)
			}
		}
	} else {
		if since != "" {
			rows, err = mr.DBInstance.DBConn.Query(
				context.Background(),
				query,
				username,
				limit,
				find,
				domain,
				since,
			)
		} else {
			rows, err = mr.DBInstance.DBConn.Query(
				context.Background(),
				query,
				username,
				limit,
				find,
				domain,
			)
		}
	}

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
			&dialogue.AvatarURL,
			&dialogue.Body,
			&dialogue.Received_date,
			&dialogue.Unread,
		)
		if err != nil {
			return nil, err
		}
		dialogues = append(dialogues, dialogue)
	}

	return dialogues, nil
}

func (mr *PostgresMailRepository) DeleteDialogue(owner string, dialogueId int) error {
	_, err := mr.DBInstance.DBConn.Exec(
		context.Background(),
		"DELETE FROM dialogues WHERE id=$1 AND owner=$2;",
		dialogueId,
		owner,
	)
	if err != nil {
		return err
	}
	return nil
}

func (mr *PostgresMailRepository) GetMailsForUser(username string, email string, limit int, last int) ([]mail.DialogueEmail, error) {
	rows, err := mr.DBInstance.DBConn.Query(
		context.Background(),
		"SELECT id, sender, subject, received_date, body, unread, status FROM mails "+
			"WHERE ((sender=$1 AND recipient=$2) OR (sender=$2 AND recipient=$1)) "+
			"AND id > $4 "+
			"ORDER BY id DESC LIMIT $3;",
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
			&mail.Unread,
			&mail.Status,
		)
		if err != nil {
			return nil, err
		}
		mails = append(mails, mail)
	}
	return mails, nil
}

func (mr *PostgresMailRepository) AddMail(mail mail.Mail) (int, error) {
	var mailId int
	err := mr.DBInstance.DBConn.QueryRow(
		context.Background(),
		"INSERT INTO mails(sender, recipient, subject, body) VALUES($1, $2, $3, $4) RETURNING id;",
		mail.Sender,
		mail.Recipient,
		mail.Subject,
		mail.Body,
	).Scan(
		&mailId,
	)
	if err != nil {
		return 0, err
	}
	return mailId, nil
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

func (mr *PostgresMailRepository) ReadDialogue(owner, other string) error {
	_, err := mr.DBInstance.DBConn.Exec(
		context.Background(),
		"UPDATE dialogues SET unread=0 WHERE owner=$1 AND other=$2;",
		owner,
		other,
	)
	if err != nil {
		return err
	}
	return nil
}

func (mr *PostgresMailRepository) CreateFolder(ownerId int, folderName string) (mail.Folder, error) {
	var folder mail.Folder
	err := mr.DBInstance.DBConn.QueryRow(
		context.Background(),
		"INSERT INTO folders(folder_name, owner) "+
			"VALUES($2, $1) RETURNING *;",
		ownerId,
		folderName,
	).Scan(
		&folder.Id,
		&folder.FolderName,
		&folder.Owner,
	)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "folders_owner_fkey" {
				return mail.Folder{}, common.InvalidUserError{"User doesn't exist"}
			}
		}
		return mail.Folder{}, err
	}
	return folder, nil
}

func (mr *PostgresMailRepository) GetFolders(ownerId int) ([]mail.Folder, error) {
	rows, err := mr.DBInstance.DBConn.Query(
		context.Background(),
		"SELECT * FROM folders "+
			"WHERE owner=$1;",
		ownerId,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	folders := make([]mail.Folder, 0, 0)
	for rows.Next() {
		var folder mail.Folder
		err = rows.Scan(
			&folder.Id,
			&folder.FolderName,
			&folder.Owner,
		)
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}
	return folders, nil
}

func (mr *PostgresMailRepository) AddDialogueToFolder(owner string, folderId, dialogueId int) error {

	query := "UPDATE dialogues SET folder="
	if folderId == 0 {
		query += "NULL "
	} else {
		query += "$3 "
	}
	query += "WHERE id=$1 AND owner=$2;"
	var commandTag pgconn.CommandTag
	var err error
	if folderId == 0 {
		commandTag, err = mr.DBInstance.DBConn.Exec(
			context.Background(),
			query,
			dialogueId,
			owner,
		)
	} else {
		commandTag, err = mr.DBInstance.DBConn.Exec(
			context.Background(),
			query,
			dialogueId,
			owner,
			folderId,
		)
	}
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "dialogues_folder_fkey" {
				return mail.InvalidEmailError{"Folder doesn't exist"}
			}
		}
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return mail.InvalidEmailError{"Dialogue doesn't exist"}
	}
	return nil
}

func (mr *PostgresMailRepository) UpdateFolderName(owner, folderId int, folderName string) (mail.Folder, error) {
	var folder mail.Folder
	err := mr.DBInstance.DBConn.QueryRow(
		context.Background(),
		"UPDATE folders SET folder_name=$1 WHERE id=$2 AND owner=$3 RETURNING *;",
		folderName,
		folderId,
		owner,
	).Scan(
		&folder.Id,
		&folder.FolderName,
		&folder.Owner,
	)
	if err != nil {
		return mail.Folder{}, err
	}
	return folder, nil
}

func (mr *PostgresMailRepository) ShiftToMainFolderDialogues(owner string, folderId int) error {
	_, err := mr.DBInstance.DBConn.Exec(
		context.Background(),
		"UPDATE dialogues SET folder=null WHERE folder=$1 AND owner=$2;",
		folderId,
		owner,
	)
	if err != nil {
		return err
	}
	return nil
}

func (mr *PostgresMailRepository) DeleteFolder(owner, folderId int) error {
	_, err := mr.DBInstance.DBConn.Exec(
		context.Background(),
		"DELETE FROM folders WHERE id=$1 AND owner=$2;",
		folderId,
		owner,
	)
	if err != nil {
		return err
	}
	return nil
}

func (mr *PostgresMailRepository) ReadMail(owner, other string) error {
	_, err := mr.DBInstance.DBConn.Exec(
		context.Background(),
		"UPDATE mails SET unread=FALSE WHERE recipient=$1 AND sender=$2;",
		owner,
		other,
	)
	if err != nil {
		return err
	}
	return nil
}

func (mr *PostgresMailRepository) UpdateMailStatus(mailId, status int) error {
	_, err := mr.DBInstance.DBConn.Exec(
		context.Background(),
		"UPDATE mails SET status=$2 WHERE id=$1;",
		mailId,
		status,
	)
	if err != nil {
		return err
	}
	return nil
}
