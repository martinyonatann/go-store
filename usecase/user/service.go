package user

import (
	"context"
	"errors"
	"strings"

	"github.com/martinyonatann/go-invoice/infrastructure/repository/user_repository"
	logger "github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

func New(r user_repository.Repository) *Feat {
	return &Feat{userRepository: r}
}

type Feat struct {
	userRepository user_repository.Repository
}

var ErrUserIdNotFound = errors.New("user_id not found")

func (x *Feat) CreateUser(ctx context.Context,
	request CreateUserRequest,
) (response CreateUserResponse, err error) {
	defer func() {
		logger.Err(err).Interface("userData", response).Msg("CreateUser_UseCase")
	}()

	newPassword, err := generatePassword(request.Password)
	if err != nil {
		return response, errors.New("failed generate password")
	}

	createUserPayload := user_repository.CreateUserRequest{
		FullName: request.FullName,
		Email:    request.Email,
		Password: newPassword,
	}

	userData, err := x.userRepository.CreateUser(ctx, createUserPayload)
	if err != nil {
		// if strings(err.Error())
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return response, errors.New("email already used")
		}

		return response, err
	}

	response = CreateUserResponse{
		ID:        userData.ID,
		FullName:  userData.FullName,
		Email:     userData.Email,
		CreatedAt: userData.CreatedAt,
	}

	return response, err
}

func (x *Feat) GetUser(
	ctx context.Context,
	request GetUserRequest,
) (response GetUserResponse, err error) {
	defer func() {
		logger.Err(err).Interface("userData", response).Msg("GetUser_UseCase")
	}()

	// validation
	{
		if request.UserID == 0 {
			return response, ErrUserIdNotFound
		}
	}

	userData, err := x.userRepository.GetUserById(ctx, user_repository.GetUserRequest{
		UserID: request.UserID,
	})
	if err != nil {
		return response, err
	}

	response = GetUserResponse{
		ID:        userData.ID,
		FullName:  userData.FullName,
		Email:     userData.Email,
		CreatedAt: userData.CreatedAt,
	}

	return response, err
}

func generatePassword(raw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
