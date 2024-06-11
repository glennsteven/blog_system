package postgres

import (
	"blog-system/internal/config"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	// Create a mock logger
	logger := logrus.New()

	// Create a mock database config
	cfg := config.Database{
		Driver:   "postgres",
		Username: "postgres",
		Password: "postgres",
		Host:     "localhost",
		Port:     5432,
		DbName:   "testdb",
	}

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	mock.ExpectClose().WillReturnError(nil)

	conn := NewDatabase(cfg, logger)

	assert.NotNil(t, conn)

	// Close the mock connection
	db.Close()
}
