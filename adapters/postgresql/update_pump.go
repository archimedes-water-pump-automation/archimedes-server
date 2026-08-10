package postgresql

import (
	"context"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type updatePumpStatusRepository struct {
	pool *pgxpool.Pool
}

func NewUpdatePumpStatusRepository(pool *pgxpool.Pool) *updatePumpStatusRepository {
	return &updatePumpStatusRepository{
		pool: pool,
	}
}

func (repository *updatePumpStatusRepository) StartPump(ctx context.Context, pumpID string, timestamp time.Time) error {
	connection, err := repository.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()

	rows, err := connection.Query(ctx, `
		INSERT INTO
			archimedes.pump_status (id, pump_id, started_at)
		VALUES
			($1, $2, $3)
		RETURNING pump_id;
	`, uuid.NewString(), pumpID, timestamp)

	if err != nil {
		return err
	}
	defer rows.Close()

	return pgxscan.ScanOne(&pumpID, rows)
}

func (repository *updatePumpStatusRepository) StopPump(ctx context.Context, pumpID string, timestamp time.Time, stopReason string) error {
	connection, err := repository.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()

	rows, err := connection.Query(ctx, `
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
	defer rows.Close()

	return pgxscan.ScanOne(&pumpID, rows)
}
