package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/lib/pq"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

// Err are known errors for PostgresSQ.
const (
	ErrUniqueViolation     = Error("unique_violation")
	ErrNullValueNotAllowed = Error("null_value_not_allowed")
	ErrorUndefinedTable    = Error("undefined_table")
	ErrNoRowsFound         = Error("no rows found")
)

// canceledMessage is an error that occurs when deadline exceeded.
const canceledMessage = "pq: canceling statement due to user request"

// parseSQLError converts the error from pq driver using postgresSQL error codes.
// https://www.postgresql.org/docs/9.3/errcodes-appendix.html
func parsePostgresSQLError(err error) error {
	// Parse by value
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrNoRowsFound
	case errors.Is(err, driver.ErrBadConn):
		return context.DeadlineExceeded
	}

	// Parse by type
	var et *pq.Error
	switch {
	case errors.As(err, &et):
		switch et.Code {
		case "23505":
			return ErrUniqueViolation
		case "42P01":
			return ErrorUndefinedTable
		case "22004":
			return ErrNullValueNotAllowed
		}
	}

	// Parse by message
	switch err.Error() {
	case canceledMessage:
		return context.DeadlineExceeded
	}

	return err
}
