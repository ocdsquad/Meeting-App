package handler

import "github.com/labstack/echo/v4"

type SnackHandler interface {
	FindAllSnack(c echo.Context) error
}
