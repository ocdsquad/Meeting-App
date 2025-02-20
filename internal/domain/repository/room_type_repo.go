package repository

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/pkg/utils"
	"context"
)

type RoomTypeRepository interface {
	GetOneByID(ctx context.Context, roomTypeID int64) (*entity.RoomType, error)
	FindAll(ctx context.Context, queryPagination utils.QueryPageLimit) (*entity.RoomTypeResultDataAccessObject, error)
}
