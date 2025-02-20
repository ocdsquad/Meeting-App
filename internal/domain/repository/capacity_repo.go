package repository

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/pkg/utils"
	"context"
)

type CapacityRepository interface {
	GetOneByID(ctx context.Context, capacityID int64) (*entity.Capacity, error)
	FindAll(ctx context.Context, queryPagination utils.QueryPageLimit) (*entity.CapacityResultDataAccessObject, error)
}
