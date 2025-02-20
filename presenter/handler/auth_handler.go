package handler

import "github.com/labstack/echo/v4"

type AuthHandler interface {
	Save(c echo.Context) error
	Login(c echo.Context) error
	Logout(c echo.Context) error
	ResetPassword(c echo.Context) error
	ForgotPassword(c echo.Context) error
}
