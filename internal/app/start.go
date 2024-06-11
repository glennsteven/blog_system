package app

import (
	"blog-system/internal/config"
	"blog-system/internal/router"
	"blog-system/pkg/logger"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Start() error {
	r := mux.NewRouter()
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logs, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		return fmt.Errorf("cannot load logger: %w", err)
	}

	logrus.Infof("%s already running with port %s", cfg.AppName, cfg.AppPort)

	defer func() {
		if r := recover(); r != nil {
			logrus.Println("Recovered blog service. Error:\n", r)
		}
	}()
	router.Router(r, cfg, logs)

	logs.Info(http.ListenAndServe(":"+cfg.AppPort, r))

	return nil
}
