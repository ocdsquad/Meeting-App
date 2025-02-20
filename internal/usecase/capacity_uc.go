package usecase

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/pkg/utils"
	"context"
)

type CapacityUseCase interface {
	FindAllCapacity(ctx context.Context, queryPagination utils.QueryPageLimit) (*entity.CapacityResultDataAccessObject, error)
}
