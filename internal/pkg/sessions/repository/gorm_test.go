package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"liokor_mail/internal/pkg/common"
	"regexp"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	DB *gorm.DB
	mock sqlmock.Sqlmock
	s common.Session

	gsr GormPostgresSessionRepository
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


	s.gsr = GormPostgresSessionRepository{
		common.GormPostgresDataBase{
			s.DB,
		},
	}
	s.s = common.Session{
		UserId: 1,
		SessionToken: "sessionToken",
		Expiration: time.Now().Add(10 * time.Hour),
	}
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGormCreate() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectExec("INSERT INTO").
		WithArgs(
			s.s.UserId,
			s.s.SessionToken,
			s.s.Expiration,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.gsr.Create(s.s)
	require.NoError(s.T(), err)


	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectExec("INSERT INTO").
		WithArgs(
			s.s.UserId,
			s.s.SessionToken,
			s.s.Expiration,
		).
		WillReturnError(&pgconn.PgError{ConstraintName: "sessions_user_id_fkey"})
	s.mock.ExpectRollback()
	err = s.gsr.Create(s.s)
	require.Error(s.T(), err)
}

func (s *Suite) TestGormGet() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "sessions" WHERE token = $1 ORDER BY "sessions"."user_id" LIMIT 1`)).
		WithArgs(s.s.SessionToken).
		WillReturnRows(sqlmock.NewRows([]string{
			"user_id",
			"token",
			"expiration",
		}).
			AddRow(
				s.s.UserId,
				s.s.SessionToken,
				s.s.Expiration,
			))

	_, err := s.gsr.Get(s.s.SessionToken)
	require.NoError(s.T(), err)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "sessions" WHERE token = $1 ORDER BY "sessions"."user_id" LIMIT 1`)).
		WithArgs(s.s.SessionToken).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err = s.gsr.Get(s.s.SessionToken)
	require.Error(s.T(), err)
}

func (s *Suite) TestGormDelete() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "sessions" WHERE token = $1 ORDER BY "sessions"."user_id" LIMIT 1`)).
		WithArgs(s.s.SessionToken).
		WillReturnRows(sqlmock.NewRows([]string{
			"user_id",
			"token",
			"expiration",
		}).
			AddRow(
				s.s.UserId,
				s.s.SessionToken,
				s.s.Expiration,
			))
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectExec("DELETE").
		WithArgs(
			s.s.SessionToken,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.gsr.Delete(s.s.SessionToken)
	require.NoError(s.T(), err)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "sessions" WHERE token = $1 ORDER BY "sessions"."user_id" LIMIT 1`)).
		WithArgs(s.s.SessionToken).
		WillReturnError(gorm.ErrRecordNotFound)
	err = s.gsr.Delete(s.s.SessionToken)
	require.Error(s.T(), err)
}