// Package interfaces declares the tank repository contracts, letting the
// core use cases depend on an abstraction instead of the PostgreSQL adapter.
package interfaces

import (
	"archimedes-server/core/tank/domain"
	"context"
	"time"
)

// IUpdateTank persists tank volume changes. Implemented by
// adapters/database/postgresql.
type IUpdateTank interface {
	// UpdateVolume overwrites the tank's stored volume and updated-at
	// timestamp. It is a no-op if tankID does not exist.
	UpdateVolume(ctx context.Context, tankID string, newVolume float64, updatedAt time.Time) error
}

// IReadTank reads tanks and their current status. Implemented by
// adapters/database/postgresql and consumed by the HTTP read API.
type IReadTank interface {
	// GetTanks lists every registered tank.
	GetTanks(ctx context.Context) ([]domain.Tank, error)
	// GetTankByID returns the latest status for tankID, or nil if the tank
	// does not exist.
	GetTankByID(ctx context.Context, tankID string) (*domain.TankStatus, error)
}
