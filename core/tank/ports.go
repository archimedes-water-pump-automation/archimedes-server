package tank

import (
	"context"
	"time"
)

type IUpdateTank interface {
	UpdateVolume(ctx context.Context, tankID string, newVolume float64, updatedAt time.Time) error
}

type IReadTank interface {
	GetTanks(ctx context.Context) ([]TankSummary, error)
	GetTankByID(ctx context.Context, tankID string) (*Tank, error)
}
