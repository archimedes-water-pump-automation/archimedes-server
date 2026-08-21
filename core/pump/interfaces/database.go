// Package interfaces declares the pump repository contracts, letting the
// core use cases depend on an abstraction instead of the PostgreSQL adapter.
package interfaces

import (
	"archimedes-server/core/pump/domain"
	"context"
	"time"
)

// IUpdatePumpStatus records pump start/stop events. Implemented by
// adapters/database/postgresql.
type IUpdatePumpStatus interface {
	// StartPump records a new run beginning at timestamp.
	StartPump(ctx context.Context, pumpID string, timestamp time.Time) error
	// StopPump closes the pump's most recent open run with the given
	// timestamp and reason.
	StopPump(ctx context.Context, pumpID string, timestamp time.Time, stopReason string) error
}

// IReadPumpStatus reads pumps and their run history. Implemented by
// adapters/database/postgresql and consumed by the HTTP read API.
type IReadPumpStatus interface {
	// GetPumps lists every registered pump.
	GetPumps(ctx context.Context) ([]domain.Pump, error)
	// GetPumpStatus returns the pump's most recent run, or nil if the pump
	// has no recorded runs.
	GetPumpStatus(ctx context.Context, pumpID string) (*domain.PumpStatus, error)
	// GetPumpStatusHistory returns every recorded run for the pump, most
	// recent first.
	GetPumpStatusHistory(ctx context.Context, pumpID string) ([]domain.PumpStatus, error)
}
