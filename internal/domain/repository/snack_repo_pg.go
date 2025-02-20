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

type snackRepo struct {
	DB *sqlx.DB
}

func NewSnackRepository(db *sqlx.DB) SnackRepository {
	return &snackRepo{
		DB: db,
	}
}

func (repo *snackRepo) FindALl(_ context.Context, queryPagination utils.QueryPageLimit) (*entity.SnackResultDataAccessObject, error) {

	queryCount := utils.QueryCount{
		TotalData:    0,
		TotalContent: 0,
	}

	snackDataAccessObject := entity.SnackResultDataAccessObject{
		Snacks:     nil,
		QueryCount: queryCount,
	}

	queryLimit := ""
	if (queryPagination.Page > 0) && (queryPagination.Limit > 0) {

		queryPagination.Page = (queryPagination.Page - 1) * queryPagination.Limit

		queryLimit = fmt.Sprintf("LIMIT %d OFFSET %d", queryPagination.Limit, queryPagination.Page)
	}

	queryOrderSort := fmt.Sprintf("ORDER BY %s %s", queryPagination.OrderBy, queryPagination.SortBy)

	var snacks []entity.Snack
	query := fmt.Sprintf("SELECT id, name, price, currency, uom, is_active, created_at, updated_at, created_by, updated_by FROM public.category_snacks %s %s", queryOrderSort, queryLimit)

	err := repo.DB.Select(&snacks, query)
	if err != nil {
		return nil, err
	}

	if len(snacks) < 1 {
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
	snackDataAccessObject.Snacks = snacks
	snackDataAccessObject.QueryCount.TotalData = totalCount
	snackDataAccessObject.QueryCount.TotalContent = len(snacks)

	return &snackDataAccessObject, nil
}
func (repo *snackRepo) GetOneByID(ctx context.Context, snackID int64) (*entity.Snack, error) {

	var snack entity.Snack

	query := fmt.Sprintf("SELECT id, name, price, currency, uom, is_active, created_at, updated_at, created_by, updated_by FROM category_snacks WHERE id = $1")

	err := repo.DB.GetContext(ctx, &snack, query, snackID)
	if err != nil {
		return nil, err
	}

	if snack.ID == 0 {
		return nil, reasonError.ErrDataNotFound
	}

	return &snack, nil
}
