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

type capacityRepo struct {
	DB *sqlx.DB
}

func NewCapacityRepository(db *sqlx.DB) CapacityRepository {
	return &capacityRepo{
		DB: db,
	}
}

func (repo *capacityRepo) GetOneByID(ctx context.Context, capacityID int64) (*entity.Capacity, error) {

	var capacity entity.Capacity
	query := fmt.Sprintf("SELECT id, value_minimum, value_maximum, uom, is_active FROM capacities WHERE id = $1")

	err := repo.DB.GetContext(ctx, &capacity, query, capacityID)
	if err != nil {
		return nil, err
	}

	if capacity.ID == 0 {
		return nil, reasonError.ErrDataNotFound
	}

	return &capacity, nil
}

func (repo *capacityRepo) FindAll(_ context.Context, queryPagination utils.QueryPageLimit) (*entity.CapacityResultDataAccessObject, error) {

	queryCount := utils.QueryCount{
		TotalData:    0,
		TotalContent: 0,
	}

	resultDataAccessObject := entity.CapacityResultDataAccessObject{
		Capacity:   nil,
		QueryCount: queryCount,
	}

	queryLimit := ""
	if (queryPagination.Page > 0) && (queryPagination.Limit > 0) {

		queryPagination.Page = (queryPagination.Page - 1) * queryPagination.Limit

		queryLimit = fmt.Sprintf("LIMIT %d OFFSET %d", queryPagination.Limit, queryPagination.Page)
	}

	queryOrderSort := fmt.Sprintf("ORDER BY %s %s", queryPagination.OrderBy, queryPagination.SortBy)

	var capacities []entity.Capacity
	query := fmt.Sprintf("SELECT id, value_minimum, value_maximum, uom, is_active, created_at, updated_at, created_by, updated_by FROM public.capacities %s %s", queryOrderSort, queryLimit)

	err := repo.DB.Select(&capacities, query)
	if err != nil {
		return nil, err
	}

	if len(capacities) < 1 {
		return nil, reasonError.ErrDataNotFound
	}

	// Query for total count
	sqlCount := fmt.Sprintf("SELECT COUNT(*) FROM public.capacities")

	var totalCount int
	err = repo.DB.Get(&totalCount, sqlCount)

	if err != nil {
		log.Println(fmt.Sprintf("message : error in query count | repo : capacity_repo_pg | error : %s", err))
		return nil, err
	}
	resultDataAccessObject.Capacity = capacities
	resultDataAccessObject.QueryCount.TotalData = totalCount
	resultDataAccessObject.QueryCount.TotalContent = len(capacities)

	return &resultDataAccessObject, nil
}
