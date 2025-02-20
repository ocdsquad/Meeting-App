package usecase

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/internal/domain/repository"
	"E-Meeting/pkg/utils"
	"context"
	"fmt"
	"log"
)

type roomTypeUseCase struct {
	repo repository.RoomTypeRepository
}

func NewRoomTypeUseCase(repo repository.RoomTypeRepository) RoomTypeUseCase {
	return &roomTypeUseCase{repo: repo}
}

func (u *roomTypeUseCase) FindAllRoomType(ctx context.Context, queryPageLimit utils.QueryPageLimit) (*entity.RoomTypeResultDataAccessObject, error) {

	rooms, err := u.repo.FindAll(ctx, queryPageLimit)
	if err != nil {
		log.Println(fmt.Sprintf("message : error in service | service : room_type_usecase_impl | error : %s", err))
		return nil, err
	}

	return rooms, nil

}
