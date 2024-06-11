package postgres

import (
	"blog-system/internal/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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

	err = conn.Ping()
	if err != nil {
		log.Fatalf("connection database : %v", err)
	}

	return conn
}
