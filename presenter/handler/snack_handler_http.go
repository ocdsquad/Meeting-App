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

type snackHandler struct {
	useCase usecase.SnackUseCase
}

func NewSnackHandler(uc usecase.SnackUseCase) SnackHandler {
	return &snackHandler{useCase: uc}
}

// @Summary Get all snacks
// @Description Fetch all snacks with pagination and sorting options
// @Tags snack
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param order_by query string false "Order by field"
// @Param sort_by query string false "Sort by order (asc/desc)"
// @Router /api/v1/snacks [get]
func (h *snackHandler) FindAllSnack(c echo.Context) error {
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

	snacks, err := h.useCase.FindAllSnack(ctx, queryPageLimit)
	if err != nil {
		log.Println(fmt.Sprintf("message : error in query count | handler : snack_handler_http | error : %s", err))
		return utils.JSONErrorResponse(c, 404, "data not found")
	}

	return utils.JSONResponse(c, http.StatusOK, "success", snacks.Snacks, snacks.QueryCount.TotalData, snacks.QueryCount.TotalContent)
}
