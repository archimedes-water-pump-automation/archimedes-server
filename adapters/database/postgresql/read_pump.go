package postgresql

import (
	"archimedes-server/core/pump/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// readPumpStatusRepository implements
// archimedes-server/core/pump/interfaces.IReadPumpStatus against
// PostgreSQL.
type readPumpStatusRepository struct {
	pool *archimedesPool
}

// NewReadPumpStatusRepository builds an IReadPumpStatus backed by pool.
func NewReadPumpStatusRepository(pool *pgxpool.Pool) *readPumpStatusRepository {
	return &readPumpStatusRepository{
		pool: NewPool(pool),
	}
}

// GetPumps lists every registered pump.
func (repository *readPumpStatusRepository) GetPumps(ctx context.Context) (pumps []domain.Pump, err error) {
	pumps = make([]domain.Pump, 0)

	err = repository.pool.ReadArchimedes(ctx, &pumps, `
		SELECT
			id,
			name
		FROM
			archimedes.pump;
	`)
	if err != nil {
		return nil, err
	}

	return pumps, nil
}

// GetPumpStatus returns pumpID's most recent run, or nil if the pump has no
// recorded runs.
func (repository *readPumpStatusRepository) GetPumpStatus(
	ctx context.Context,
	pumpID string,
) (*domain.PumpStatus, error) {
	pumpStatus := make([]domain.PumpStatus, 0)

	err := repository.pool.ReadArchimedes(ctx, &pumpStatus, `
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
		return nil, err
	}
	if len(pumpStatus) == 0 {
		return nil, nil
	}

	return &pumpStatus[0], nil
}

// GetPumpStatusHistory returns every recorded run for pumpID, most recent
// first.
func (repository *readPumpStatusRepository) GetPumpStatusHistory(
	ctx context.Context,
	pumpID string,
) ([]domain.PumpStatus, error) {
	pumpStatus := make([]domain.PumpStatus, 0)

	err := repository.pool.ReadArchimedes(ctx, &pumpStatus, `
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
			started_at DESC;
	`, pumpID)
	if err != nil {
		return nil, err
	}

	return pumpStatus, nil

}
