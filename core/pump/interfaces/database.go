package interfaces

import (
	"archimedes-server/core/pump/domain"
	"context"
	"time"
)

type IUpdatePumpStatus interface {
	StartPump(ctx context.Context, pumpID string, timestamp time.Time) error
	StopPump(ctx context.Context, pumpID string, timestamp time.Time, stopReason string) error
}

type IReadPumpStatus interface {
	GetPumps(ctx context.Context) ([]domain.Pump, error)
	GetPumpStatus(ctx context.Context, pumpID string) (*domain.PumpStatus, error)
	GetPumpStatusHistory(ctx context.Context, pumpID string) ([]domain.PumpStatus, error)
}
