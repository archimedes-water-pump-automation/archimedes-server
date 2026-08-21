// Package interfaces declares the contracts for resolving and running a
// tank's volume calculation, letting core use cases depend on an
// abstraction instead of the PostgreSQL adapter or a specific tank shape.
package interfaces

import "context"

// IGetVolumeType resolves the volume calculator for a tank based on its
// stored shape and dimensions. Implemented by adapters/database/postgresql.
type IGetVolumeType interface {
	// GetVolumeFromShape looks up tankID's shape and dimensions and returns
	// the matching IVolumeCalculator. It returns an error if the tank does
	// not exist or its shape has no registered calculator.
	GetVolumeFromShape(ctx context.Context, tankID string) (IVolumeCalculator, error)
}
