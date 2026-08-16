package interfaces

import (
	"archimedes-server/core/tank/domain"
	"context"
	"time"
)

type IUpdateTank interface {
	UpdateVolume(ctx context.Context, tankID string, newVolume float64, updatedAt time.Time) error
}

type IReadTank interface {
	GetTanks(ctx context.Context) ([]domain.Tank, error)
	GetTankByID(ctx context.Context, tankID string) (*domain.TankStatus, error)
}
