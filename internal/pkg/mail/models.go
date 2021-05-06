package mail

import (
	"liokor_mail/internal/pkg/common"
	"time"
)

type Mail struct {
	Sender        string    `json:"-"`
	Recipient     string    `json:"recipient"`
	Subject       string    `json:"subject"`
	Body          string    `json:"body"`
	Received_date time.Time `json:"-"`
}

type DialogueEmail struct {
	Id            int       `json:"id"`
	Sender        string    `json:"sender"`
	Subject       string    `json:"title"`
	Received_date time.Time `json:"time"`
	Body          string    `json:"body"`
	Unread        bool      `json:"new"`
	Status        int       `json:"status"`
}

type Dialogue struct {
	Id            int               `json:"id"`
	Email         string            `json:"username"`
	AvatarURL     common.NullString `json:"avatarUrl"`
	Body          string            `json:"body"`
	Received_date time.Time         `json:"time"`
	Unread        int               `json:"new"`
	Owner         string
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
