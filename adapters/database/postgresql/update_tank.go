package postgresql

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// updateTankRepository implements
// archimedes-server/core/tank/interfaces.IUpdateTank against PostgreSQL.
type updateTankRepository struct {
	pool *archimedesPool
}

// NewUpdateTankRepository builds an IUpdateTank backed by pool.
func NewUpdateTankRepository(pool *pgxpool.Pool) *updateTankRepository {
	return &updateTankRepository{
		pool: NewPool(pool),
	}
}

// UpdateVolume overwrites tankID's stored volume and updated-at timestamp.
// It is a no-op if tankID does not exist.
func (repository *updateTankRepository) UpdateVolume(
	ctx context.Context,
	tankID string,
	newVolume float64,
	updatedAt time.Time,
) error {
	err := repository.pool.WriteArchimedes(ctx, `
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

	return nil
}
