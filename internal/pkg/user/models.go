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
	Id           int               `json:"-"`
	Username     string            `json:"username"`
	HashPassword string            `json:"-"`
	AvatarURL    common.NullString `json:"avatarUrl"`
	FullName     string            `json:"fullname"`
	ReserveEmail string            `json:"reserveEmail"`
	RegisterDate string            `json:"-"`
	IsAdmin      bool              `json:"-"`
}

type ChangePassword struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}
