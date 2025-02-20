package handler

import "github.com/labstack/echo/v4"

type DashboardHandler interface {
	GetDashboard(c echo.Context) error
}
