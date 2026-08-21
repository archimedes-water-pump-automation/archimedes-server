package postgresql

import (
	"archimedes-server/core/tank/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// readTankRepository implements
// archimedes-server/core/tank/interfaces.IReadTank against PostgreSQL.
type readTankRepository struct {
	pool *archimedesPool
}

// NewReadTankRepository builds an IReadTank backed by pool.
func NewReadTankRepository(pool *pgxpool.Pool) *readTankRepository {
	return &readTankRepository{
		pool: NewPool(pool),
	}
}

// GetTanks lists every registered tank.
func (repository *readTankRepository) GetTanks(ctx context.Context) ([]domain.Tank, error) {
	waterTanks := make([]domain.Tank, 0)

	err := repository.pool.ReadArchimedes(ctx, &waterTanks, `
		SELECT
			id,
			name
		FROM
			archimedes.water_tank;
	`)
	if err != nil {
		return nil, err
	}

	return waterTanks, nil
}

// GetTankByID returns tankID's latest status, or nil if the tank does not
// exist.
func (repository *readTankRepository) GetTankByID(ctx context.Context, tankID string) (*domain.TankStatus, error) {
	waterTanks := make([]domain.TankStatus, 0)

	err := repository.pool.ReadArchimedes(ctx, &waterTanks, `
		SELECT
			capacity,
			current_volume,
			created_at,
			updated_at
		FROM
			archimedes.water_tank
		WHERE
			id = $1;
	`, tankID)
	if err != nil {
		return nil, err
	}
	if len(waterTanks) == 0 {
		return nil, nil
	}

	return &waterTanks[0], nil
}
