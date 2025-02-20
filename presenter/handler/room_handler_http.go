package handler

import (
	"E-Meeting/internal/usecase"
	"E-Meeting/pkg/reason"
	"E-Meeting/pkg/utils"
	"E-Meeting/presenter/model"
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

type roomHandler struct {
	useCase usecase.RoomUseCase
}

func NewRoomHandler(uc usecase.RoomUseCase) RoomHandler {
	return &roomHandler{useCase: uc}
}

// @Summary Get all rooms
// @Description Retrieve a list of rooms with optional pagination and sorting
// @Tags room
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param capacity query int false "filter by capacity of rooms per page (default: 0 for no limit)"
// @Param room_type query int false "filter by room_type of rooms per page (default: 0 for no limit)"
// @Param page query int false "Page number (default: 0)"
// @Param limit query int false "Number of rooms per page (default: 0 for no limit)"
// @Param order_by query string false "Field to order by (default: 'id')"
// @Param sort_by query string false "Sort direction (default: 'desc', can be 'asc')"
// @Router /api/v1/rooms [get]
func (h *roomHandler) FindAllRoom(c echo.Context) error {
	ctx := context.Background()
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	orderBy := c.QueryParam("order_by")
	sortBy := c.QueryParam("sort_by")
	roomType := c.QueryParam("room_type")
	capacity := c.QueryParam("capacity")

	queryPageLimit := utils.QueryPageLimit{
		Page:    0,
		Limit:   0,
		OrderBy: "id",
		SortBy:  "desc",
	}

	filterData := model.FilterDataRoomRequest{
		Capacity: nil,
		RoomType: nil,
	}

	if page != "" {
		pageInt, _ := strconv.Atoi(page)
		queryPageLimit.Page = pageInt
	}

	if limit != "" {
		limitInt, _ := strconv.Atoi(limit)
		queryPageLimit.Limit = limitInt
	}

	if orderBy != "" {
		queryPageLimit.OrderBy = orderBy
	}

	if sortBy != "" {
		queryPageLimit.SortBy = sortBy
	}

	if roomType != "" {
		roomTypeInt, _ := strconv.ParseInt(roomType, 10, 64)
		filterData.RoomType = &roomTypeInt
	}

	if capacity != "" {
		capacityInt, _ := strconv.ParseInt(capacity, 10, 64)
		filterData.Capacity = &capacityInt
	}

	result, err := h.useCase.FindAllRoom(ctx, queryPageLimit, &filterData)
	if err != nil {
		log.Println(fmt.Sprintf("message : data nott found | handler : room_handler_http | error : %s", err))
		return utils.JSONErrorResponse(c, 404, "data not found")
	}

	// Mapping Data
	var rooms []model.Room

	if result.Rooms != nil {
		if len(result.Rooms) > 0 {
			for _, room := range result.Rooms {
				rooms = append(rooms, model.Room{
					ID:            room.ID,
					Name:          room.Name,
					PriceHour:     room.PriceHour,
					IsActive:      room.IsActive,
					Description:   room.Description,
					RoomTypeID:    room.RoomTypeID,
					CapacityID:    room.CapacityID,
					Capacity:      room.Capacity.Int64,
					RoomTypeName:  room.RoomTypeName.String,
					AttachmentURL: room.AttachmentURL,
				})
			}

		}

	}

	return utils.JSONResponse(c, http.StatusOK, "success", rooms, result.QueryCount.TotalData, result.QueryCount.TotalContent)
}

// @Summary Insert a new room
// @Description Upload a new room along with its details
// @Tags room
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Room file upload"
// @Param name formData string true "Room name"
// @Param description formData string false "Room description"
// @Param price formData number true "Room price"
// @Param room_type_id formData number true "Room type ID"
// @Param capacity formData number true "Room capacity"
// @Param user body model.RoomRequestSwagger true "Request Input"
// @Router /api/v1/rooms [post]
func (h *roomHandler) Insert(c echo.Context) error {

	ctx := context.Background()

	userToken := c.Get("user").(interface{})

	claims, ok := userToken.(*jwt.MapClaims)
	if !ok {
		log.Println("failed parse data token")
		return utils.JSONErrorResponse(c, http.StatusUnauthorized, "invalid user token format")

	}

	email, ok := (*claims)["email"].(string)
	if !ok {
		log.Println("failed parse data token")
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, "email not found")
	}

	ctx = context.WithValue(ctx, "email", email)
	var input model.RoomRequest

	file, _ := c.FormFile("files")
	// if err != nil {
	// 	return utils.JSONErrorResponse(c, http.StatusBadRequest, "file not found")
	// }

	if err := c.Bind(&input); err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err := h.useCase.Insert(ctx, input, file)
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "success", nil, 0, 0)
}

// @Summary Update a room
// @Description Update a room along with its details
// @Tags room
// @Security ApiKeyAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Room ID"
// @Param files formData file true "Room file upload"
// @Param name formData string true "Room name"
// @Param description formData string false "Room description"
// @Param price formData number true "Room price"
// @Param room_type_id formData int true "Room type ID"
// @Param capacity_id formData int true "Room capacity"
// @Router /api/v1/rooms/{id} [put]
func (h *roomHandler) Update(c echo.Context) error {

	ctx := context.Background()

	userToken := c.Get("user").(interface{})

	claims, ok := userToken.(*jwt.MapClaims)
	if !ok {
		log.Println("failed parse data token")
		return utils.JSONErrorResponse(c, http.StatusUnauthorized, "invalid user token format")

	}

	email, ok := (*claims)["email"].(string)
	if !ok {
		log.Println("failed parse data token")
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, "email not found")
	}

	ctx = context.WithValue(ctx, "email", email)

	var input model.RoomRequest

	roomID := c.Param("id")

	roomIDInt, _ := strconv.Atoi(roomID)

	file, _ := c.FormFile("files")
	// if err != nil {
	// 	return utils.JSONErrorResponse(c, http.StatusBadRequest, "file not found")
	// }

	if err := c.Bind(&input); err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err := h.useCase.UpdateOneByID(ctx, input, file, int64(roomIDInt))
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "success", nil, 0, 0)
}

// @Summary Delete a room
// @Description Delete a room by its ID
// @Tags room
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Router /api/v1/rooms/{id} [delete]
func (h *roomHandler) DeleteOneByID(c echo.Context) error {

	ctx := context.Background()

	userToken := c.Get("user").(interface{})

	claims, ok := userToken.(*jwt.MapClaims)
	if !ok {
		log.Println("failed parse data token")
		return utils.JSONErrorResponse(c, http.StatusUnauthorized, "invalid user token format")

	}

	email, ok := (*claims)["email"].(string)
	if !ok {
		log.Println("failed parse data token")
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, "email not found")
	}

	ctx = context.WithValue(ctx, "email", email)

	roomID := c.Param("id")

	roomIDInt, _ := strconv.Atoi(roomID)

	err := h.useCase.DeleteOneByID(ctx, int64(roomIDInt))
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, reason.ErrDataNotFound) {
			statusCode = http.StatusNotFound
		}

		return utils.JSONErrorResponse(c, statusCode, err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "success", nil, 0, 0)
}

// @Summary Get a room
// @Description Retrieve details of a room by its ID
// @Tags room
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Room ID"
// @Router /api/v1/rooms/{id} [get]
func (h *roomHandler) FindOneByID(c echo.Context) error {
	ctx := context.Background()

	roomID := c.Param("id")

	roomIDInt, _ := strconv.Atoi(roomID)

	result, err := h.useCase.FindOneByID(ctx, int64(roomIDInt))
	if err != nil {
		statusCode := http.StatusBadRequest
		if errors.Is(err, reason.ErrDataNotFound) {
			statusCode = http.StatusNotFound
		}

		return utils.JSONErrorResponse(c, statusCode, err.Error())
	}

	var room *model.Room
	var reservations []model.Reservation

	if result != nil {

		if len(result.Reservations) > 0 {

			for _, reservation := range result.Reservations {
				reservations = append(reservations, model.Reservation{
					ReservationRoomID: reservation.ReservationRoomID,
					Date:              reservation.Date,
					StartTime:         reservation.StartTime,
					EndTime:           reservation.EndTime,
				})
			}

		}

		room = &model.Room{
			ID:            result.ID,
			Name:          result.Name,
			PriceHour:     result.PriceHour,
			IsActive:      result.IsActive,
			Description:   result.Description,
			RoomTypeID:    result.RoomTypeID,
			CapacityID:    result.CapacityID,
			Capacity:      result.Capacity.Int64,
			AttachmentURL: result.AttachmentURL,
			RoomTypeName:  result.RoomTypeName.String,
			Reservations:  &reservations,
		}
	}

	return utils.JSONResponse(c, http.StatusOK, "success", room, 0, 0)
}
