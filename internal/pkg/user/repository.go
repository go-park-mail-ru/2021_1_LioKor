package user

import "liokor_mail/internal/pkg/common"

type UserRepository interface {
	AddUploadedFile(username string, filePath string) error
	GetUserByUsername(username string) (User, error)
	GetUserById(id int) (User, error)
	CreateUser(user User) error
	UpdateUser(newData User) (User, error)
	UpdateAvatar(username string, newAvatar common.NullString) (User, error)
	ChangePassword(username string, newPSWD string) error
	RemoveUser(username string) error
}
