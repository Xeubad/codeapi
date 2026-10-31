package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestSumAccountConsumption(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	mock.ExpectQuery("WITH bounds").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "consumed_since", "consumed_today"}).
			AddRow(int64(1), 2.5, 1.25).
			AddRow(int64(2), 0.0, 0.0))

	repo := &upstreamVendorRepository{sql: db}
	out, err := repo.SumAccountConsumption(context.Background(), []int64{1, 2}, time.Now())
	require.NoError(t, err)
	require.Equal(t, 2.5, out[1].SinceSet)
	require.Equal(t, 1.25, out[1].Today)
	require.Zero(t, out[2].SinceSet)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestSumAccountConsumptionEmptyInput(t *testing.T) {
	repo := &upstreamVendorRepository{sql: nil}
	out, err := repo.SumAccountConsumption(context.Background(), nil, time.Now())
	require.NoError(t, err)
	require.Empty(t, out)
}

func TestSumAccountConsumptionNilDB(t *testing.T) {
	var db *sql.DB
	repo := &upstreamVendorRepository{sql: db}
	out, err := repo.SumAccountConsumption(context.Background(), []int64{1}, time.Now())
	require.NoError(t, err)
	require.Empty(t, out)
}
