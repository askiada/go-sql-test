package model

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// DB is the interface that can perform database-related queries.
type DB interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}
