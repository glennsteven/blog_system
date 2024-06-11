package config

import (
	"blog-system/pkg/viper"
	"fmt"
)

func Load() (*Configurations, error) {
	v, err := viper.GlobalConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	cfg := Configurations{
		AppPort: v.GetString("APP_PORT"),
		AppEnv:  v.GetString("APP_ENV"),
		AppName: v.GetString("APP_NAME"),
		Database: Database{
			Driver:   v.GetString("DB_DRIVER"),
			Username: v.GetString("DB_USERNAME"),
			Password: v.GetString("DB_PASSWORD"),
			Host:     v.GetString("DB_HOST"),
			Port:     v.GetInt32("DB_PORT"),
			DbName:   v.GetString("DB_NAME"),
		},
		Logger: Logger{
			Level: v.GetString("LOG_LEVEL"),
		},
		Jwt: Jwt{Key: v.GetString("JWT_KEY")},
	}

	return &cfg, nil
}

type Configurations struct {
	AppPort  string
	AppEnv   string
	AppName  string
	Database Database
	Logger   Logger
	Jwt      Jwt
}

type Database struct {
	Driver   string
	Username string
	Password string
	Host     string
	Port     int32
	DbName   string
}

type Logger struct {
	Level string
}

type Jwt struct {
	Key string
}
