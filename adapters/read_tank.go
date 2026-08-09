package adapters

import (
	"archimedes-server/core/tank"
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type readTankRepository struct {
	pool *pgxpool.Pool
}

func NewReadTankRepository(pool *pgxpool.Pool) *readTankRepository {
	return &readTankRepository{
		pool: pool,
	}
}

func (repository *readTankRepository) GetTanks(ctx context.Context) (tankNames []tank.TankSummary, err error) {
	connection, err := repository.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer connection.Release()

	rows, err := connection.Query(ctx, `
		SELECT
			name,
			id
		FROM
			archimedes.water_tank;
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = pgxscan.ScanAll(&tankNames, rows)
	if err != nil {
		return nil, err
	}

	return tankNames, nil
}

func (repository *readTankRepository) GetTankByID(ctx context.Context, tankID string) (*tank.Tank, error) {
	connection, err := repository.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer connection.Release()

	var waterTank tank.Tank

	rows, err := connection.Query(ctx, `
		SELECT
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
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	err = pgxscan.ScanOne(&waterTank, rows)
	if err != nil {
		return nil, err
	}

	return &waterTank, nil
}
