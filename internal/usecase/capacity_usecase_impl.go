package usecase

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/internal/domain/repository"
	"E-Meeting/pkg/utils"
	"context"
	"fmt"
	"log"
)

type capacityUseCase struct {
	repo repository.CapacityRepository
}

func NewCapacityUseCase(repo repository.CapacityRepository) CapacityUseCase {
	return &capacityUseCase{repo: repo}
}

func (u *capacityUseCase) FindAllCapacity(ctx context.Context, queryPagination utils.QueryPageLimit) (*entity.CapacityResultDataAccessObject, error) {

	result, err := u.repo.FindAll(ctx, queryPagination)
	if err != nil {
		log.Println(fmt.Sprintf("message : error in service | service : capacity_usecase_impl | error : %s", err))
		return nil, err
	}

	return result, nil

}
