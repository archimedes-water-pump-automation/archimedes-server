package postgresql

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type updateTankRepository struct {
	pool *pgxpool.Pool
}

func NewUpdateTankRepository(pool *pgxpool.Pool) *updateTankRepository {
	return &updateTankRepository{
		pool: pool,
	}
}

func (repository *updateTankRepository) UpdateVolume(ctx context.Context, tankID string, newVolume float64, updatedAt time.Time) error {
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
