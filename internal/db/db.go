// Package db provides a wrapper around pgxpool.Pool that retries the queries
package db

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/arsham/retry"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	dbretrier "github.com/askiada/go-sql-test/internal/db/retrier"
	"github.com/askiada/go-sql-test/internal/model"
)

// Client is the type that can perform database-related queries.
type Client struct {
	DBConnection *dbretrier.Pool
}

const dbClientTimeout = time.Second * 10

// NewClient returns a new postgres connection. The connection established will be shared throughout this application.
func NewClient(ctxWithTimeout context.Context, credentials *model.DBCredentials, retrier retry.Retry, maxConns int) (*Client, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctxWithTimeout, dbClientTimeout)
	defer cancel()

	addr := psqlBuildQueryString(
		credentials.User,
		credentials.Pass,
		credentials.Name,
		credentials.Host,
		credentials.Port,
		"disable") + " client_encoding=UTF8"

	config, err := pgxpool.ParseConfig(addr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	maxPoolSize := int32(maxConns)

	if maxConns == 0 {
		maxConns, err = getMaxConnections(addr)
		if err != nil {
			return nil, fmt.Errorf("unable to get max connections: %w", err)
		}

		maxPoolSize = int32(90 * maxConns / 100)
	}

	config.MaxConns = maxPoolSize
	config.MinConns = max(1, int32(maxConns)/10)

	pool, err := pgxpool.NewWithConfig(ctxWithTimeout, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create pool: %w", err)
	}

	dbCli := &Client{
		DBConnection: &dbretrier.Pool{DBPool: pool, Retrier: retrier},
	}

	return dbCli, nil
}

// Close the connections.
func (c *Client) Close() {
	c.DBConnection.Close()
}

// getMaxConnections looks at whether the LocalPgMaxConnections variable is
// greater than 0. If so it will use it to set max idle and open connections on
// the DB.
// Otherwise it will use the max_connections variable on the db using the
// provided address.
func getMaxConnections(addr string) (int, error) {
	config, err := pgx.ParseConfig(addr)
	if err != nil {
		return 0, fmt.Errorf("unable to parse config: %w", err)
	}

	c := stdlib.OpenDB(*config)
	defer c.Close() //nolint:errcheck // we don't care about the error here

	var tmp string

	err = c.QueryRow(`SHOW max_connections;`).Scan(&tmp)
	if err != nil {
		return 0, fmt.Errorf("unable to get max_connections: %w", err)
	}

	mc, err := strconv.Atoi(tmp)
	if err != nil {
		return 0, fmt.Errorf("unable to convert max_connections to int: %w", err)
	}

	return mc, nil
}

func psqlBuildQueryString(user, pass, dbname, host string, port int, sslmode string) string {
	parts := []string{}
	if user != "" {
		parts = append(parts, fmt.Sprintf("user=%s", user))
	}

	if pass != "" {
		parts = append(parts, fmt.Sprintf("password=%s", pass))
	}

	if dbname != "" {
		parts = append(parts, fmt.Sprintf("dbname=%s", dbname))
	}

	if host != "" {
		parts = append(parts, fmt.Sprintf("host=%s", host))
	}

	if port != 0 {
		parts = append(parts, fmt.Sprintf("port=%d", port))
	}

	if sslmode != "" {
		parts = append(parts, fmt.Sprintf("sslmode=%s", sslmode))
	}

	return strings.Join(parts, " ")
}
