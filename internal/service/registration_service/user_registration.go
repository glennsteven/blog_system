package registration_service

import (
	"blog-system/internal/entities"
	"blog-system/internal/repositories"
	"blog-system/internal/requests"
	"blog-system/internal/resources"
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type userRegistrationService struct {
	repoUser repositories.UserRepositories
	log      *logrus.Logger
}

func NewUserRegistrationService(repoUser repositories.UserRepositories, log *logrus.Logger) UserRegistration {
	return &userRegistrationService{repoUser: repoUser, log: log}
}

func (u *userRegistrationService) UserRegistration(ctx context.Context, payload requests.RegisterUserRequest, hashingPassword string) (resources.Response, error) {
	saveUser, err := u.repoUser.Store(ctx, entities.User{
		Email:     payload.Email,
		FullName:  payload.FullName,
		Password:  hashingPassword,
		Address:   payload.Address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		errLog := errors.Wrap(err, "error creating user")
		u.log.Error(errLog)
		switch errors.Cause(err) {
		case entities.ErrUserAlreadyExist:
			return resources.Response{Code: http.StatusUnprocessableEntity, Message: err.Error()}, err
		default:
			return resources.Response{Code: http.StatusInternalServerError}, err
		}
	}

	return resources.Response{
		Code:    http.StatusCreated,
		Message: "successfully register",
		Data: resources.RegistrationResource{
			FullName:  saveUser.FullName,
			Email:     saveUser.Email,
			CreatedAt: saveUser.CreatedAt.Format(`2006-01-02 15:04:05`),
			UpdatedAt: saveUser.UpdatedAt.Format(`2006-01-02 15:04:05`),
		},
	}, nil
}
