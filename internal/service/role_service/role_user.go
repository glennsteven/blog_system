package role_service

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

type roleUserService struct {
	repoRole     repositories.RoleRepositories
	repoRoleUser repositories.RoleUser
	repoUser     repositories.UserRepositories
	log          *logrus.Logger
}

func NewRoleUserService(
	repoRole repositories.RoleRepositories,
	repoRoleUser repositories.RoleUser,
	repoUser repositories.UserRepositories,
	log *logrus.Logger,
) RoleUser {
	return &roleUserService{
		repoRole:     repoRole,
		repoRoleUser: repoRoleUser,
		repoUser:     repoUser,
		log:          log,
	}
}

func (r *roleUserService) RoleUser(ctx context.Context, payload requests.RoleRequest) (resources.Response, error) {
	saveRole, err := r.repoRole.Store(ctx, entities.Role{
		Name:      payload.RoleName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		errLog := errors.Wrap(err, "error creating role")
		r.log.Error(errLog)
		switch errors.Cause(err) {
		case entities.ErrRoleAlreadyExist:
			return resources.Response{Code: http.StatusUnprocessableEntity, Message: err.Error()}, err
		default:
			return resources.Response{Code: http.StatusInternalServerError}, err
		}
	}

	return resources.Response{
		Code:    http.StatusCreated,
		Message: "successfully register",
		Data:    saveRole,
	}, nil
}

func (r *roleUserService) AssignRole(ctx context.Context, payload requests.AssignRoleRequest) (resources.Response, error) {
	user, err := r.repoUser.FindUserId(ctx, payload.UserId)
	if err != nil {
		errLog := errors.Wrap(err, "error finding user")
		r.log.Error(errLog)
		switch errors.Cause(err) {
		case entities.ErrUserNotFound:
			return resources.Response{Code: http.StatusUnprocessableEntity, Message: err.Error()}, err
		default:
			return resources.Response{Code: http.StatusInternalServerError}, err
		}
	}

	role, err := r.repoRole.FindRoleId(ctx, payload.RoleId)
	if err != nil {
		errLog := errors.Wrap(err, "error finding role")
		r.log.Error(errLog)
		switch errors.Cause(err) {
		case entities.ErrRoleNotFound:
			return resources.Response{Code: http.StatusUnprocessableEntity, Message: err.Error()}, err
		default:
			return resources.Response{Code: http.StatusInternalServerError}, err
		}
	}

	_, err = r.repoRoleUser.FindUserRole(ctx, entities.UserRole{
		UserId: user.Id,
		RoleId: role.Id,
	})

	if err != nil {
		errLog := errors.Wrap(err, "error finding role")
		r.log.Error(errLog)
		switch errors.Cause(err) {
		case entities.ErrRoleNotFound:
			return resources.Response{Code: http.StatusUnprocessableEntity, Message: err.Error()}, err
		default:
			return resources.Response{Code: http.StatusInternalServerError}, err
		}
	}

	//if roleAssigned != nil {
	//	r.log.Infof("data role user existing: %v", err)
	//	return resources.Response{
	//		Code:    http.StatusUnprocessableEntity,
	//		Message: "role user duplicate",
	//	}, err
	//}

	saveRoleUser, err := r.repoRoleUser.Store(ctx, entities.UserRole{
		UserId:    payload.UserId,
		RoleId:    payload.RoleId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		errLog := errors.Wrap(err, "error creating user role")
		r.log.Error(errLog)
		switch errors.Cause(err) {
		case entities.ErrUserRoleAlreadyExist:
			return resources.Response{Code: http.StatusUnprocessableEntity, Message: err.Error()}, err
		default:
			return resources.Response{Code: http.StatusInternalServerError}, err
		}
	}

	return resources.Response{
		Code:    http.StatusCreated,
		Message: "successfully assign role permission",
		Data:    saveRoleUser,
	}, nil
}
