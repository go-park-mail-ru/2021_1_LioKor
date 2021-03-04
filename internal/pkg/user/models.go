package user

import "time"

type Credentials struct {
	Username string
	Password string
}


type SessionToken struct {
	Value  string
	Expiration time.Time
}

type UserSignUp struct {
	Username string `json:"username"`
	Password string `json:"password"`
	AvatarURL string `json:"avatar_url"`
	FullName string `json:"fullname"`
	ReserveEmail string `json:"reserve_email"`
}

type User struct {
	Username string `json:"username"`
	HashPassword string `json:"-"`
	AvatarURL string `json:"avatar_url"`
	FullName string `json:"fullname"`
	ReserveEmail string `json:"reserve_email"`
	RegisterDate string `json:"-"`
	IsAdmin bool `json:"-"`
}

type Session struct {
	Username string
	SessionToken string
	Expiration time.Time
}

type ChangePassword struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type InvalidSessionError struct {
	Message string
}

func (e InvalidSessionError) Error() string{
	return e.Message
}

type InvalidUserError struct {
	Message string
}

func (e InvalidUserError) Error() string{
	return e.Message
}

