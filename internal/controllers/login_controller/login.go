package login_controller

import (
	"blog-system/internal/helper"
	"blog-system/internal/requests"
	"blog-system/internal/resources"
	"blog-system/internal/service/login_service"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type loginController struct {
	loginService login_service.LoginUser
	log          *logrus.Logger
}

func NewLogin(
	loginService login_service.LoginUser,
	log *logrus.Logger,
) Login {
	return &loginController{
		loginService: loginService,
		log:          log,
	}
}

func (l *loginController) Login(w http.ResponseWriter, r *http.Request) {
	var (
		payload  requests.LoginRequest
		response resources.Response
	)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	resultLogin, err := l.loginService.LoginUser(r.Context(), requests.LoginRequest{
		Email:    payload.Email,
		Password: payload.Password,
	})

	if err != nil {
		response.Code = resultLogin.Code
		response.Message = resultLogin.Message
		helper.ResponseJSON(w, resultLogin.Code, response)
		return
	}

	helper.ResponseJSON(w, http.StatusOK, resultLogin)
	return
}
