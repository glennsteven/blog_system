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
			Driver:            v.GetString("DB_DRIVER"),
			Name:              v.GetString("DB_NAME"),
			Host:              v.GetString("DB_HOST"),
			User:              v.GetString("DB_USER"),
			Password:          v.GetString("DB_PASSWORD"),
			Timezone:          v.GetString("DB_TIMEZONE"),
			SSLEnabled:        v.GetBool("DB_SSL_ENABLED"),
			MaxOpenConnection: v.GetInt("DB_MAX_OPEN_CONN"),
			MaxIdleConnection: v.GetInt("DB_MAX_IDLE_CONN"),
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
	Driver            string
	Name              string
	Host              string
	User              string
	Password          string
	Timezone          string
	SSLEnabled        bool
	MaxOpenConnection int
	MaxIdleConnection int
}

type Logger struct {
	Level string
}

type Jwt struct {
	Key string
}
