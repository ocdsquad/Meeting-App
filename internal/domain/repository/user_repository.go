package repository

import (
	"E-Meeting/internal/domain/entity"
	"context"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	Login(ctx context.Context, username string, password string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	// Delete(ctx context.Context, id int) error
	ResetPassword(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}
