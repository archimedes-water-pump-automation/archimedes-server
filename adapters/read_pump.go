package adapters

import (
	"archimedes-server/core/pump"
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type readPumpStatusRepository struct {
	pool *pgxpool.Pool
}

func NewReadPumpStatusRepository(pool *pgxpool.Pool) *readPumpStatusRepository {
	return &readPumpStatusRepository{
		pool: pool,
	}
}

func (repository *readPumpStatusRepository) GetPumps(ctx context.Context) (pumps []pump.Pump, err error) {
	connection, err := repository.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer connection.Release()

	rows, err := connection.Query(ctx, `
		SELECT
			id,
			name
		FROM
			archimedes.pump;
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = pgxscan.ScanAll(&pumps, rows)
	if err != nil {
		return nil, err
	}

	return pumps, nil
}

func (repository *readPumpStatusRepository) GetPumpStatus(ctx context.Context, pumpID string) (*pump.PumpStatus, error) {
	connection, err := repository.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer connection.Release()

	var pumpStatus pump.PumpStatus

	rows, err := connection.Query(ctx, `
		SELECT
			id,
			pump_id,
			started_at,
			stopped_at,
			stop_reason
		FROM
			archimedes.pump_status
		WHERE
			pump_id = $1
		ORDER BY
			started_at DESC
		LIMIT 1;
	`, pumpID)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	err = pgxscan.ScanOne(&pumpStatus, rows)
	if err != nil {
		return nil, err
	}

	return &pumpStatus, nil
}
