package repository

import (
	"errors"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	"liokor_mail/internal/pkg/common"
	"liokor_mail/internal/pkg/user"
)


type GormPostgresUserRepository struct {
	DBInstance common.GormPostgresDataBase
}

func (ur *GormPostgresUserRepository) GetUserByUsername(username string) (user.User, error) {
	var u user.User
	err := ur.DBInstance.DB.Where("LOWER(username) = LOWER(?)", username).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.User{}, common.InvalidUserError{"user doesn't exist"}
		}
		return user.User{}, err
	}
	return u, nil
}

func (ur *GormPostgresUserRepository) GetUserById(id int) (user.User, error) {
	var u user.User
	err := ur.DBInstance.DB.First(&u, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.User{}, common.InvalidUserError{"user doesn't exist"}
		}
		return user.User{}, err
	}
	return u, nil
}

func (ur *GormPostgresUserRepository) CreateUser(user user.User) error {
	result := ur.DBInstance.DB.Create(&user)
	if err := result.Error; err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "users_username_key" {
				return common.InvalidUserError{"username exists"}
			}
		}
		return err
	}
	return nil
}

func (ur *GormPostgresUserRepository) UpdateUser(newData user.User) (user.User, error) {
	result := ur.DBInstance.DB.Save(&newData)
	if err := result.Error; err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "users_username_key" {
				return user.User{}, common.InvalidUserError{"username"}
			}
		}
		return user.User{}, err
	}
	if result.RowsAffected != 1 {
		return user.User{}, common.InvalidUserError{"Cannot update user"}
	}
	return newData, nil
}

func (ur *GormPostgresUserRepository) UpdateAvatar(username string, newAvatar common.NullString) (user.User, error) {
	u, err := ur.GetUserByUsername(username)
	if err != nil {
		return user.User{}, err
	}
	ur.DBInstance.DB.Model(&u).Update("avatar_url", newAvatar.String)
	if err = ur.DBInstance.DB.Error; err != nil{
		return user.User{}, err
	}
	return u, nil
}

func (ur *GormPostgresUserRepository) ChangePassword(username string, newPSWD string) error {
	u, err := ur.GetUserByUsername(username)
	if err != nil {
		return err
	}
	ur.DBInstance.DB.Model(&u).Update("password_hash", newPSWD)
	if err = ur.DBInstance.DB.Error; err != nil{
		return err
	}
	return nil
}

func (ur *GormPostgresUserRepository) RemoveUser(username string) error {
	u, err := ur.GetUserByUsername(username)
	if err != nil {
		return err
	}
	ur.DBInstance.DB.Delete(&u)
	if err = ur.DBInstance.DB.Error; err != nil || ur.DBInstance.DB.RowsAffected == 0{
		return err
	}
	return nil
}
