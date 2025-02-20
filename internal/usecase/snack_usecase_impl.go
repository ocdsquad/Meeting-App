package usecase

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/internal/domain/repository"
	"E-Meeting/pkg/utils"
	"context"
	"fmt"
	"log"
)

type snackUseCase struct {
	repo repository.SnackRepository
}

func NewSnackUseCase(repo repository.SnackRepository) SnackUseCase {
	return &snackUseCase{repo: repo}
}

func (u *snackUseCase) FindAllSnack(ctx context.Context, queryPageLimit utils.QueryPageLimit) (*entity.SnackResultDataAccessObject, error) {

	snacks, err := u.repo.FindALl(ctx, queryPageLimit)
	if err != nil {
		log.Println(fmt.Sprintf("message : error in service | service : snack_usecase_impl | error : %s", err))
		return nil, err
	}

	return snacks, nil

}
