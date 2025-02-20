package handler

import (
	"E-Meeting/internal/usecase"
	"E-Meeting/pkg/utils"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

type roomTypeHandler struct {
	useCase usecase.RoomTypeUseCase
}

func NewRoomTypeHandler(uc usecase.RoomTypeUseCase) RoomTypeHandler {
	return &roomTypeHandler{useCase: uc}
}

// @Summary Get all room types
// @Description Fetch all room types with pagination and sorting options
// @Tags room_type
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param order_by query string false "Order by field"
// @Param sort_by query string false "Sort by order (asc/desc)"
// @Router /api/v1/room-types [get]
func (h *roomTypeHandler) FindAllRoomType(c echo.Context) error {
	ctx := context.Background()
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	orderBy := c.QueryParam("order_by")
	sortBy := c.QueryParam("sort_by")

	const (
		defaultOrderBy = "id"
		defaultSortBy  = "desc"
	)

	queryPageLimit := utils.QueryPageLimit{
		Page:    0,
		Limit:   0,
		OrderBy: defaultOrderBy,
		SortBy:  defaultSortBy,
	}

	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			return utils.JSONErrorResponse(c, 400, "invalid page parameter")
		}
		queryPageLimit.Page = pageInt
	}

	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			return utils.JSONErrorResponse(c, 400, "invalid limit parameter")
		}
		queryPageLimit.Limit = limitInt
	}

	if orderBy != "" {
		queryPageLimit.OrderBy = orderBy
	}

	if sortBy != "" {
		queryPageLimit.SortBy = sortBy
	}

	results, err := h.useCase.FindAllRoomType(ctx, queryPageLimit)
	if err != nil {
		log.Println(fmt.Sprintf("message : error in query count | handler : room_type_handler_http | error : %s", err))
		return utils.JSONErrorResponse(c, 404, "data not found")
	}

	return utils.JSONResponse(c, http.StatusOK, "success", results.RoomTypes, results.QueryCount.TotalData, results.QueryCount.TotalContent)
}
