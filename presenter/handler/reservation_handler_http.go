package handler

import (
	"E-Meeting/internal/usecase"
	"E-Meeting/pkg/utils"
	"E-Meeting/presenter/model"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/guregu/null"
	"github.com/labstack/echo/v4"
)

type reservationHandler struct {
	uc usecase.ReservationUseCase
}

func NewReservationHandler(uc usecase.ReservationUseCase) ReservationHandler {
	return &reservationHandler{
		uc: uc,
	}
}

// Save handles saving reservation code
// @Summary Save a reservation code
// @Description Save the reservation code based on the provided details
// @Tags reservation
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param reservation body model.ReservationCodeRequest true "Reservation Code Request"
// @Router /api/v1/reservations [post]
func (h *reservationHandler) Save(c echo.Context) error {
	userToken := c.Get("user").(interface{})

	claims, ok := userToken.(*jwt.MapClaims)
	if !ok {
		log.Println("failed parse data token")
		return utils.JSONErrorResponse(c, http.StatusUnauthorized, "invalid user token format")

	}
	email, ok := (*claims)["email"].(string)
	if !ok {
		log.Println("failed parse data token")
		return utils.JSONErrorResponse(c, http.StatusUnauthorized, "email not found")
	}

	ctx := context.WithValue(c.Request().Context(), "email", email)

	var request model.ReservationCodeRequest
	if err := c.Bind(&request); err != nil {
		return utils.JSONErrorResponse(c, http.StatusUnprocessableEntity, err.Error())
	}
	log.Printf("Reservation data: %+v\n", request)

	if err := h.uc.Save(ctx, &request); err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "success", nil, 0, 0)
}

// Inquiry handles reservation inquiries
// @Summary Inquiry reservation details
// @Description Get reservation details based on the provided input
// @Tags reservation
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param reservation body model.ReservationCreateRequestSwagger true "Reservation Code Request"
// @Router /api/v1/reservations/inquiry [post]
func (h *reservationHandler) Inquiry(c echo.Context) error {

	userToken := c.Get("user").(interface{})

	claims, ok := userToken.(*jwt.MapClaims)
	if !ok {
		log.Println("failed parse data token")
		return utils.JSONErrorResponse(c, http.StatusUnauthorized, "invalid user token format")

	}
	email, ok := (*claims)["email"].(string)
	if !ok {
		log.Println("failed parse data token")
		return utils.JSONErrorResponse(c, http.StatusUnauthorized, "email not found")
	}

	userID, ok := (*claims)["user_id"].(float64)
	if !ok {
		log.Println("failed parse data token")
		return utils.JSONErrorResponse(c, http.StatusUnauthorized, "unauthorized")
	}

	ctx := context.WithValue(c.Request().Context(), "email", email)

	var request model.ReservationCreateRequest
	if err := c.Bind(&request); err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	result, err := h.uc.Inquiry(ctx, &request, int64(userID))
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	var snack *model.Snack

	if result.Snack != nil {
		snack = &model.Snack{
			ID:       result.Snack.ID,
			Name:     null.StringFrom(result.Snack.Name.String),
			Price:    null.FloatFrom(result.Snack.Price.Float64),
			Currency: null.StringFrom(result.Snack.Currency.String),
			Uom:      null.StringFrom(result.Snack.Uom.String),
		}

	}

	resultResponse := model.ReservationCreateResponse{
		RoomName:                      result.Room.Name.String,
		RoomType:                      result.Room.RoomTypeName.String,
		RoomCapacity:                  result.Room.CapacityID.Int64,
		RoomPrice:                     result.Room.PriceHour.Float64,
		ReservationName:               request.Name,
		ReservationPhone:              request.Phone,
		ReservationOrganization:       request.Organization,
		ReservationDate:               request.Date.String(),
		Duration:                      result.Duration,
		TotalParticipant:              int64(request.TotalParticipant),
		Snack:                         snack,
		GrandTotalPrice:               result.GrandTotalPrice,
		TotalPriceReservationDuration: result.TotalPriceReservationDuration,
		TotalPriceSnack:               result.TotalPriceSnack,
		Note:                          result.RequestInput.Note,
		ReservationCode:               result.ReservationCode,
	}

	return utils.JSONResponse(c, http.StatusOK, "success", resultResponse, 0, 0)
}

// @Summary Get reservation history
// @Description Get reservation history
// @Tags reservation
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/reservations/histories [get]
func (h *reservationHandler) GetHistoryReservation(c echo.Context) error {
	userToken := c.Get("user").(*jwt.MapClaims)
	userID := int((*userToken)["user_id"].(float64))
	isAdmin := bool((*userToken)["is_admin"].(bool))

	log.Printf("User ID: %d\n", userID)

	reservations, err := h.uc.GetHistoryReservation(c.Request().Context(), userID, isAdmin)

	log.Printf("Reservation data: %+v\n", reservations)

	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return utils.JSONResponse(c, http.StatusOK, "reservation retrieved successfully", reservations, 0, 0)
}

// @Summary Update reservation status
// @Description Update reservation status
// @Tags reservation
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path number true "Reservation ID"
// @Param status body model.ReservationUpdateRequest true "Reservation Status Request"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/reservation/{id}/statuses [put]
func (h *reservationHandler) UpdateStatusReservation(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, "invalid reservation ID")
	}

	var status model.ReservationUpdateRequest
	if err := c.Bind(&status); err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.uc.UpdateStatusReservation(c.Request().Context(), intID, status.Status); err != nil {
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "reservation status updated successfully", nil, 0, 0)
}

// @Summary Delete reservation
// @Description Delete reservation
// @Tags reservation
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path number true "Reservation ID"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/reservation/{id} [delete]
func (h *reservationHandler) DeleteReservation(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, "invalid reservation ID")
	}

	if err := h.uc.DeleteReservation(c.Request().Context(), intID); err != nil {
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "reservation deleted successfully", nil, 0, 0)
}

// @Summary Get reservation detail
// @Description Get reservation detail
// @Tags reservation
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path number true "Reservation ID"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/reservations/{id} [get]
func (h *reservationHandler) GetDetailReservation(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, "invalid reservation ID")
	}

	reservation, err := h.uc.GetDetailReservation(c.Request().Context(), intID)
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "reservation retrieved successfully", reservation, 0, 0)
}

// @Summary Get List Reservations by room id
// @Description Get List Reservations by room id
// @Tags room
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path number true "Room ID"
// @Param start_date query string false "Start Date in YYYY-MM-DD format"
// @Param end_date query string false "End Date in YYYY-MM-DD format"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/rooms/{id}/reservations [get]
func (h *reservationHandler) GetListReservationByRoomID(c echo.Context) error {
	id := c.Param("id")
	startTime := c.QueryParam("start_date")
	endTime := c.QueryParam("end_date")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, "invalid room ID")
	}
	var startDateTime, endDateTime time.Time

	if startTime != "" {
		startDateTime, err = time.Parse(time.RFC3339, startTime)
		if err != nil {
			return utils.JSONErrorResponse(c, http.StatusBadRequest, "invalid start_time format")
		}
	}

	if endTime != "" {
		endDateTime, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			return utils.JSONErrorResponse(c, http.StatusBadRequest, "invalid end_time format")
		}
	}

	reservations, err := h.uc.GetListReservationByRoomID(c.Request().Context(), intID, startDateTime, endDateTime)
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "reservation retrieved successfully", reservations, 0, 0)
}

// @Summary Get All Reservations
// @Description Get All Reservations
// @Tags reservation
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param start_date query string false "Start Date in  YYYY-MM-DD format"
// @Param end_date query string false "End Date in  YYYY-MM-DD format"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/reservations [get]
func (h *reservationHandler) GetAllReservation(c echo.Context) error {
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

	reservations, err := h.uc.GetAllReservation(c.Request().Context(), startDateTime, endDateTime)

	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "reservation retrieved successfully", reservations, 0, 0)
}
