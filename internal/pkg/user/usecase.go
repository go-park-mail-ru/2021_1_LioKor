package user

import "liokor_mail/internal/pkg/common"

type UseCase interface {
	Login(credentials Credentials) error
	CreateSession(username string) (common.Session, error)
	GetUserByUsername(username string) (User, error)
	GetUserById(id int) (User, error)
	SignUp(newUser UserSignUp) error
	UpdateUser(username string, newData User) (User, error)
	UpdateAvatar(username string, newAvatar string) (User, error)
	ChangePassword(sessionUser User, changePSWD ChangePassword) error
	Logout(sessionToken string) error
}
