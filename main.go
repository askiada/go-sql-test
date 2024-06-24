package main

import (
	"context"
	"testing"
	"time"

	"github.com/arsham/retry"
	"github.com/stretchr/testify/require"

	"github.com/askiada/go-sql-test/internal/db"
	"github.com/askiada/go-sql-test/internal/parser"
)

func TestRunSQL(t *testing.T) { //nolint:paralleltest // This test is not parallel
	env, err := parser.GetEnv()
	require.NoError(t, err)

	ctx := context.Background()

	pool, err := db.NewClient(ctx, &env.DBCredentials, retry.Retry{
		Method:   retry.IncrementalDelay,
		Attempts: 3,
		Delay:    100 * time.Millisecond,
	}, 1)

	require.NoError(t, err)
	pairs, err := parser.Run(ctx, env.SQLFile, pool.DBConnection)
	require.NoError(t, err)

	for _, pair := range pairs {
		t.Run(pair.Name, func(t *testing.T) {
			pair, err := parser.PrepairPair(pair)
			require.NoError(t, err)
			require.Equal(t, pair.Expected, pair.Actual)
		})
	}
}

func main() {
	testing.Main(
		nil,
		[]testing.InternalTest{
			{"TestRunSQL", TestRunSQL},
		},
		nil, nil,
	)
}
