package user

import (
	"liokor_mail/internal/pkg/common"
)

type Credentials struct {
	Username string
	Password string
}

type UserSignUp struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	AvatarURL    string `json:"avatarUrl"`
	FullName     string `json:"fullname"`
	ReserveEmail string `json:"reserveEmail"`
}

type User struct {
	Id           int               `json:"-" gorm:"column:id"`
	Username     string            `json:"username" gorm:"column:username"`
	HashPassword string            `json:"-" gorm:"column:password_hash"`
	AvatarURL    common.NullString `json:"avatarUrl" gorm:"column:avatar_url"`
	FullName     string            `json:"fullname" gorm:"column:fullname"`
	ReserveEmail string            `json:"reserveEmail" gorm:"column:reserve_email"`
	RegisterDate string            `json:"-" gorm:"-"`
	IsAdmin      bool              `json:"-" gorm:"-"`
}

type ChangePassword struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type InvalidUsernameError struct {
	Message string
}

func (e InvalidUsernameError) Error() string {
	return e.Message
}

type WeakPasswordError struct {
	Message string
}

func (e WeakPasswordError) Error() string {
	return e.Message
}

