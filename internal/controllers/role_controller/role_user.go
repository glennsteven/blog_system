package role_controller

import (
	"blog-system/internal/helper"
	"blog-system/internal/requests"
	"blog-system/internal/resources"
	"blog-system/internal/service/role_service"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type roleController struct {
	roleUserService role_service.RoleUser
	log             *logrus.Logger
}

func NewRoleController(
	roleUserService role_service.RoleUser,
	log *logrus.Logger,
) Role {
	return &roleController{
		roleUserService: roleUserService,
		log:             log,
	}
}

func (ro *roleController) RoleUser(w http.ResponseWriter, r *http.Request) {
	var (
		payload  requests.RoleRequest
		response resources.Response
	)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		ro.log.Infof("decode request failed: %v", err)
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	result, err := ro.roleUserService.RoleUser(r.Context(), requests.RoleRequest{RoleName: payload.RoleName})
	if err != nil {
		ro.log.Infof("processing store role failed: %v", err)
		response.Code = http.StatusInternalServerError
		response.Message = err.Error()
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	helper.ResponseJSON(w, http.StatusCreated, result)
	return
}

func (ro *roleController) AssignRole(w http.ResponseWriter, r *http.Request) {
	var (
		payload  requests.AssignRoleRequest
		response resources.Response
	)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		ro.log.Infof("decode request failed: %v", err)
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	result, err := ro.roleUserService.AssignRole(r.Context(), requests.AssignRoleRequest{
		UserId: payload.UserId,
		RoleId: payload.RoleId,
	})

	if err != nil {
		ro.log.Infof("processing assign role failed: %v", err)
		response.Code = http.StatusInternalServerError
		response.Message = err.Error()
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	helper.ResponseJSON(w, http.StatusCreated, result)
	return
}
