package repository

import (
	"E-Meeting/internal/domain/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UserRepo struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &UserRepo{
		DB: db,
	}
}

func (repo *UserRepo) Save(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	errorQuery := repo.DB.QueryRowx(query, user.Username, user.Email, user.Password).Scan(&user.ID)
	if errorQuery != nil {
		return fmt.Errorf("could not create user: %v", errorQuery)
	}
	return nil
}

func (repo *UserRepo) Login(ctx context.Context, username, _ string) (*entity.User, error) {
	query := `SELECT id, username, password,email, is_admin, language FROM users WHERE username = $1`
	row := repo.DB.QueryRowContext(ctx, query, username)

	var user entity.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.IsAdmin, &user.Language)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &user, errors.New("user not found")
		}

		return &user, fmt.Errorf("could not login user: %v", err)
	}
	return &user, nil
}

func (repo *UserRepo) GetByID(ctx context.Context, id int) (*entity.User, error) {
	query := `SELECT id, username, email, is_active, is_admin, language, avatar_url FROM users WHERE id = $1`
	row, err := repo.DB.QueryContext(ctx, query, id)

	if err != nil {
		panic(err)
	}
	var user entity.User
	if !row.Next() {
		return &user, errors.New("user not found")
	}

	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.IsActive, &user.IsAdmin, &user.Language, &user.AvatarUrl)
	if err != nil {
		return nil, fmt.Errorf("could not get user: %v", err)
	}
	return &user, nil

}

func (repo *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `SELECT id, email, password FROM users WHERE email = $1`
	row, err := repo.DB.QueryContext(ctx, query, email)

	if err != nil {
		panic(err)
	}
	var user entity.User
	if !row.Next() {
		return &user, errors.New("user not found")
	}

	err = row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("could not get user: %v", err)
	}
	return &user, nil

}

func (repo *UserRepo) ResetPassword(_ context.Context, user *entity.User) error {

	query := `UPDATE users SET password = $1 WHERE email = $2`
	_, err := repo.DB.Exec(query, user.Password, user.Email)
	if err != nil {
		return fmt.Errorf("could not reset password: %v", err)
	}
	return nil
}

func (repo *UserRepo) Update(_ context.Context, user *entity.User) error {
	query := `UPDATE users SET username = $1, email = $2, language = $3, avatar_url = $4  WHERE id = $5`
	_, err := repo.DB.Exec(query, user.Username, user.Email, user.Language, user.AvatarUrl, user.ID)
	if err != nil {
		return fmt.Errorf("could not update user: %v", err)
	}
	return nil
}

// func (repo *UserRepo) Delete(ctx context.Context, id int) error {
// 	query := `DELETE FROM users WHERE id = $1`
// 	_, err := repo.DB.Exec(query, id)
// 	if err != nil {
// 		return fmt.Errorf("could not delete user: %v", err)
// 	}
// 	return nil
// }
