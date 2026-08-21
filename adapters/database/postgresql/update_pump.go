package postgresql

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// updatePumpStatusRepository implements
// archimedes-server/core/pump/interfaces.IUpdatePumpStatus against
// PostgreSQL.
type updatePumpStatusRepository struct {
	pool *archimedesPool
}

// NewUpdatePumpStatusRepository builds an IUpdatePumpStatus backed by pool.
func NewUpdatePumpStatusRepository(pool *pgxpool.Pool) *updatePumpStatusRepository {
	return &updatePumpStatusRepository{
		pool: NewPool(pool),
	}
}

// StartPump inserts a new run for pumpID starting at timestamp.
func (repository *updatePumpStatusRepository) StartPump(ctx context.Context, pumpID string, timestamp time.Time) error {
	err := repository.pool.WriteArchimedes(ctx, `
		INSERT INTO
			archimedes.pump_status (id, pump_id, started_at)
		VALUES
			($1, $2, $3)
		RETURNING pump_id;
	`, uuid.NewString(), pumpID, timestamp)
	if err != nil {
		return err
	}

	return nil
}

// StopPump closes pumpID's most recent open run with the given timestamp
// and reason. It updates whichever run is most recent regardless of
// whether it is already stopped, since the query is not scoped to open
// runs only.
func (repository *updatePumpStatusRepository) StopPump(
	ctx context.Context,
	pumpID string,
	timestamp time.Time,
	stopReason string,
) error {
	err := repository.pool.WriteArchimedes(ctx, `
		WITH status AS (
			SELECT
				id
			FROM
				archimedes.pump_status
			WHERE
				pump_id=$1
			ORDER BY started_at DESC
			LIMIT 1
		)
		UPDATE
			archimedes.pump_status
		SET
			stopped_at = $2,
			stop_reason = $3
		FROM
			status
		WHERE
			archimedes.pump_status.id = status.id
		RETURNING pump_id;
	`, pumpID, timestamp, stopReason)
	if err != nil {
		return err
	}

	return nil
}
