package handler

import "github.com/labstack/echo/v4"

type UserHandler interface {
	GetByID(c echo.Context) error
	GetByToken(c echo.Context) error
	Update(c echo.Context) error
}
