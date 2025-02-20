package handler

import (
	"E-Meeting/internal/usecase"
	"E-Meeting/pkg/mailer"
	"E-Meeting/pkg/utils"
	"E-Meeting/presenter/model"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// authHandler handles authentication related HTTP requests.
type authHandler struct {
	useCase usecase.UserUseCase
}

// NewAuthHandler creates a new instance of authHandler.
func NewAuthHandler(uc usecase.UserUseCase) AuthHandler {
	return &authHandler{useCase: uc}
}

// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.UserCreateRequest true "User Create Request"
// @Success 201 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/auth/register [post]
func (h *authHandler) Save(c echo.Context) error {
	user := new(model.UserCreateRequest)

	if err := c.Bind(user); err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	log.Printf("User data: %+v\n", user)
	if err := h.useCase.Save(c.Request().Context(), user); err != nil {
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, err.Error())

	}
	return utils.JSONResponse(c, http.StatusCreated, "user created successfully", nil, 0, 0)
}

// @Summary Login
// @Description Login with the provided username and password, user has role admin or user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.UserLoginRequest true "User Login Request"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/auth/login [post]
func (h *authHandler) Login(c echo.Context) error {
	var user model.UserLoginRequest
	if err := c.Bind(&user); err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	token, err := h.useCase.Login(c.Request().Context(), &user)
	if err != nil {
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "login successful", token, 0, 0)
}

func (h *authHandler) Logout(c echo.Context) error {
	return utils.JSONResponse(c, http.StatusOK, "logout successful", nil, 0, 0)
}

// @Summary Reset Password
// @Description Reset password with the provided new password
// @Tags auth
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Param user body model.UserResetPasswordRequest true "User Reset Password Request"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/auth/reset-password/{id} [post]
func (h *authHandler) ResetPassword(c echo.Context) error {
	otp := c.Param("id")
	log.Printf("OTP: %s\n", otp)
	var user model.UserResetPasswordRequest
	if err := c.Bind(&user); err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := h.useCase.ResetPassword(c.Request().Context(), &user, otp); err != nil {
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "reset password successful", nil, 0, 0)
}

// @Summary Forgot Password
// @Description Send OTP to email for password reset
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.UserForgotPasswordRequest true "User Forgot Password Request"
// @Success 200 {object} utils.APIResponse
// @Failure 400 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /api/v1/auth/forgot-password [post]
func (h *authHandler) ForgotPassword(c echo.Context) error {
	var user model.UserForgotPasswordRequest
	if err := c.Bind(&user); err != nil {
		return utils.JSONErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	log.Printf("User data: %+v\n", user)
	mailer := mailer.NewMailer(c.Logger())

	if err := h.useCase.ForgotPassword(c.Request().Context(), &user, mailer); err != nil {
		return utils.JSONErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return utils.JSONResponse(c, http.StatusOK, "forgot password request successful, OTP sent to email", nil, 0, 0)
}
