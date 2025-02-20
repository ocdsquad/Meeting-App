package usecase

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/pkg/utils"
	"context"
)

type RoomTypeUseCase interface {
	FindAllRoomType(ctx context.Context, queryPageLimit utils.QueryPageLimit) (*entity.RoomTypeResultDataAccessObject, error)
}
