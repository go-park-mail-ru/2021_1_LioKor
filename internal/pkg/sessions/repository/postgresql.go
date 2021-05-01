package repository

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"liokor_mail/internal/pkg/common"
)

type PostgresSessionsRepository struct {
	DBInstance common.PostgresDataBase
}

func (sr *PostgresSessionsRepository) Create(session common.Session) error {
	_, err := sr.DBInstance.DBConn.Exec(
		context.Background(),
		"INSERT INTO sessions(user_id, token, expiration) VALUES ($1, $2, $3);",
		session.UserId,
		session.SessionToken,
		session.Expiration,
	)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "sessions_user_id_fkey" {
				return common.InvalidUserError{"user doesn't exist"}
			} else if pgerr.ConstraintName == "sessions_pkey" {
				return common.InvalidSessionError{"sessionToken exists"}
			}
		}
		return err
	}
	return nil
}

func (sr *PostgresSessionsRepository) Get(token string) (common.Session, error) {
	var session common.Session
	err := sr.DBInstance.DBConn.QueryRow(
		context.Background(),
		"SELECT * FROM sessions "+
			"WHERE token=$1 LIMIT 1;",
		token,
	).Scan(
		&session.UserId,
		&session.SessionToken,
		&session.Expiration,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return common.Session{}, common.InvalidSessionError{"session doesn't exist"}
		} else {
			return common.Session{}, err
		}
	}

	return session, nil

}

func (sr *PostgresSessionsRepository) Delete(token string) error {
	commandTag, err := sr.DBInstance.DBConn.Exec(
		context.Background(),
		"DELETE FROM sessions WHERE token=$1;",
		token,
	)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return common.InvalidSessionError{"session doesn't exist"}
	}

	return nil

}