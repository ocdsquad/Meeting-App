package usecase

import (
	"E-Meeting/internal/domain/entity"
	"E-Meeting/internal/domain/repository"
	"E-Meeting/pkg/cache"
	"E-Meeting/pkg/common"
	"E-Meeting/pkg/helper"
	"E-Meeting/pkg/mailer"
	"E-Meeting/pkg/reason"
	"E-Meeting/pkg/utils"
	"E-Meeting/presenter/model"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/guregu/null"
)

type userUseCase struct {
	repo     repository.UserRepository
	Validate *validator.Validate
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return &userUseCase{
		repo:     repo,
		Validate: validator.New(),
	}
}

func (u *userUseCase) Save(ctx context.Context, request *model.UserCreateRequest) error {
	log.Printf("request: %v\n", request)
	err := u.Validate.Struct(request)
	if err != nil {
		return err
	}
	hashedPassword, err := helper.HashPassword(request.Password)
	if err != nil {
		return err
	}
	request.Password = hashedPassword
	user := entity.User{
		Email:    request.Email,
		Username: request.Username,
		Password: request.Password,
	}

	err = u.repo.Save(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUseCase) Login(ctx context.Context, request *model.UserLoginRequest) (*model.UserLoginResponse, error) {
	err := u.Validate.Struct(request)
	if err != nil {
		return nil, err
	}
	user, err := u.repo.Login(ctx, request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	err = helper.VerifyPassword(user.Password, request.Password)
	if err != nil {
		log.Printf("error verify password: %v\n", err)
		return nil, err
	}
	token, err := helper.GenerateJWT(user.ID, user.IsAdmin, user.Language, user.Email)
	if err != nil {
		log.Printf("error generate token: %v\n", err)
		return nil, err
	}

	response := &model.UserLoginResponse{
		Token: token,
	}

	return response, nil
}

func (u *userUseCase) GetByID(ctx context.Context, id int) (*model.UserGetProfileResponse, error) {
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.UserGetProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		IsActive:  user.IsActive,
		IsAdmin:   user.IsAdmin,
		Language:  user.Language,
		AvatarUrl: user.AvatarUrl,
	}, nil
}

func (u *userUseCase) ResetPassword(ctx context.Context, request *model.UserResetPasswordRequest, otp string) error {
	err := u.Validate.Struct(request)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("user_%s", otp)
	log.Printf("cache key: %s\n", cacheKey)
	data, _ := cache.MyCache.Get(cacheKey)
	log.Printf("data: %v\n", data)

	user, err := u.repo.GetByEmail(ctx, data.(common.UserOTP).Email)

	if err != nil {
		return err
	}

	hashedPassword, err := helper.HashPassword(request.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	err = u.repo.ResetPassword(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUseCase) ForgotPassword(ctx context.Context, request *model.UserForgotPasswordRequest, mailer *mailer.Mailer) error {
	err := u.Validate.Struct(request)
	if err != nil {
		return err
	}
	user, err := u.repo.GetByEmail(ctx, request.Email)
	if err != nil {
		return err
	}

	otp, err := helper.GenerateRandomID(6)
	if err != nil {
		return err
	}
	data := common.UserOTP{
		Email: user.Email,
		OTP:   otp,
	}
	cacheKey := fmt.Sprintf("user_%s", data.OTP)
	log.Printf("cache key: %s\n", cacheKey)
	cache.MyCache.Set(cacheKey, data, 5*time.Minute)

	err = mailer.SendMail(user.Email, "Reset Password", fmt.Sprintf("Your OTP: %s", otp))
	if err != nil {
		return fmt.Errorf("error send email: %v", err)
	}
	return nil
}

func (u *userUseCase) Update(ctx context.Context, id int, request *model.UserUpdateProfileRequest, file *multipart.FileHeader) error {
	err := u.Validate.Struct(request)
	if err != nil {
		return err
	}
	user, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	user.Username = request.Username
	user.Email = request.Email
	user.Language = request.Language

	isFile := false

	if file != nil {
		isFile = true
	}

	if isFile {

		fileURL, err := utils.SaveFile(file, "user")
		if err != nil {
			log.Println(fmt.Sprintf("message : capacity is not active | service : user_usecase_impl | validate : capacity_is_active"))
			return reason.ErrFailedInsertData
		}

		user.AvatarUrl = null.StringFrom(fileURL)
	}
	err = u.repo.Update(ctx, user)
	if err != nil {
		return err
	}
	return nil
}
