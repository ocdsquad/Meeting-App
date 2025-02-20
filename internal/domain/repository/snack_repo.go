package repository

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/pkg/utils"
	"context"
)

type SnackRepository interface {
	FindALl(ctx context.Context, queryPagination utils.QueryPageLimit) (*entity.SnackResultDataAccessObject, error)
	GetOneByID(ctx context.Context, snackID int64) (*entity.Snack, error)
}
