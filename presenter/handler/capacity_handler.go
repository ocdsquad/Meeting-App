package handler

import "github.com/labstack/echo/v4"

type CapacityHandler interface {
	FindAllCapacity(c echo.Context) error
}
