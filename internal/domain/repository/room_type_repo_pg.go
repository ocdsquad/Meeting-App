package repository

import (
	"E-Meeting/internal/domain/entity"
	reasonError "E-Meeting/pkg/reason"
	"E-Meeting/pkg/utils"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

type roomTypeRepo struct {
	DB *sqlx.DB
}

func NewRoomTypeRepository(db *sqlx.DB) RoomTypeRepository {
	return &roomTypeRepo{
		DB: db,
	}
}

func (repo *roomTypeRepo) GetOneByID(ctx context.Context, roomTypeID int64) (*entity.RoomType, error) {

	var roomType entity.RoomType
	query := fmt.Sprintf("SELECT id, name, is_active FROM room_types WHERE id = $1")

	err := repo.DB.GetContext(ctx, &roomType, query, roomTypeID)
	if err != nil {
		return nil, err
	}

	if roomType.ID == 0 {
		return nil, reasonError.ErrDataNotFound
	}

	return &roomType, nil
}

func (repo *roomTypeRepo) FindAll(_ context.Context, queryPagination utils.QueryPageLimit) (*entity.RoomTypeResultDataAccessObject, error) {

	queryCount := utils.QueryCount{
		TotalData:    0,
		TotalContent: 0,
	}

	roomTypeDataAccessObject := entity.RoomTypeResultDataAccessObject{
		RoomTypes:  nil,
		QueryCount: queryCount,
	}

	queryLimit := ""
	if (queryPagination.Page > 0) && (queryPagination.Limit > 0) {

		queryPagination.Page = (queryPagination.Page - 1) * queryPagination.Limit

		queryLimit = fmt.Sprintf("LIMIT %d OFFSET %d", queryPagination.Limit, queryPagination.Page)
	}

	queryOrderSort := fmt.Sprintf("ORDER BY %s %s", queryPagination.OrderBy, queryPagination.SortBy)

	var roomTypes []entity.RoomType
	query := fmt.Sprintf("SELECT id, name, is_active, created_at, updated_at, created_by, updated_by FROM public.room_types %s %s", queryOrderSort, queryLimit)

	err := repo.DB.Select(&roomTypes, query)
	if err != nil {
		return nil, err
	}

	if len(roomTypes) < 1 {
		return nil, reasonError.ErrDataNotFound
	}

	// Query for total count
	sqlCount := fmt.Sprintf("SELECT COUNT(*) FROM public.category_snacks")

	var totalCount int
	err = repo.DB.Get(&totalCount, sqlCount)

	if err != nil {
		log.Println(fmt.Sprintf("message : error in query count | repo : snack_repo_pg | error : %s", err))
		return nil, err
	}
	roomTypeDataAccessObject.RoomTypes = roomTypes
	roomTypeDataAccessObject.QueryCount.TotalData = totalCount
	roomTypeDataAccessObject.QueryCount.TotalContent = len(roomTypes)

	return &roomTypeDataAccessObject, nil
}
