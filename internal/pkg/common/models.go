package common

import (
	"database/sql"
	"encoding/json"
)

type NullString struct {
	sql.NullString
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
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