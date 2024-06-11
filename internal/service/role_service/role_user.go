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
	repoRole repositories.RoleRepositories
	log      *logrus.Logger
}

func NewRoleUserService(repoRole repositories.RoleRepositories, log *logrus.Logger) RoleUser {
	return &roleUserService{repoRole: repoRole, log: log}
}

func (r *roleUserService) RoleUser(ctx context.Context, payload requests.RoleRequest) (resources.Response, error) {
	role, err := r.repoRole.FindRole(ctx, payload.RoleName)
	if err != nil {
		r.log.Infof("finding role_controller name: %v", err)
		return resources.Response{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
		}, err
	}

	if role != nil {
		r.log.Infof("data role_controller existing: %v", err)
		return resources.Response{
			Code:    http.StatusBadRequest,
			Message: "role_controller is existing",
		}, err
	}

	saveRole, err := r.repoRole.Store(ctx, entities.Role{
		Name:      payload.RoleName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		r.log.Infof("failed process insert new role_controller: %v", err)
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
