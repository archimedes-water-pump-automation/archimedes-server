package interfaces

import "context"

type IVolumeCalculator interface {
	Calculate(ctx context.Context, fluidDistance float64) (float64, error)
}
