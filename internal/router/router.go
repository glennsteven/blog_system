package router

import (
	checkHealth "blog-system"
	"blog-system/internal/config"
	"blog-system/internal/controllers/registration_controller"
	"blog-system/internal/controllers/role_controller"
	"blog-system/internal/repositories"
	"blog-system/internal/service/registration_service"
	"blog-system/internal/service/role_service"
	"blog-system/pkg/database/postgres"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Router(r *mux.Router, cfg *config.Configurations, log *logrus.Logger) {
	r.HandleFunc("/",
		checkHealth.HealthCheck,
	).Methods(http.MethodGet)

	db := postgres.NewDatabase(cfg.Database, log)

	userRepository := repositories.NewUsers(db, log)
	roleRepository := repositories.NewRoles(db, log)

	userRegistrationService := registration_service.NewUserRegistrationService(userRepository, log)
	roleService := role_service.NewRoleUserService(roleRepository, log)

	userRegistrationController := registration_controller.NewRegistrationController(userRegistrationService, log)
	roleController := role_controller.NewRoleController(roleService, log)

	r.HandleFunc("/api/register",
		userRegistrationController.UserRegistration,
	).Methods(http.MethodPost)

	r.HandleFunc("/api/role",
		roleController.RoleUser,
	).Methods(http.MethodPost)

}
