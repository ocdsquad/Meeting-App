package repository

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/pkg/utils"
	"E-Meeting/presenter/model"
	"context"
)

type RoomRepository interface {
	FindALl(ctx context.Context, queryPageCount utils.QueryPageLimit, filter *model.FilterDataRoomRequest) (*entity.RoomsDataAccessObject, error)
	Insert(ctx context.Context, input *entity.Room) (insertID int, err error)
	GetDashboard(ctx context.Context) (rooms []*model.RoomDashboards, err error)
	Update(ctx context.Context, input *entity.Room, roomID int64) (rowAffected int, err error)
	DeleteByID(ctx context.Context, roomID int64) (err error)
	FindOneByID(ctx context.Context, roomID int64) (room *entity.Room, err error)
}
