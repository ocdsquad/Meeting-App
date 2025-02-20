package handler

import (
	"E-Meeting/internal/usecase"
	"E-Meeting/pkg/reason"
	"E-Meeting/pkg/utils"
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

type capacityHandler struct {
	useCase usecase.CapacityUseCase
}

func NewCapacityHandler(uc usecase.CapacityUseCase) CapacityHandler {
	return &capacityHandler{useCase: uc}
}

// @Summary Get all capacities
// @Description Fetch all capacity with pagination and sorting options
// @Tags capacity
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param order_by query string false "Order by field"
// @Param sort_by query string false "Sort by order (asc/desc)"
// @Router /api/v1/capacities [get]
func (h *capacityHandler) FindAllCapacity(c echo.Context) error {
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

	results, err := h.useCase.FindAllCapacity(ctx, queryPageLimit)
	if err != nil {

		statusCode := http.StatusNotFound
		message := "data not found"
		if !errors.Is(err, reason.ErrDataNotFound) {
			log.Println(fmt.Sprintf("message : error get data | handler : capacity_handler_http | error : %s", err))
			statusCode = http.StatusInternalServerError
			message = "internal server error"
		}

		return utils.JSONErrorResponse(c, statusCode, message)
	}

	return utils.JSONResponse(c, http.StatusOK, "success", results.Capacity, results.QueryCount.TotalData, results.QueryCount.TotalContent)
}
