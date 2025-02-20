package handler

import "github.com/labstack/echo/v4"

type AttachmentHandler interface {
	Insert(c echo.Context) error
}
