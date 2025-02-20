package handler

import "github.com/labstack/echo/v4"

type RoomHandler interface {
	FindAllRoom(c echo.Context) error
	Insert(c echo.Context) error
	Update(c echo.Context) error
	DeleteOneByID(c echo.Context) error
	FindOneByID(c echo.Context) error
}
