package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/mail"
	"regexp"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	DB *gorm.DB
	mock sqlmock.Sqlmock
	gmr GormPostgresMailRepository

	domain string
	owner string
	other string
	email mail.Mail
	dialogue mail.Dialogue
	folder mail.Folder
	dialogueEmail mail.DialogueEmail
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(s.T(), err)


	s.gmr = GormPostgresMailRepository{
		common.GormPostgresDataBase{
			s.DB,
		},
	}

	s.domain = "liokor.ru"
	s.owner = "liokor"
	s.other = "otherMail@ya.ru"
	s.email = mail.Mail{
		Id: 1,
		Sender: "liokor@liokor.ru",
		Recipient: "otherMail@ya.ru",
		Subject: "Test",
		Body: "Testing test",
	}
	s.folder = mail.Folder{
		Id: 1,
		FolderName: "Cool folder name",
		Owner: 1,
	}
	s.dialogue = mail.Dialogue{
		Id : 1,
		Email: "otherMail@ya.ru",
		Owner: "liokor",
		Body: "Testing test",
		Received_date: s.email.Received_date,
		Unread: 1,
		AvatarURL:common.NullString{ sql.NullString{"/media/test", true}},
	}
	s.dialogueEmail = mail.DialogueEmail{
		Id: 1,
		Sender: "liokor@liokor.ru",
		Subject: "Test",
		Received_date: time.Now(),
		Body: "Testing test",
		Unread: true,
		Status: 1,
	}
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestAddMail() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").
		WithArgs(
			s.email.Sender,
			s.email.Recipient,
			s.email.Subject,
			s.email.Body,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))
	s.mock.ExpectCommit()
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT "id" FROM "dialogues" WHERE owner=$1 AND other=$2 LIMIT 1`)).
		WithArgs(s.owner, s.other).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
		}).
			AddRow(
				1,
			))
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{
		"id",
		"sender",
		"subject",
		"received_date",
		"body",
		"unread",
		"status",
	}).AddRow(
		s.dialogueEmail.Id,
		s.dialogueEmail.Sender,
		s.dialogueEmail.Subject,
		s.dialogueEmail.Received_date,
		s.dialogueEmail.Body,
		s.dialogueEmail.Unread,
		s.dialogueEmail.Status,
		))

	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()

	id, err := s.gmr.AddMail(s.email, s.domain)
	require.NoError(s.T(), err)
	require.Equal(s.T(), s.email.Id, id)


	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").
		WithArgs(
			s.email.Sender,
			s.email.Recipient,
			s.email.Subject,
			s.email.Body,
		).
		WillReturnError(errors.New("Error"))
	s.mock.ExpectRollback()
	_, err = s.gmr.AddMail(s.email, s.domain)
	require.Error(s.T(), err)


	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").
		WithArgs(
			s.email.Recipient,
			s.email.Sender,
			s.email.Subject,
			s.email.Body,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))
	s.mock.ExpectCommit()
	s.mock.ExpectQuery("SELECT").
		WillReturnError(gorm.ErrRecordNotFound)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))
	s.mock.ExpectCommit()
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"sender",
			"subject",
			"received_date",
			"body",
			"unread",
			"status",
		}).AddRow(
			s.dialogueEmail.Id,
			s.dialogueEmail.Sender,
			s.dialogueEmail.Subject,
			s.dialogueEmail.Received_date,
			s.dialogueEmail.Body,
			s.dialogueEmail.Unread,
			s.dialogueEmail.Status,
		))

	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").WillReturnError(gorm.ErrRecordNotFound)
	s.mock.ExpectRollback()

	newEmail := s.email
	newEmail.Sender = s.email.Recipient
	newEmail.Recipient = s.email.Sender
	id, err = s.gmr.AddMail(newEmail, s.domain)
	require.Error(s.T(), err)
}

func (s *Suite) TestGetMailsForUser() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT "id", "sender", "subject", "received_date", "body", "unread", "status" FROM "dialogues" 
			WHERE ((sender=$1 AND recipient=$2 AND deleted_by_sender=FALSE) OR
			(sender=$3 AND recipient=$4 AND deleted_by_recipient=FALSE)) AND
			(id > $5) ORDER BY "id" DESC LIMIT $6`)).
		WithArgs(
			s.email.Sender,
			s.email.Recipient,
			s.email.Recipient,
			s.email.Sender,
			0,
			10,
			).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"sender",
			"subject",
			"received_date",
			"body",
			"unread",
			"status",
		}).AddRow(
			s.dialogueEmail.Id,
			s.dialogueEmail.Sender,
			s.dialogueEmail.Subject,
			s.dialogueEmail.Received_date,
			s.dialogueEmail.Body,
			s.dialogueEmail.Unread,
			s.dialogueEmail.Status,
		))
	_, err := s.gmr.GetMailsForUser(s.email.Sender, s.email.Recipient, 10, 0)
	require.NoError(s.T(), err)
}

func (s *Suite) TestReadMail() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WithArgs(
		false,
			s.email.Sender,
			s.email.Recipient,
			).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.gmr.ReadMail(s.email.Sender, s.email.Recipient)
	require.NoError(s.T(), err)
}

func (s *Suite) TestUpdateMailStatus() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WithArgs(
		2,
			s.email.Id,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.gmr.UpdateMailStatus(s.email.Id, 2)
	require.NoError(s.T(), err)
}

func (s *Suite) TestDeleteMails() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("SELECT").
		WithArgs(s.email.Id).
		WillReturnRows(sqlmock.NewRows([]string{"sender","recipient"}).
			AddRow(s.email.Sender, s.email.Recipient))
	s.mock.ExpectExec("UPDATE").
		WithArgs(
			true,
			s.email.Id,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"sender",
			"subject",
			"received_date",
			"body",
			"unread",
			"status",
		}).AddRow(
			s.dialogueEmail.Id,
			s.dialogueEmail.Sender,
			s.dialogueEmail.Subject,
			s.dialogueEmail.Received_date,
			s.dialogueEmail.Body,
			s.dialogueEmail.Unread,
			s.dialogueEmail.Status,
		))
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectCommit()

	err := s.gmr.DeleteMail(s.owner, []int{1}, s.domain)
	require.NoError(s.T(), err)
}

func (s *Suite) TestCountMailFromUser() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT count(1) FROM "mails" WHERE sender=$1 AND received_date>$2`)).
		WithArgs(s.owner, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	c, err := s.gmr.CountMailsFromUser(s.owner, time.Minute)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, c)
}

func (s *Suite) TestDialogueExists() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT "id" FROM "dialogues" WHERE owner=$1 AND other=$2 LIMIT 1`)).
		WithArgs(s.owner, s.other).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	exists := s.gmr.DialogueExists(s.owner, s.other)
	require.Equal(s.T(), true, exists)


	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT "id" FROM "dialogues" WHERE owner=$1 AND other=$2 LIMIT 1`)).
		WithArgs(s.owner, s.other).
		WillReturnError(gorm.ErrRecordNotFound)
	exists = s.gmr.DialogueExists(s.owner, s.other)
	require.Equal(s.T(), false, exists)
}

func (s *Suite) TestCreateDialogue() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").
		WithArgs(
		s.other,
		s.owner,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))
	s.mock.ExpectCommit()
	d, err := s.gmr.CreateDialogue(s.owner, s.other)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, d.Id)


	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").
		WithArgs(
			s.other,
			s.owner,
		).
		WillReturnError(&pgconn.PgError{ConstraintName: "dialogues_owner_fkey"})
	s.mock.ExpectRollback()
	_, err = s.gmr.CreateDialogue(s.owner, s.other)
	require.Error(s.T(), common.InvalidUserError{"username doesn't exist"}, err)
}

func (s *Suite) TestUpdateDialogueLastMail() {
	s.mock.ExpectQuery("SELECT").
		WithArgs(s.email.Sender, s.email.Recipient, s.email.Recipient, s.email.Sender).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"sender",
			"recipient",
			"received_date",
			"body",
			"unread",
			"status",
		}).AddRow(
			s.email.Id,
			s.email.Sender,
			s.email.Recipient,
			s.dialogueEmail.Received_date,
			s.dialogueEmail.Body,
			s.dialogueEmail.Unread,
			s.dialogueEmail.Status,
		))
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WithArgs(
			s.email.Body,
			s.email.Id,
			s.dialogueEmail.Received_date,
			0,
			s.owner,
			s.other,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.gmr.UpdateDialogueLastMail(s.owner, s.other, "liokor.ru")
	require.NoError(s.T(), err)

	s.mock.ExpectQuery("SELECT").
		WithArgs(s.email.Sender, s.email.Recipient, s.email.Recipient, s.email.Sender).
		WillReturnError(gorm.ErrRecordNotFound)
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err = s.gmr.UpdateDialogueLastMail(s.owner, s.other, "liokor.ru")
	require.NoError(s.T(), err)
}

func (s *Suite) TestGetDialoguesInFolder() {
	since := time.Now()
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{
		"dialogues.id",
		"dialogues.other",
		"users.avatar_url",
		"dialogues.body",
		"dialogues.received_date",
		"dialogues.unread",
		}).AddRow(
			s.dialogue.Id,
			s.dialogue.Email,
			s.dialogue.AvatarURL.String,
			s.dialogue.Body,
			s.dialogue.Received_date,
			s.dialogue.Unread,
	))
	_, err := s.gmr.GetDialoguesInFolder(s.owner, 10, 0, s.domain, since)
	require.NoError(s.T(), err)
}

func (s *Suite) TestFindDialogues() {
	since := time.Now()
	s.mock.ExpectQuery("SELECT").
		WillReturnRows(sqlmock.NewRows([]string{
			"dialogues.id",
			"dialogues.other",
			"users.avatar_url",
			"dialogues.body",
			"dialogues.received_date",
			"dialogues.unread",
		}).AddRow(
			s.dialogue.Id,
			s.dialogue.Email,
			s.dialogue.AvatarURL.String,
			s.dialogue.Body,
			s.dialogue.Received_date,
			s.dialogue.Unread,
		))
	_, err := s.gmr.FindDialogues(s.owner, "o", 10, s.domain, since)
	require.NoError(s.T(), err)

}

func (s *Suite) TestReadDialogue() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WithArgs(
			0,
			s.owner,
			s.other,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.gmr.ReadDialogue(s.owner, s.other)
	require.NoError(s.T(), err)

}

func (s *Suite) TestDeleteDialogue() {
	s.mock.ExpectQuery("SELECT").
		WithArgs(s.dialogue.Id, s.owner).
		WillReturnRows(sqlmock.NewRows([]string{
		"id",
		"owner",
		"other",
		"last_mail_id",
		"received_date",
		"unread",
		"folder",
		"dialogues.body",
	}).AddRow(
		s.dialogue.Id,
		s.dialogue.Owner,
		s.dialogue.Email,
		1,
		s.dialogue.Received_date,
		s.dialogue.Unread,
		1,
		s.dialogue.Body,
	))
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectExec("DELETE").
		WithArgs(
		s.dialogue.Id,
		s.owner,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()


	err := s.gmr.DeleteDialogue(s.owner, s.dialogue.Id, s.domain)
	require.NoError(s.T(), err)

	s.mock.ExpectQuery("SELECT").
		WithArgs(s.dialogue.Id, s.owner).
		WillReturnError(gorm.ErrRecordNotFound)
	err = s.gmr.DeleteDialogue(s.owner, s.dialogue.Id, s.domain)
	require.Error(s.T(), err)
}

func (s *Suite) TestDeleteDialogueMails() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WithArgs(
			true,
			s.email.Sender,
			s.email.Recipient,
			).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WithArgs(
			true,
			s.email.Recipient,
			s.email.Sender,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.gmr.DeleteDialogueMails(s.owner, s.other, s.domain)
	require.Error(s.T(), err)
}

func (s *Suite) TestCreateFolder() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").
		WithArgs(
			s.folder.FolderName,
			s.folder.Owner,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))
	s.mock.ExpectCommit()
	f, err := s.gmr.CreateFolder(s.folder.Owner, s.folder.FolderName)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, f.Id)


	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").
		WithArgs(
			s.folder.FolderName,
			s.folder.Owner,
		).
		WillReturnError(&pgconn.PgError{ConstraintName: "folders_owner_fkey"})
	s.mock.ExpectRollback()
	_, err = s.gmr.CreateFolder(s.folder.Owner, s.folder.FolderName)
	require.Error(s.T(), common.InvalidUserError{"user doesn't exist"}, err)
}

func (s *Suite) TestGetFolders() {
	s.mock.ExpectQuery("SELECT").
		WithArgs(s.folder.Id).
		WillReturnRows(sqlmock.NewRows([]string{
		"id",
		"fodler_name",
		"owner",
		"unread",
	}).AddRow(
		s.folder.Id,
		s.folder.FolderName,
		s.folder.Owner,
		0,
		))
   _, err := s.gmr.GetFolders(s.folder.Owner)
	require.NoError(s.T(), err)
}

func (s *Suite) TestAddDialogueToFolder() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WithArgs(
			s.folder.Id,
			s.dialogue.Id,
			s.owner,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.gmr.AddDialogueToFolder(s.owner, s.folder.Id, s.dialogue.Id)
	require.NoError(s.T(), err)
}

func (s *Suite) TestUpdateFolderName() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WithArgs(
			"New folder name",
			s.folder.Id,
			s.folder.Owner,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	_, err := s.gmr.UpdateFolderName(s.folder.Owner, s.folder.Id, "New folder name")
	require.NoError(s.T(), err)
}

func (s *Suite) TestShiftToMainFolderDialogues() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectCommit()
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WithArgs(
			nil,
			s.owner,
			s.folder.Id,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.gmr.ShiftToMainFolderDialogues(s.owner, s.folder.Id)
	require.NoError(s.T(), err)
}

func (s *Suite) TestDeleteFolder() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectExec("DELETE").
		WithArgs(
			s.folder.Id,
			s.folder.Owner,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.gmr.DeleteFolder(s.folder.Owner, s.folder.Id)
	require.NoError(s.T(), err)
}