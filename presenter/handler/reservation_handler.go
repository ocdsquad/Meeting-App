package handler

import "github.com/labstack/echo/v4"

type ReservationHandler interface {
	Save(c echo.Context) error
	GetHistoryReservation(c echo.Context) error
	GetListReservationByRoomID(c echo.Context) error
	GetDetailReservation(c echo.Context) error
	GetAllReservation(c echo.Context) error
	Inquiry(c echo.Context) error
	UpdateStatusReservation(c echo.Context) error
	DeleteReservation(c echo.Context) error
}
