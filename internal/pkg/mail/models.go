package mail

import (
	"database/sql"
	"time"
)

type Mail struct {
	Sender string `json:"-"`
	Recipient string `json:"recipient"`
	Subject string `json:"subject"`
	Body string `json:"body"`
	Received_date  time.Time `json:"-"`
}

type DialogueEmail struct {
	Sender string `json:"username"`
	Subject string `json:"title"`
	Received_date  time.Time `json:"time"`
	Body string `json:"body"`
}

type Dialogue struct {
	Email string `json:"username"`
	AvatarURLDB sql.NullString `json:"-"`
	AvatarURL    string `json:"avatarUrl"`
	Body string `json:"body"`
	Received_date  time.Time `json:"time"`
}

type InvalidEmailError struct {
	Message string
}

func (e InvalidEmailError) Error() string {
	return e.Message
}
