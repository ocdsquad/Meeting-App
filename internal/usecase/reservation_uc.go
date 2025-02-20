package usecase

import (
	"E-Meeting/presenter/model"
	"context"
	"time"
)

type ReservationUseCase interface {
	Save(ctx context.Context, code *model.ReservationCodeRequest) error
	Inquiry(ctx context.Context, request *model.ReservationCreateRequest, userID int64) (*model.ReservationCreateServiceResponse, error)
	GetDetailReservation(ctx context.Context, reservationID int) (*model.ReservationDetailResponse, error)
	GetHistoryReservation(ctx context.Context, userID int, isAdmin bool) ([]model.ReservationHistoryResponse, error)
	GetDashboard(ctx context.Context, startDateTime time.Time, endDateTime time.Time) (*model.DashboardResponse, error)
	GetAllReservation(ctx context.Context, startDate time.Time, endDate time.Time) ([]model.ReservationGetAllResponse, error)
	GetListReservationByRoomID(ctx context.Context, roomID int, startDateTime time.Time, endDateTime time.Time) ([]model.ReservationListByRoomIdResponse, error)
	UpdateStatusReservation(ctx context.Context, reservationID int, status string) error
	DeleteReservation(ctx context.Context, reservationID int) error
}
