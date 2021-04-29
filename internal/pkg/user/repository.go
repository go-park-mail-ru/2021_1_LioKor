package user

import "liokor_mail/internal/pkg/common"

type UserRepository interface {
	GetUserByUsername(username string) (User, error)
	GetUserById(id int) (User, error)
	CreateUser(user User) error
	UpdateUser(username string, newData User) (User, error)
	UpdateAvatar(username string, newAvatar common.NullString) (User, error)
	ChangePassword(username string, newPSWD string) error
	RemoveUser(username string) error
}
