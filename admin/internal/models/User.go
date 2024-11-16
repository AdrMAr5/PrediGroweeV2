package models

import (
	"encoding/json"
	"io"
)

type User struct {
	ID        int      `json:"id"`
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Email     string   `json:"email"`
	Role      UserRole `json:"role"`
	GoogleID  string   `json:"google_id"`
}

func (u *User) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(u)
}
