// RoutingRestAPI sets up the routing for the REST API using Echo framework.
// It includes dependency injection, middleware, and route definitions.
//
// @Summary Set up routing for REST API
// @Description This function sets up the routing for the REST API, including dependency injection, middleware, and route definitions.
// @Tags Routing
// @Accept  json
// @Produce  json
// @Param config body configs.Config true "Configuration"
// @Success 200 {object} echo.Echo
// @Failure 500 {object} echo.HTTPError
// @Router /api/v1 [get]
package handler

import (
	"E-Meeting/configs"
	"E-Meeting/internal/domain/repository"
	"E-Meeting/internal/usecase"
	"E-Meeting/pkg/database"
	"E-Meeting/pkg/middleware"
	_ "E-Meeting/presenter/handler/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RoutingRestAPI(e *echo.Echo, config configs.Config) error {

	// Dependency Injection
	dbConn, err := database.NewPostgresConnection(config.DB.Host, config.DB.User, config.DB.Pass, config.DB.Name, config.DB.Port)
	if err != nil {
		e.Logger.Error(err)
		return err
	}
	snackRepo := repository.NewSnackRepository(dbConn)
	snackUsecase := usecase.NewSnackUseCase(snackRepo)
	snackHandlerImpl := NewSnackHandler(snackUsecase)

	roomTypeRepo := repository.NewRoomTypeRepository(dbConn)
	roomTypeUsecase := usecase.NewRoomTypeUseCase(roomTypeRepo)
	roomTypeHandlerImpl := NewRoomTypeHandler(roomTypeUsecase)

	capacityRepo := repository.NewCapacityRepository(dbConn)
	capacityUsecase := usecase.NewCapacityUseCase(capacityRepo)
	capacityHandlerImpl := NewCapacityHandler(capacityUsecase)

	roomRepo := repository.NewRoomRepository(dbConn)
	roomUseCase := usecase.NewRoomUseCase(roomRepo, roomTypeRepo, capacityRepo)
	roomHandlerImpl := NewRoomHandler(roomUseCase)

	attachmentRepo := repository.NewAttachmentRepository(dbConn)
	attachmentUseCase := usecase.NewAttachmentUseCase(attachmentRepo)
	attachmentHandlerImpl := NewAttachmentHandler(attachmentUseCase)

	userRepo := repository.NewUserRepository(dbConn)
	userUsecase := usecase.NewUserUseCase(userRepo)
	userHandler := NewUserHandler(userUsecase)
	authHandler := NewAuthHandler(userUsecase)

	reservationRepo := repository.NewReservationRepository(dbConn)
	reservationUsecase := usecase.NewReservationUseCase(reservationRepo, roomRepo, snackRepo)
	reservationHandler := NewReservationHandler(reservationUsecase)
	dashboardHandler := NewDashboardHandler(reservationUsecase)

	// static route
	e.Static("/uploads", "./uploads")
	// Routing
	apiV1 := e.Group("/api/v1")

	apiV1.POST("/auth/login", authHandler.Login)

	apiV1.POST("/auth/register", authHandler.Save)
	apiV1.POST("/auth/forgot-password", authHandler.ForgotPassword)
	apiV1.POST("/auth/reset-password/:id", authHandler.ResetPassword)

	// Middleware
	// apiV1.Use(middleware.AuthMiddleware)
	// Profile
	apiV1.GET("/users/:id", userHandler.GetByID, middleware.AuthMiddleware)
	apiV1.PUT("/users", userHandler.Update, middleware.AuthMiddleware)
	apiV1.POST("/attachments", attachmentHandlerImpl.Insert)

	// Room
	apiV1.GET("/rooms", roomHandlerImpl.FindAllRoom, middleware.AuthMiddleware)
	apiV1.GET("/rooms/:id/reservations", reservationHandler.GetListReservationByRoomID, middleware.AuthMiddleware)

	// Reservation
	apiV1.POST("/reservations/inquiry", reservationHandler.Inquiry, middleware.AuthMiddleware)
	apiV1.POST("/reservations", reservationHandler.Save, middleware.AuthMiddleware)
	apiV1.GET("/reservations/:id", reservationHandler.GetDetailReservation, middleware.AuthMiddleware)
	apiV1.PUT("/reservation/:id/statuses", reservationHandler.UpdateStatusReservation, middleware.AuthMiddleware)
	apiV1.GET("/reservations/histories", reservationHandler.GetHistoryReservation, middleware.AuthMiddleware)

	// Snack
	apiV1.GET("/snacks", snackHandlerImpl.FindAllSnack)

	// Capacity
	apiV1.GET("/capacities", capacityHandlerImpl.FindAllCapacity)

	// Room type
	apiV1.GET("/room-types", roomTypeHandlerImpl.FindAllRoomType)

	// Admin Routes
	apiV1.GET("/dashboard", dashboardHandler.GetDashboard, middleware.AuthMiddleware, middleware.IsAdminMiddleware)
	apiV1.GET("/reservations", reservationHandler.GetAllReservation, middleware.AuthMiddleware, middleware.IsAdminMiddleware)
	apiV1.POST("/rooms", roomHandlerImpl.Insert, middleware.AuthMiddleware, middleware.IsAdminMiddleware)
	apiV1.PUT("/rooms/:id", roomHandlerImpl.Update, middleware.AuthMiddleware, middleware.IsAdminMiddleware)
	apiV1.DELETE("/rooms/:id", roomHandlerImpl.DeleteOneByID, middleware.AuthMiddleware, middleware.IsAdminMiddleware)
	apiV1.GET("/rooms/:id", roomHandlerImpl.FindOneByID, middleware.AuthMiddleware, middleware.IsAdminMiddleware)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return nil
}
