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

func (repository *readTankRepository) GetTanks(ctx context.Context) ([]tank.Tank, error) {
	var waterTanks []tank.Tank

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
	if len(waterTanks) == 0 {
		return nil, nil
	}

	return waterTanks, nil
}

func (repository *readTankRepository) GetTankByID(ctx context.Context, tankID string) (*tank.TankStatus, error) {
	var waterTanks []tank.TankStatus

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

func (repository *readTankRepository) GetTankShape(ctx context.Context, tankID string) (*tank.TankShape, error) {
	var waterTanks []tank.TankShape

	err := repository.pool.ReadArchimedes(ctx, &waterTanks, `
		SELECT
			tank_shape,
			dimensions
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
