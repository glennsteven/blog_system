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
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Router(r *mux.Router, cfg *config.Configurations, log *logrus.Logger) error {
	r.HandleFunc("/",
		checkHealth.HealthCheck,
	).Methods(http.MethodGet)

	dbDriver := postgres.DriverType(cfg.Database.Driver)
	db, err := postgres.DatabaseOpen(dbDriver, cfg.Database)
	if err != nil {
		return errors.Wrap(err, "cannot connect to third party database")
	}

	userRepository := repositories.NewUsers(db, log)
	roleRepository := repositories.NewRoles(db, log)
	roleUserRepository := repositories.NewUserRole(db, log)
	//postRepository := repositories.NewPost(db, log)
	//tagRepository := repositories.NewTags(db, log)
	//postTagRepository := repositories.NewPostTag(db, log)

	userRegistrationService := registration_service.NewUserRegistrationService(userRepository, log)
	roleService := role_service.NewRoleUserService(roleRepository, roleUserRepository, userRepository, log)
	//loginService := login_service.NewLoginUserService(userRepository, roleUserRepository, log, cfg)
	//postService := post_service.NewPostService(postRepository, tagRepository, postTagRepository, log)

	userRegistrationController := registration_controller.NewRegistrationController(userRegistrationService, log)
	roleController := role_controller.NewRoleController(roleService, log)
	//loginController := login_controller.NewLogin(loginService, log)
	//postController := post_controller.NewPostController(postService, log, cfg)

	sub := r.PathPrefix("/api").Subrouter()

	sub.HandleFunc("/register",
		userRegistrationController.UserRegistration,
	).Methods(http.MethodPost)

	sub.HandleFunc("/role",
		roleController.RoleUser,
	).Methods(http.MethodPost)

	sub.HandleFunc("/role/assign",
		roleController.AssignRole,
	).Methods(http.MethodPost)

	//sub.HandleFunc("/login",
	//	loginController.Login,
	//).Methods(http.MethodPost)
	//
	//subs := r.PathPrefix("/api/v1").Subrouter()
	//subs.HandleFunc("/posts",
	//	postController.Post,
	//).Methods(http.MethodPost)
	//
	//subs.HandleFunc("/posts/{post_id}",
	//	postController.UpdatePost,
	//).Methods(http.MethodPut)
	//
	//subs.HandleFunc("/posts/{post_id}",
	//	postController.GetPost,
	//).Methods(http.MethodGet)
	//
	//subs.HandleFunc("/posts",
	//	postController.GetPostFromTag,
	//).Methods(http.MethodGet)
	//
	//subs.HandleFunc("/posts/{post_id}",
	//	postController.DestroyPost,
	//).Methods(http.MethodDelete)
	//
	//subs.Use(middlewares.AuthMiddleware)
	return nil
}
