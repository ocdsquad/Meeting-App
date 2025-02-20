package usecase

import (
	"E-Meeting/pkg/mailer"
	"E-Meeting/presenter/model"
	"context"
	"mime/multipart"
)

type UserUseCase interface {
	Save(ctx context.Context, request *model.UserCreateRequest) error
	Login(ctx context.Context, request *model.UserLoginRequest) (*model.UserLoginResponse, error)
	Update(ctx context.Context, id int, request *model.UserUpdateProfileRequest, file *multipart.FileHeader) error
	ResetPassword(ctx context.Context, request *model.UserResetPasswordRequest, otp string) error
	ForgotPassword(ctx context.Context, request *model.UserForgotPasswordRequest, mailer *mailer.Mailer) error
	GetByID(ctx context.Context, id int) (*model.UserGetProfileResponse, error)
}
