package models

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password,omitempty" validate:"required"`
	Role      string `json:"role"`
}

func (u *User) Validate() error {
	return validator.New().Struct(u)
}
func (u *User) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(u)
}
func (u *User) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(u)
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (l *LoginUserPayload) Validate() error {
	return validator.New().Struct(l)
}
func (l *LoginUserPayload) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(l)
}

type UserResponse struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
