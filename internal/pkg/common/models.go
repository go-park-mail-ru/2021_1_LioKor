package common

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullString struct {
	sql.NullString
}

func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

func (ns *NullString) UnmarshalJSON(data []byte) error {
	var b *string
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	if b != nil {
		ns.Valid = true
		ns.String = *b
	} else {
		ns.Valid = false
	}
	return nil
}


type Session struct {
	UserId       int `gorm:"column:user_id"`
	SessionToken string `gorm:"column:token"`
	Expiration   time.Time `gorm:"column:expiration"`
}

type InvalidSessionError struct {
	Message string
}

func (e InvalidSessionError) Error() string {
	return e.Message
}

type InvalidUserError struct {
	Message string
}

func (e InvalidUserError) Error() string {
	return e.Message
}

type InvalidImageError struct {
	Message string
}

func (e InvalidImageError) Error() string {
	return e.Message
}
