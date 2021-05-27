package repository

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
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