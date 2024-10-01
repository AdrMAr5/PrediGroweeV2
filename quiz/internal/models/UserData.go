package models

import (
	"encoding/json"
	"io"
)

type UserData struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
}

func (u *UserData) FromJSON(ioReader io.Reader) error {
	return json.NewDecoder(ioReader).Decode(u)
}
