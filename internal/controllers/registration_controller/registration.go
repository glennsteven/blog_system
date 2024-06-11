package registration_controller

import (
	"blog-system/internal/helper"
	"blog-system/internal/requests"
	"blog-system/internal/resources"
	"blog-system/internal/service/registration_service"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type registrationController struct {
	regisUserService registration_service.UserRegistration
	log              *logrus.Logger
}

func NewRegistrationController(
	regisUserService registration_service.UserRegistration,
	log *logrus.Logger,
) Registration {
	return &registrationController{
		regisUserService: regisUserService,
		log:              log,
	}
}

func (re *registrationController) UserRegistration(w http.ResponseWriter, r *http.Request) {
	var (
		payload  requests.RegisterUserRequest
		response resources.Response
	)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		re.log.Infof("decode request failed: %v", err)
		response.Code = http.StatusBadRequest
		response.Message = err.Error()
		helper.ResponseJSON(w, http.StatusBadRequest, response)
		return
	}

	isValid := helper.CheckUniquePassword(payload.Password)
	if !isValid {
		response.Code = http.StatusUnprocessableEntity
		response.Message = "password must contain uppercase, lowercase letters, numbers, special character and a minimum of 8 characters"
		helper.ResponseJSON(w, http.StatusUnprocessableEntity, response)
		return
	}

	resultHashing, err := hashingPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		re.log.Infof("password hashing got error: %v", err)
		response.Code = http.StatusBadRequest
		response.Message = "Internal Server Error"
		helper.ResponseJSON(w, http.StatusUnprocessableEntity, response)
		return
	}

	result, err := re.regisUserService.UserRegistration(r.Context(), requests.RegisterUserRequest{
		FullName: payload.FullName,
		Email:    payload.Email,
		Address:  payload.Address,
	}, string(resultHashing))

	if err != nil {
		re.log.Infof("processing user failed: %v", err)
		response.Code = http.StatusInternalServerError
		response.Message = err.Error()
		helper.ResponseJSON(w, http.StatusInternalServerError, response)
		return
	}

	helper.ResponseJSON(w, http.StatusCreated, result)
	return
}
