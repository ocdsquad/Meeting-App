package entity

import (
	"time"

	"github.com/guregu/null"
)

type User struct {
	ID        int         `json:"id" db:"id"`
	Username  string      `json:"username" db:"username"`
	Email     string      `json:"email" db:"email"`
	Password  string      `json:"-" db:"password"`
	IsActive  bool        `json:"is_active" db:"is_active"`
	IsAdmin   bool        `json:"is_admin" db:"is_admin"`
	Language  string      `json:"language" db:"language"`
	AvatarUrl null.String `json:"avatar_url" db:"avatar_url"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}
