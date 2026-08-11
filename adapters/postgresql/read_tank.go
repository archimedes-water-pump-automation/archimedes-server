package postgresql

import (
	"archimedes-server/core/tank"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type readTankRepository struct {
	pool *archimedesPool
}

func NewReadTankRepository(pool *pgxpool.Pool) *readTankRepository {
	return &readTankRepository{
		pool: NewPool(pool),
	}
}

func (repository *readTankRepository) GetTanks(ctx context.Context) (tankNames []tank.TankSummary, err error) {
	err = repository.pool.QueryArchimedes(ctx, &tankNames, `
		SELECT
			name,
			id
		FROM
			archimedes.water_tank;
	`)

	if err != nil {
		return nil, err
	}

	return tankNames, nil
}

func (repository *readTankRepository) GetTankByID(ctx context.Context, tankID string) (*tank.Tank, error) {
	var waterTanks []tank.Tank

	err := repository.pool.QueryArchimedes(ctx, &waterTanks, `
		SELECT
			id,
			name,
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
