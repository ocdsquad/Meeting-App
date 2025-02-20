package usecase

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/pkg/utils"
	"E-Meeting/presenter/model"
	"context"
	"mime/multipart"
)

type RoomUseCase interface {
	FindAllRoom(ctx context.Context, queryPageLimit utils.QueryPageLimit, filter *model.FilterDataRoomRequest) (*entity.RoomsDataAccessObject, error)
	Insert(ctx context.Context, input model.RoomRequest, file *multipart.FileHeader) error
	GetDashboard(ctx context.Context) ([]*model.RoomDashboards, error)
	UpdateOneByID(ctx context.Context, input model.RoomRequest, file *multipart.FileHeader, roomID int64) error
	DeleteOneByID(ctx context.Context, roomID int64) error
	FindOneByID(ctx context.Context, roomID int64) (*entity.Room, error)
}
