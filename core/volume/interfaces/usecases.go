package interfaces

import "context"

// IVolumeCalculator computes a tank's fluid volume for a given shape.
// Implementations hold the tank's dimensions and are obtained through
// IGetVolumeType; see core/volume/usecases/cylindrical_cone for the current
// implementation.
type IVolumeCalculator interface {
	// Calculate returns the fluid volume for a sensor reading of
	// fluidDistance, the distance from the sensor down to the fluid
	// surface. Returns an error if the tank's stored dimensions are
	// missing or invalid.
	Calculate(ctx context.Context, fluidDistance float64) (float64, error)
}
