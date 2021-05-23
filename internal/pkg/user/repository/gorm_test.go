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
	"liokor_mail/internal/pkg/user"
	"regexp"
	"testing"
)

type Suite struct {
	suite.Suite
	DB *gorm.DB
	mock sqlmock.Sqlmock
	u user.User
	gur GormPostgresUserRepository
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


	s.gur = GormPostgresUserRepository{
		common.GormPostgresDataBase{
			s.DB,
		},
	}
	s.u = user.User{
		Id: 1,
		Username:     "TestUser",
		HashPassword: "hashPassword",
		AvatarURL:    common.NullString{sql.NullString{String: "/media/test", Valid: true}},
		FullName:     "Test User",
		ReserveEmail: "test@test.test",
		RegisterDate: "",
		IsAdmin:      false,
	}
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGormGetUserByUsername() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE LOWER(username) = LOWER($1) ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(s.u.Username).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"username",
			"password_hash",
			"avatar_url",
			"fullname",
			"reserve_email",
	}).
			AddRow(
				s.u.Id,
				s.u.Username,
				s.u.HashPassword,
				s.u.AvatarURL.String,
				s.u.FullName,
				s.u.ReserveEmail,
				))
	_, err := s.gur.GetUserByUsername(s.u.Username)
	require.NoError(s.T(), err)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE LOWER(username) = LOWER($1) ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(s.u.Username).
		WillReturnError(gorm.ErrRecordNotFound)
	_, err = s.gur.GetUserByUsername(s.u.Username)
	require.Error(s.T(), err)
}

func (s *Suite) TestGormGetUserById() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(s.u.Id).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"username",
			"password_hash",
			"avatar_url",
			"fullname",
			"reserve_email",
		}).
			AddRow(
				s.u.Id,
				s.u.Username,
				s.u.HashPassword,
				s.u.AvatarURL.String,
				s.u.FullName,
				s.u.ReserveEmail,
			))
	_, err := s.gur.GetUserById(s.u.Id)
	require.NoError(s.T(), err)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE "users"."id" = $1 ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(s.u.Id).
		WillReturnError(gorm.ErrRecordNotFound)
	_, err = s.gur.GetUserById(s.u.Id)
	require.Error(s.T(), err)
}

func (s *Suite) TestGormCreateUser() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").
		WithArgs(
		s.u.Username,
		s.u.HashPassword,
		s.u.AvatarURL.String,
		s.u.FullName,
		s.u.ReserveEmail,
		s.u.Id,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(s.u.Id))
	s.mock.ExpectCommit()
	err := s.gur.CreateUser(s.u)
	require.NoError(s.T(), err)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").
		WithArgs(
			s.u.Username,
			s.u.HashPassword,
			s.u.AvatarURL.String,
			s.u.FullName,
			s.u.ReserveEmail,
			s.u.Id,
		).
		WillReturnError(&pgconn.PgError{ConstraintName: "users_username_key"})
	s.mock.ExpectRollback()
	err = s.gur.CreateUser(s.u)
	require.Error(s.T(), err)
}

func (s *Suite) TestGormUpdateUser() {
	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WithArgs(
			s.u.Username,
			s.u.HashPassword,
			s.u.AvatarURL.String,
			s.u.FullName,
			s.u.ReserveEmail,
			s.u.Id,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	_, err := s.gur.UpdateUser(s.u)
	require.NoError(s.T(), err)

	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WithArgs(
			s.u.Username,
			s.u.HashPassword,
			s.u.AvatarURL.String,
			s.u.FullName,
			s.u.ReserveEmail,
			s.u.Id,
		).
		WillReturnError(&pgconn.PgError{ConstraintName: "users_username_key"})
	s.mock.ExpectRollback()
	_, err = s.gur.UpdateUser(s.u)
	require.Error(s.T(), err)
}

func (s *Suite) TestGormUpdateAvatar() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE LOWER(username) = LOWER($1) ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(s.u.Username).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"username",
			"password_hash",
			"avatar_url",
			"fullname",
			"reserve_email",
		}).
			AddRow(
				s.u.Id,
				s.u.Username,
				s.u.HashPassword,
				s.u.AvatarURL.String,
				s.u.FullName,
				s.u.ReserveEmail,
			))
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WithArgs(
			s.u.AvatarURL.String,
			s.u.Id,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	_, err := s.gur.UpdateAvatar(s.u.Username, s.u.AvatarURL)
	require.NoError(s.T(), err)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE LOWER(username) = LOWER($1) ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(s.u.Username).
		WillReturnError(gorm.ErrRecordNotFound)
	_, err = s.gur.UpdateAvatar(s.u.Username, s.u.AvatarURL)
	require.Error(s.T(), err)
}

func (s *Suite) TestGormChangePassword() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE LOWER(username) = LOWER($1) ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(s.u.Username).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"username",
			"password_hash",
			"avatar_url",
			"fullname",
			"reserve_email",
		}).
			AddRow(
				s.u.Id,
				s.u.Username,
				s.u.HashPassword,
				s.u.AvatarURL.String,
				s.u.FullName,
				s.u.ReserveEmail,
			))
	s.mock.ExpectBegin()
	s.mock.ExpectExec("UPDATE").
		WithArgs(
			s.u.HashPassword,
			s.u.Id,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.gur.ChangePassword(s.u.Username, s.u.HashPassword)
	require.NoError(s.T(), err)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE LOWER(username) = LOWER($1) ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(s.u.Username).
		WillReturnError(gorm.ErrRecordNotFound)
	err = s.gur.ChangePassword(s.u.Username, s.u.HashPassword)
	require.Error(s.T(), err)
}

func (s *Suite) TestGormRemoveUser() {
	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE LOWER(username) = LOWER($1) ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(s.u.Username).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"username",
			"password_hash",
			"avatar_url",
			"fullname",
			"reserve_email",
		}).
			AddRow(
				s.u.Id,
				s.u.Username,
				s.u.HashPassword,
				s.u.AvatarURL.String,
				s.u.FullName,
				s.u.ReserveEmail,
			))
	s.mock.ExpectBegin()
	s.mock.ExpectExec("DELETE").
		WithArgs(
			s.u.Id,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.gur.RemoveUser(s.u.Username)
	require.NoError(s.T(), err)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "users" WHERE LOWER(username) = LOWER($1) ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(s.u.Username).
		WillReturnError(gorm.ErrRecordNotFound)
	err = s.gur.RemoveUser(s.u.Username)
	require.Error(s.T(), err)
}
