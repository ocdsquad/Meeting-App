package handler

import "github.com/labstack/echo/v4"

type RoomTypeHandler interface {
	FindAllRoomType(c echo.Context) error
}
