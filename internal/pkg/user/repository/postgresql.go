package repository

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/user"
)

type PostgresUserRepository struct {
	DBInstance common.PostgresDataBase
}

func (ur *PostgresUserRepository) GetUserByUsername(username string) (user.User, error) {
	var u user.User
	err := ur.DBInstance.DBConn.QueryRow(
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
			return user.User{}, common.InvalidUserError{"user doesn't exist"}
		} else {
			return user.User{}, err
		}
	}
	return u, nil
}

func (ur *PostgresUserRepository) GetUserById(id int) (user.User, error) {
	var u user.User
	err := ur.DBInstance.DBConn.QueryRow(
		context.Background(),
		"SELECT * FROM users WHERE id=$1 LIMIT 1;",
		id,
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
			return user.User{}, common.InvalidUserError{"user doesn't exist"}
		} else {
			return user.User{}, err
		}
	}
	return u, nil
}

func (ur *PostgresUserRepository) CreateUser(u user.User) error {
	_, err := ur.DBInstance.DBConn.Exec(
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
				return common.InvalidUserError{"username exists"}
			}
		}
		return err
	}

	return nil
}

func (ur *PostgresUserRepository) UpdateUser(newData user.User) (user.User, error) {
	err := ur.DBInstance.DBConn.QueryRow(
		context.Background(),
		"UPDATE users SET fullname=$1, reserve_email=$2 "+
			"WHERE LOWER(username)=LOWER($3) "+
			"RETURNING *;",
		newData.FullName,
		newData.ReserveEmail,
		newData.Username,
	).Scan(
		&newData.Id,
		&newData.Username,
		&newData.HashPassword,
		&newData.AvatarURL,
		&newData.FullName,
		&newData.ReserveEmail,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return user.User{}, common.InvalidUserError{"user doesn't exist"}
		} else if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "users_username_key" {
				return user.User{}, common.InvalidUserError{"username"}
			}
		}
		return user.User{}, err
	}

	return newData, nil
}

func (ur *PostgresUserRepository) ChangePassword(username string, newPSWD string) error {
	commandTag, err := ur.DBInstance.DBConn.Exec(
		context.Background(),
		"UPDATE users SET password_hash=$1 WHERE LOWER(username)=LOWER($2);",
		newPSWD,
		username,
	)

	if commandTag.RowsAffected() != 1 {
		return common.InvalidUserError{"user doesn't exist"}
	}
	if err != nil {
		return err
	}
	return nil
}

func (ur *PostgresUserRepository) UpdateAvatar(username string, newAvatar common.NullString) (user.User, error) {
	newData := user.User{}
	err := ur.DBInstance.DBConn.QueryRow(
		context.Background(),
		"UPDATE users SET avatar_url=$1 "+
			"WHERE LOWER(username)=LOWER($2) "+
			"RETURNING *;",
		newAvatar,
		username,
	).Scan(
		&newData.Id,
		&newData.Username,
		&newData.HashPassword,
		&newData.AvatarURL,
		&newData.FullName,
		&newData.ReserveEmail,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return user.User{}, common.InvalidUserError{"user doesn't exist"}
		}
		return user.User{}, err
	}

	return newData, nil

}

func (ur *PostgresUserRepository) RemoveUser(username string) error {
	//deletes referenced sessions if exists
	commandTag, err := ur.DBInstance.DBConn.Exec(
		context.Background(),
		"DELETE FROM users WHERE username=$1;",
		username,
	)

	if err != nil {
		return err
	}

	if commandTag.RowsAffected() != 1 {
		return common.InvalidUserError{"user doesn't exist"}
	}

	return nil
}
