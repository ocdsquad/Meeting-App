package handler

import (
	"E-Meeting/internal/usecase"
	"E-Meeting/pkg/utils"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type dashboardHandler struct {
	uc usecase.ReservationUseCase
}

func NewDashboardHandler(uc usecase.ReservationUseCase) DashboardHandler {
	return &dashboardHandler{uc}
}

// @Summary Get dashboard
// @Description Get dashboard
// @Tags dashboard
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param start_date query string false "Start Date in YYYY-MM-DD format"
// @Param end_date query string false "End Date in YYYY-MM-DD format"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/dashboard [get]
func (h *dashboardHandler) GetDashboard(c echo.Context) error {
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	var startDateTime, endDateTime time.Time
	var err error

	if startDate != "" {
		startDateTime, err = time.Parse("2006-01-02", startDate)
		if err != nil {
			return utils.JSONErrorResponse(c, http.StatusBadRequest, "invalid start_date format")
		}
	}
	log.Println("startDateTime", startDateTime)

	if endDate != "" {
		endDateTime, err = time.Parse("2006-01-02", endDate)
		if err != nil {
			return utils.JSONErrorResponse(c, http.StatusBadRequest, "invalid end_date format")
		}
	}
	log.Println("endDateTime", endDateTime)
	dashboard, err := h.uc.GetDashboard(c.Request().Context(), startDateTime, endDateTime)
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "data retrieved", dashboard, 0, 0)

}
