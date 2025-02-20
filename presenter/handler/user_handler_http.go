package handler

import (
	"E-Meeting/internal/usecase"
	"E-Meeting/pkg/helper"
	"E-Meeting/pkg/utils"
	"E-Meeting/presenter/model"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// authHandler handles authentication related HTTP requests.
type userHandler struct {
	useCase usecase.UserUseCase
}

// NewAuthHandler creates a new instance of authHandler.
func NewUserHandler(uc usecase.UserUseCase) UserHandler {
	return &userHandler{useCase: uc}
}

// @Summary Get user by ID
// @Description Get user by ID
// @Tags user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path number true "User ID"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/users/{id} [get]
func (h *userHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, "invalid user ID")
	}

	user, err := h.useCase.GetByID(c.Request().Context(), intID)

	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return utils.JSONResponse(c, http.StatusOK, "user retrieved successfully", user, 0, 0)
}

func (h *userHandler) GetByToken(c echo.Context) error {
	userToken := c.Get("user").(*helper.TokenClaims)
	user, err := h.useCase.GetByID(c.Request().Context(), userToken.UserID)
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return utils.JSONResponse(c, http.StatusOK, "user retrieved successfully", user, 0, 0)
}

// @Summary Update user profile
// @Description Update user profile
// @Tags user
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param user body model.UserUpdateProfileRequest true "User Update Profile Request"
// @Param files formData file false "Profile picture"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/users [put]
func (h *userHandler) Update(c echo.Context) error {
	var user model.UserUpdateProfileRequest
	file, _ := c.FormFile("files")
	if err := c.Bind(&user); err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	id := c.Get("user").(*jwt.MapClaims)
	userID := int((*id)["user_id"].(float64)) // assuming user_id is stored as float64 in the token claims

	err := h.useCase.Update(c.Request().Context(), userID, &user, file)
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return utils.JSONResponse(c, http.StatusOK, "user updated successfully", nil, 0, 0)
}
