package parser

import (
	"context"
	"testing"
	"time"

	"github.com/arsham/retry"
	"github.com/stretchr/testify/require"

	"github.com/askiada/go-sql-test/internal/db"
)

func TestRunSQL(t *testing.T) { //nolint:paralleltest // This test is not parallel
	env, err := getEnv()
	require.NoError(t, err)

	ctx := context.Background()

	pool, err := db.NewClient(ctx, &env.DBCredentials, retry.Retry{
		Method:   retry.IncrementalDelay,
		Attempts: 3,
		Delay:    100 * time.Millisecond,
	}, 1)

	require.NoError(t, err)
	pairs, err := run(ctx, env.sqlFile, pool.DBConnection)
	require.NoError(t, err)

	for _, pair := range pairs {
		pair, err := prepairPair(pair)
		require.NoError(t, err)
		require.Equal(t, pair.expected, pair.actual)
	}
}
