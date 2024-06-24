package parser

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestRun1(t *testing.T) {
	t.Parallel()

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ctx := context.Background()
	defer mock.Close(ctx)

	mock.ExpectQuery(".*").WillReturnRows(mock.NewRows([]string{"count"}).AddRow(int64(3)))

	pairs, err := Run(ctx, "testdata/1.sql", mock)
	require.NoError(t, err)

	for _, pair := range pairs {
		pair, err := PrepairPair(pair)
		require.NoError(t, err)
		require.Equal(t, pair.Expected, pair.Actual)
	}
}

func TestRun2(t *testing.T) {
	t.Parallel()

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ctx := context.Background()
	defer mock.Close(ctx)

	mrows := mock.NewRows([]string{"1", "2", "3", "4", "5"})

	mrows.AddRow("coucou", true, 5, 3.14, "{'a': 'b'}")
	mrows.AddRow("coucou2", false, 45, 18, "{'m': 'n'}")

	mock.ExpectQuery(".*").WillReturnRows(mrows)

	mrows2 := mock.NewRows([]string{"1", "2", "3", "4", "5"})

	mrows2.AddRow("coucou", true, 5, 3.14, "{'a': 'b'}")
	mrows2.AddRow("coucou2", false, 45, 18, "{'m': 'n'}")

	mock.ExpectQuery(".*").WillReturnRows(mrows2)
	pairs, err := Run(ctx, "testdata/2.sql", mock)
	require.NoError(t, err)

	for _, pair := range pairs {
		pair, err := PrepairPair(pair)
		require.NoError(t, err)
		require.Equal(t, pair.Expected, pair.Actual)
	}
}

func TestRun3(t *testing.T) {
	t.Parallel()

	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	ctx := context.Background()
	defer mock.Close(ctx)

	mrows := mock.NewRows([]string{"1", "2", "3", "4", "5"})

	mrows.AddRow("coucou", true, 5, 3.14, "{'a': 'b'}")
	mrows.AddRow("coucou2", false, 45, 18, "{'m': 'n'}")

	mock.ExpectQuery(".*").WillReturnRows(mrows)

	mrows2 := mock.NewRows([]string{"1", "2", "3", "4", "5"})

	mrows2.AddRow("coucou", true, 5, 3.14, "{'a': 'b'}")
	mrows2.AddRow("coucou2", false, 45, 18, "{'m': 'n'}")

	mock.ExpectQuery(".*").WillReturnRows(mrows2)
	pairs, err := Run(ctx, "testdata/3.sql", mock)
	require.NoError(t, err)

	for _, pair := range pairs {
		pair, err := PrepairPair(pair)
		require.NoError(t, err)
		require.Equal(t, pair.Expected, pair.Actual)
	}
}
