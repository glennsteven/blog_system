package postgres

import (
	"blog-system/internal/config"
	"context"
	"database/sql/driver"
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // required import library postgres
	"net/url"
	"time"
)

type DriverType string

const Postgres = DriverType("postgres")

type Conn struct {
	*sqlx.DB
	driver DriverType
	dsn    string
}

// parseToPostgresSQLDSN parses the given config into PostgresSQL DSN string.
func parseToPostgresSQLDSN(cfg *config.Database) (string, error) {
	timezone := "utc"
	if len(cfg.Timezone) != 0 {
		timezone = cfg.Timezone
	}

	sslMode := "disable"
	if cfg.SSLEnabled {
		sslMode = "required"
	}

	q := make(url.Values)
	q.Set("timezone", timezone)
	q.Set("sslmode", sslMode)

	source := url.URL{
		Scheme:   string(Postgres),
		Host:     cfg.Host,
		Path:     cfg.Name,
		User:     url.UserPassword(cfg.User, cfg.Password),
		RawQuery: q.Encode(),
	}

	return source.String(), nil
}

func (c *Conn) Driver() DriverType {
	return c.driver
}

func (c *Conn) DSN() string {
	return c.dsn
}

// ParseSQLError parses driver specific error into known errors.
func (c *Conn) ParseSQLError(err error) error {
	switch c.driver {
	case Postgres:
		return parsePostgresSQLError(err)
	}
	return err
}

// CheckConnection returns an error if connection is not ready.
// Otherwise, return nil.
func CheckConnection(ctx context.Context, conn *Conn) error {
	pingErr := driver.ErrBadConn
	for tries := 0; pingErr != nil; tries++ {
		pingErr = conn.DB.Ping()
		time.Sleep(time.Duration(tries) * 100 * time.Millisecond)
		// Cancel by deadline
		if ctx.Err() != nil {
			break
		}
	}

	// Make sure no context error
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Make one round trip to database to make sure if the database ready to
	// handle query.
	_, err := conn.DB.ExecContext(ctx, "SELECT true;")
	if err != nil {
		return err
	}

	return nil
}

func DatabaseOpen(driver DriverType, cfg config.Database) (*Conn, error) {
	var dsn string
	var err error

	switch driver {
	default:
		return nil, errors.New("unsupported driver")
	case Postgres:
		dsn, err = parseToPostgresSQLDSN(&cfg)
	}

	db, err := sqlx.Open(string(driver), dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConnection)
	db.SetMaxIdleConns(cfg.MaxIdleConnection)

	conn := Conn{DB: db, driver: driver, dsn: dsn}

	return &conn, nil
}
