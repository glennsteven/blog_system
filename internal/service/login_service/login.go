package login_service

import (
	"blog-system/internal/config"
	"blog-system/internal/repositories"
	"blog-system/internal/requests"
	"blog-system/internal/resources"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type JWTClaim struct {
	Id       int64
	FullName string
	RoleId   int64
	jwt.RegisteredClaims
}

type loginUserService struct {
	repoUser     repositories.UserRepositories
	repoRoleUser repositories.RoleUser
	log          *logrus.Logger
	cfg          *config.Configurations
}

func NewLoginUserService(
	repoUser repositories.UserRepositories,
	repoRoleUser repositories.RoleUser,
	log *logrus.Logger,
	cfg *config.Configurations,
) LoginUser {
	return &loginUserService{
		repoUser:     repoUser,
		repoRoleUser: repoRoleUser,
		log:          log,
		cfg:          cfg,
	}
}

func (l *loginUserService) LoginUser(ctx context.Context, payload requests.LoginRequest) (resources.Response, error) {
	findUser, err := l.repoUser.FindUser(ctx, payload.Email)
	if err != nil {
		l.log.Infof("finding user login: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}

	if findUser == nil {
		return resources.Response{
			Code:    http.StatusUnauthorized,
			Message: "invalid email or password",
		}, err
	}

	role, err := l.repoRoleUser.FindUserIdRole(ctx, findUser.Id)
	if err != nil {
		l.log.Infof("find role user: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}

	if role == nil {
		return resources.Response{
			Code:    http.StatusBadRequest,
			Message: "Please assign role first",
		}, err
	}

	check := ComparePasswords(findUser.Password, []byte(payload.Password))
	if !check {
		l.log.Info("compare hash password failed")
		return resources.Response{
			Code:    http.StatusUnauthorized,
			Message: "invalid email or password",
		}, err
	}

	exp := time.Now().Add(time.Hour * 15)

	claims := &JWTClaim{
		Id:       findUser.Id,
		FullName: findUser.FullName,
		RoleId:   role.RoleId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "xyz-issuer",
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	tokenAlgorithm := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenAlgorithm.SignedString([]byte(l.cfg.Jwt.Key))
	if err != nil {
		l.log.Infof("algoritm token: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}

	return resources.Response{
		Code:    http.StatusCreated,
		Message: "Token created",
		Data:    map[string]string{"token": token},
	}, nil
}
