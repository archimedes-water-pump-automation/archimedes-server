package adapters

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tankRepository struct {
	pool *pgxpool.Pool
}

func NewTankRepository(pool *pgxpool.Pool) *tankRepository {
	return &tankRepository{
		pool: pool,
	}
}

func (repository *tankRepository) UpdateVolume(ctx context.Context, tankID int64, newVolume float64, updatedAt time.Time) error {
	connection, err := repository.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()

	rows, err := connection.Query(ctx, `
		UPDATE
			archimedes.water_tank
		SET
			current_volume = $2,
			updated_at = $3
		WHERE
			id = $1
		RETURNING id;
	`, tankID, newVolume, updatedAt)

	if err != nil {
		return err
	}
	defer rows.Close()

	return pgxscan.ScanOne(&tankID, rows)
}
