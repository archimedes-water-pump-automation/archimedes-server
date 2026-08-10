package pump

import (
	"context"
	"time"
)

type IUpdatePumpStatus interface {
	StartPump(ctx context.Context, pumpID string, timestamp time.Time) error
	StopPump(ctx context.Context, pumpID string, timestamp time.Time, stopReason string) error
}

type IReadPumpStatus interface {
	GetPumps(ctx context.Context) ([]Pump, error)
	GetPumpStatus(ctx context.Context, pumpID string) (*PumpStatus, error)
}
