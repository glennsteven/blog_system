package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/lib/pq"
	"io"
	"testing"
)

func TestParseSQLError(t *testing.T) {
	tests := []struct {
		in  error
		out error
	}{
		{
			in:  sql.ErrNoRows,
			out: ErrNoRowsFound,
		},
		{
			in:  driver.ErrBadConn,
			out: context.DeadlineExceeded,
		},
		{
			in:  &pq.Error{Code: "23505"},
			out: ErrUniqueViolation,
		},
		{
			in:  &pq.Error{Code: "42P01"},
			out: ErrorUndefinedTable,
		},
		{
			in:  &pq.Error{Code: "22004"},
			out: ErrNullValueNotAllowed,
		},
		{
			in:  errors.New(canceledMessage),
			out: context.DeadlineExceeded,
		},
		{
			in:  io.EOF,
			out: io.EOF,
		},
	}

	for _, test := range tests {
		if err := parsePostgresSQLError(test.in); err != test.out {
			t.Errorf("expecting %v but got %v", test.out, err)
		}
	}
}
