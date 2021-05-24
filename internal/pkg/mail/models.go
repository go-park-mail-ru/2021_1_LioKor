package mail

import (
	"liokor_mail/internal/pkg/common"
	"time"
)

type Mail struct {
	Id            int       `json:"-" gorm:"column:id"`
	Sender        string    `json:"-" gorm:"column:sender"`
	Recipient     string    `json:"recipient" gorm:"column:recipient"`
	Subject       string    `json:"subject" gorm:"column:subject"`
	Body          string    `json:"body" gorm:"column:body"`
	Received_date time.Time `json:"-" gorm:"received_date"`
}

type DialogueEmail struct {
	Id            int       `json:"id" gorm:"column:id"`
	Sender        string    `json:"sender" gorm:"column:sender"`
	Subject       string    `json:"title" gorm:"column:subject"`
	Received_date time.Time `json:"time" gorm:"column:received_date"`
	Body          string    `json:"body" gorm:"column:body"`
	Unread        bool      `json:"new" gorm:"column:unread"`
	Status        int       `json:"status" gorm:"column:status"`
}

type Dialogue struct {
	Id            int               `json:"id" gorm:"column:id"`
	Email         string            `json:"username" gorm:"column:other"`
	AvatarURL     common.NullString `json:"avatarUrl" gorm:"column:avatar_url"`
	Body          string            `json:"body" gorm:"column:body"`
	Received_date time.Time         `json:"time" gorm:"column:received_date"`
	Unread        int               `json:"new" gorm:"column:unread"`
	Owner         string			`gorm:"column:owner"`
}

type Folder struct {
	Id         int    `json:"id"`
	FolderName string `json:"name"`
	Owner      int    `json:"owner"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type InvalidEmailError struct {
	Message string
}

func (e InvalidEmailError) Error() string {
	return e.Message
}
