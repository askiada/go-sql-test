// Package retrier provides a wrapper around pgxpool.Pool that retries the queries
package retrier

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/arsham/retry"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool is a pgxpoooooool that has a retrier.
type Pool struct {
	DBPool  *pgxpool.Pool
	Retrier retry.Retry
}

// Query executes a query and returns the rows.
func (r *Pool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) { //nolint:ireturn // it's a wrapper
	var (
		rows pgx.Rows
		err  error
	)

	err = r.Retrier.Do(
		func() error {
			rows, err = r.DBPool.Query(ctx, sql, args...) //nolint:sqlclosecheck // rows will be closed by the caller
			if err != nil {
				return enrichPgxError(err)
			}

			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("retries couldn't fix this: %w", err)
	}

	return rows, nil
}

// QueryRowScan is a wrapper on top of QueryRow that also integrates the Scan bit.
func (r *Pool) QueryRowScan(ctx context.Context, sql string, queryParams, scanParams []interface{}) error {
	var err error

	err = r.Retrier.Do(
		func() error {
			err = r.DBPool.QueryRow(ctx, sql, queryParams...).Scan(scanParams...)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return retry.StopError{Err: enrichPgxError(err)}
				}

				return enrichPgxError(err)
			}

			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("retries couldn't fix this: %w", err)
	}

	return nil
}

// Exec executes a query.
func (r *Pool) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	var (
		commandTag pgconn.CommandTag
		err        error
	)

	err = r.Retrier.Do(
		func() error {
			commandTag, err = r.DBPool.Exec(ctx, sql, args...)
			if err != nil {
				return enrichPgxError(err)
			}

			return nil
		},
	)
	if err != nil {
		return commandTag, err //nolint:wrapcheck //there is enough info already
	}

	return commandTag, nil
}

// Close terminates the pool.
func (r *Pool) Close() {
	// I'm not going to retry this.
	r.DBPool.Close()
}

// BeginTransaction returns a transaction and an error. It'll retry till it works.
func (r *Pool) BeginTransaction(ctx context.Context) (pgx.Tx, error) { //nolint:ireturn // it's a wrapper
	var (
		pgTx pgx.Tx
		err  error
	)

	err = r.Retrier.Do(
		func() error {
			pgTx, err = r.DBPool.Begin(ctx)

			return enrichPgxError(err)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("retries couldn't fix this: %w", err)
	}

	return pgTx, nil
}

// CopyFrom executes a copy from.
// It is not recommended to retry a copy from.
func (r *Pool) CopyFrom(
	ctx context.Context,
	tableName pgx.Identifier,
	columnNames []string,
	rowSrc pgx.CopyFromSource,
) (int64, error) {
	// Never retry a copy from
	copied, err := r.DBPool.CopyFrom(ctx, tableName, columnNames, rowSrc)
	if err != nil {
		return 0, enrichPgxError(err)
	}

	return copied, nil
}

func enrichPgxError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		msg := strings.Builder{}
		msg.WriteString("\n")

		if pgErr.Code != "" {
			msg.WriteString(fmt.Sprintf("Code: %s\n", pgErr.Code))
		}

		if pgErr.Detail != "" {
			msg.WriteString(fmt.Sprintf("Detail: %s\n", pgErr.Detail))
		}

		if pgErr.Hint != "" {
			msg.WriteString(fmt.Sprintf("Hint: %s\n", pgErr.Hint))
		}

		msg.WriteString(fmt.Sprintf("Position: %d\n", pgErr.Position))

		msg.WriteString(fmt.Sprintf("InternalPosition: %d\n", pgErr.InternalPosition))

		if pgErr.InternalQuery != "" {
			msg.WriteString(fmt.Sprintf("InternalQuery: %s\n", pgErr.InternalQuery))
		}

		if pgErr.Where != "" {
			msg.WriteString(fmt.Sprintf("Where: %s\n", pgErr.Where))
		}

		if pgErr.SchemaName != "" {
			msg.WriteString(fmt.Sprintf("SchemaName: %s\n", pgErr.SchemaName))
		}

		if pgErr.TableName != "" {
			msg.WriteString(fmt.Sprintf("TableName: %s\n", pgErr.TableName))
		}

		if pgErr.ColumnName != "" {
			msg.WriteString(fmt.Sprintf("ColumnName: %s\n", pgErr.ColumnName))
		}

		if pgErr.DataTypeName != "" {
			msg.WriteString(fmt.Sprintf("DataTypeName: %s\n", pgErr.DataTypeName))
		}

		if pgErr.ConstraintName != "" {
			msg.WriteString(fmt.Sprintf("ConstraintName: %s\n", pgErr.ConstraintName))
		}

		if pgErr.File != "" {
			msg.WriteString(fmt.Sprintf("File: %s\n", pgErr.File))
		}

		if pgErr.Line > 0 {
			msg.WriteString(fmt.Sprintf("Line: %d\n", pgErr.Line))
		}

		if pgErr.Routine != "" {
			msg.WriteString(fmt.Sprintf("Routine: %s\n", pgErr.Routine))
		}

		return fmt.Errorf("%s: %w", msg.String(), err)
	}

	return err
}
