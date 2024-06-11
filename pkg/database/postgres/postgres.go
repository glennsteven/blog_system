package postgres

import (
	"blog-system/internal/config"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
)

func NewDatabase(cfg config.Database, log *logrus.Logger) *sql.DB {
	driver := cfg.Driver
	username := cfg.Username
	password := cfg.Password
	host := cfg.Host
	port := cfg.Port
	dbName := cfg.DbName

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, username, password, dbName)

	conn, err := sql.Open(driver, dsn)
	if err != nil {
		log.Printf("connection database got error: %v", err)
		return nil
	}

	defer func() {
		if err := conn.Close(); err != nil {
			log.Errorf("error closing database connection: %v", err)
		}
	}()

	return conn
}
