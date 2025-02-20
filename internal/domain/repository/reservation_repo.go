package repository

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/presenter/model"
	"context"
	"time"
)

type ReservationRepository interface {
	Save(ctx context.Context, reservation *entity.ReservationRooms) error
	GetDetailReservation(ctx context.Context, reservationID int) (*model.ReservationDetailService, error)
	GetAllReservationByRoomIDAndDate(ctx context.Context, roomID int64, date string, startTime string, endTime string) ([]entity.ReservationRooms, error)
	GetAllReservation(ctx context.Context, startDate time.Time, endDate time.Time) ([]*model.ReservationGetAllResponse, error)
	GetHistoryReservationByUserID(ctx context.Context, userID int) ([]*model.ReservationHistoryResponse, error)
	GetHistoryReservationByIsAdmin(ctx context.Context) ([]*model.ReservationHistoryResponse, error)
	GetListReservationByRoomID(ctx context.Context, roomID int, startDateTime time.Time, endDateTime time.Time) ([]*model.ReservationListByRoomIdResponse, error)
	GetDashboard(ctx context.Context, startTime time.Time, endTime time.Time) (*model.DashboardResponse, error)
	UpdateStatusReservation(ctx context.Context, reservationID int, status string) error
	DeleteReservation(ctx context.Context, reservationID int) error
}
