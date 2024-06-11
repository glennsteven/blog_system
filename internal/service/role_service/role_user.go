package role_service

import (
	"blog-system/internal/entities"
	"blog-system/internal/repositories"
	"blog-system/internal/requests"
	"blog-system/internal/resources"
	"context"
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
	role, err := r.repoRole.FindRole(ctx, payload.RoleName)
	if err != nil {
		r.log.Infof("finding role name: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}

	if role != nil {
		r.log.Infof("data role existing: %v", err)
		return resources.Response{
			Code:    http.StatusBadRequest,
			Message: "role is existing",
		}, err
	}

	saveRole, err := r.repoRole.Store(ctx, entities.Role{
		Name:      payload.RoleName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		r.log.Infof("failed process insert new role: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
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
		r.log.Infof("finding user id: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}

	if user == nil {
		r.log.Infof("user is not exist: %v", err)
		return resources.Response{
			Code:    http.StatusBadRequest,
			Message: "user not found",
		}, err
	}

	role, err := r.repoRole.FindRoleId(ctx, payload.RoleId)
	if err != nil {
		r.log.Infof("finding role id: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}

	if role == nil {
		r.log.Infof("role is not exist: %v", err)
		return resources.Response{
			Code:    http.StatusBadRequest,
			Message: "role not found",
		}, err
	}

	roleAssigned, err := r.repoRoleUser.FindUserRole(ctx, entities.UserRole{
		UserId: user.Id,
		RoleId: role.Id,
	})
	if err != nil {
		r.log.Infof("finding role user: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}

	if roleAssigned != nil {
		r.log.Infof("data role user existing: %v", err)
		return resources.Response{
			Code:    http.StatusBadRequest,
			Message: "role user duplicate",
		}, err
	}

	saveRoleUser, err := r.repoRoleUser.Store(ctx, entities.UserRole{
		UserId:    payload.UserId,
		RoleId:    payload.RoleId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		r.log.Infof("failed process insert new role user: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}

	return resources.Response{
		Code:    http.StatusCreated,
		Message: "successfully assign role permission",
		Data:    saveRoleUser,
	}, nil
}
