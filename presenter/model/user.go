package model

import "github.com/guregu/null"

type User struct {
	ID        int    `json:"id" db:"id"`
	Username  string `json:"username" validate:"required" db:"username"`
	Email     string `json:"email" validate:"required" db:"email"`
	Password  string `json:"-" validate:"required" db:"password"`
	IsActive  bool   `json:"is_active" db:"is_active"`
	IsAdmin   bool   `json:"is_admin" db:"is_admin"`
	Language  string `json:"language" db:"language"`
	AvatarUrl string `json:"avatar_url" db:"avatar_url"`
}

type UserCreateRequest struct {
	Username        string `json:"username" validate:"required" db:"username"`
	Email           string `json:"email" validate:"required,email" db:"email"`
	Password        string `json:"password" validate:"required" db:"password"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UserDetail struct {
	ID int `json:"id" db:"id"`
}

type UserResetPasswordResponse struct {
	Message string `json:"message"`
}

type UserUpdateProfileRequest struct {
	Username string `json:"username" form:"username" validate:"required" db:"username"`
	Email    string `json:"email" form:"email" validate:"required,email" db:"email"`
	Language string `json:"language" form:"language" validate:"required" db:"language"`
}

type UserGetProfileResponse struct {
	ID        int         `json:"id" db:"id"`
	Username  string      `json:"username" validate:"required" db:"username"`
	Email     string      `json:"email" validate:"required" db:"email"`
	IsActive  bool        `json:"is_active" db:"is_active"`
	IsAdmin   bool        `json:"is_admin" db:"is_admin"`
	Language  string      `json:"language" db:"language"`
	AvatarUrl null.String `json:"avatar_url" db:"avatar_url"`
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required" db:"username"`
	Password string `json:"password" validate:"required" db:"password"`
}

type UserResetPasswordRequest struct {
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UserForgotPasswordRequest struct {
	Email string `json:"email" validate:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}
