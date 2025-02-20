package usecase

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/pkg/utils"
	"context"
)

type SnackUseCase interface {
	FindAllSnack(ctx context.Context, queryPageLimit utils.QueryPageLimit) (*entity.SnackResultDataAccessObject, error)
}
