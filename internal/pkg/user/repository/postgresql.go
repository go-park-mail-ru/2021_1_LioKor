package repository

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"liokor_mail/internal/pkg/user"
)

type PostgresUserRepository struct {
	DBConn *pgxpool.Conn
	DBpool *pgxpool.Pool
}

func NewPostgresUserRepository(dbConfig string) (user.UserRepository, error){
	dbpool, err := pgxpool.Connect(context.Background(), dbConfig)
	if err != nil {
		return nil, err
	}

	conn, err := dbpool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	return &PostgresUserRepository{conn, dbpool}, nil
}

func (ur *PostgresUserRepository) Close() {
	ur.DBConn.Release()
	ur.DBpool.Close()
}

func (ur *PostgresUserRepository) GetUserDB(username string) (user.UserDB, error) {
	var u user.UserDB
	err := ur.DBConn.QueryRow(
		context.Background(),
		"SELECT * FROM users WHERE LOWER(username)=LOWER($1) LIMIT 1;",
		username,
	).Scan(
		&u.Id,
		&u.Username,
		&u.HashPassword,
		&u.AvatarURL,
		&u.FullName,
		&u.ReserveEmail,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return user.UserDB{}, user.InvalidUserError{"user doesn't exist"}
		} else {
			return user.UserDB{}, err
		}
	}
	return u, nil
}

func (ur *PostgresUserRepository) CreateSession(session user.Session) error {
	u, err := ur.GetUserDB(session.Username)
	if err != nil {
		return err
	}

	_, err = ur.DBConn.Exec(
		context.Background(),
	"INSERT INTO sessions(user_id, token, expiration) VALUES ($1, $2, $3);",
		u.Id,
		session.SessionToken,
		session.Expiration,
		)
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "sessions_user_id_fkey" {
				return user.InvalidUserError{"user doesn't exist"}
			}
		}
		return err
	}
	return nil
}

func (ur *PostgresUserRepository) GetSessionBySessionToken(token string) (user.Session, error) {
	var session user.Session
	err := ur.DBConn.QueryRow(
		context.Background(),
		"SELECT u.username, s.token, s.expiration FROM users AS u " +
			"JOIN sessions AS s ON u.id=s.user_id " +
			"WHERE token=$1 LIMIT 1;",
		token,
	).Scan(
		&session.Username,
		&session.SessionToken,
		&session.Expiration,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return user.Session{}, user.InvalidSessionError{"session doesn't exist"}
		} else {
			return user.Session{}, err
		}
	}

	return session, nil
}

func (ur *PostgresUserRepository) GetUserByUsername(username string) (user.User, error) {
	var u user.User
	err := ur.DBConn.QueryRow(
		context.Background(),
		"SELECT username, password_hash, avatar_url, fullname, reserve_email FROM users WHERE LOWER(username)=LOWER($1) LIMIT 1;",
		username,
	).Scan(
		&u.Username,
		&u.HashPassword,
		&u.AvatarURL,
		&u.FullName,
		&u.ReserveEmail,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return user.User{}, user.InvalidUserError{"user doesn't exist"}
		} else {
			return user.User{}, err
		}
	}
	return u, nil
}

func (ur *PostgresUserRepository) CreateUser(u user.User) error {
	_, err := ur.DBConn.Exec(
		context.Background(),
		"INSERT INTO users(username, password_hash, avatar_url, fullname, reserve_email) VALUES ($1, $2, $3, $4, $5);",
		u.Username,
		u.HashPassword,
		u.AvatarURL,
		u.FullName,
		u.ReserveEmail,
	)

	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "users_username_key" {
				return user.InvalidUserError{"username"}
			}
		}
		return err
	}

	return nil
}

func (ur *PostgresUserRepository) UpdateUser(username string, newData user.User) (user.User, error) {
	err := ur.DBConn.QueryRow(
		context.Background(),
		"UPDATE users SET avatar_url=$1, fullname=$2, reserve_email=$3 " +
			"WHERE LOWER(username)=LOWER($4) " +
			"RETURNING username, password_hash, avatar_url, fullname, reserve_email;",
		newData.AvatarURL,
		newData.FullName,
		newData.ReserveEmail,
		username,
	).Scan(
		&newData.Username,
		&newData.HashPassword,
		&newData.AvatarURL,
		&newData.FullName,
		&newData.ReserveEmail,
	)
	//err now rows
	if err != nil {
		if err == pgx.ErrNoRows {
			return user.User{}, user.InvalidUserError{"user doesn't exist"}
		} else if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "users_username_key" {
				return user.User{}, user.InvalidUserError{"username"}
			}
		}
		return user.User{}, err
	}

	return newData, nil
}

func (ur *PostgresUserRepository) ChangePassword(username string, newPSWD string) error {
	commandTag, err := ur.DBConn.Exec(
		context.Background(),
		"UPDATE users SET password_hash=$1 WHERE LOWER(username)=LOWER($2);",
		newPSWD,
		username,
	)

	if commandTag.RowsAffected() != 1 {
		return user.InvalidUserError{"user doesn't exist"}
	}
	if err != nil {
		return err
	}
	return nil
}

func (ur *PostgresUserRepository) RemoveSession(token string) error {
	commandTag, err := ur.DBConn.Exec(
		context.Background(),
		"DELETE FROM sessions WHERE token=$1;",
		token,
	)

	if commandTag.RowsAffected() != 1 {
		return user.InvalidSessionError{"session doesn't exist"}
	}

	if err != nil {
		return err
	}
	return nil
}

