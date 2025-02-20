package utils

import (
	"github.com/labstack/echo/v4"
)

// APIResponse represents the standardized API response structure
// @Description Standard API response format
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Meta    Meta        `json:"meta"`
	Data    interface{} `json:"data"`
}

type Meta struct {
	TotalCount   int `json:"total_count"`
	TotalContent int `json:"total_content`
}

type Data struct {
	Data interface{} `json:"data"`
}

// JSONResponse is a helper to send standardized JSON responses.
func JSONResponse(c echo.Context, status int, message string, data interface{}, totalCount, totalContent int) error {
	response := APIResponse{
		Status:  status,
		Message: message,
		Meta: Meta{
			TotalCount:   totalCount,
			TotalContent: totalContent,
		},
		Data: data,
	}
	return c.JSON(status, response)
}

// JSONErrorResponse is a helper for error responses.
func JSONErrorResponse(c echo.Context, status int, message string) error {
	response := APIResponse{
		Status:  status,
		Message: message,
		Meta: Meta{
			TotalCount:   0,
			TotalContent: 0,
		},
		Data: nil,
	}
	return c.JSON(status, response)
}
