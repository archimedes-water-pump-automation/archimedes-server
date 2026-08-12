package tank

import (
	"context"
	"time"
)

type ICalculateVolume interface {
	CalculateVolume(ctx context.Context, tankID string, fluidHeight float64) (float64, error)
}

type IShape interface {
	Calculate(ctx context.Context, tank *TankShape, fluidHeight float64) (float64, error)
}

type IUpdateTank interface {
	UpdateVolume(ctx context.Context, tankID string, newVolume float64, updatedAt time.Time) error
}

type IReadTank interface {
	GetTanks(ctx context.Context) ([]Tank, error)
}

type IReadTankShape interface {
	GetTankShape(ctx context.Context, tankID string) (*TankShape, error)
}
