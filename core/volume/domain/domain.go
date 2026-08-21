// Package domain defines the tank shapes the volume calculators understand.
package domain

// Shape identifies the geometry a tank uses for volume calculation, as
// stored in the water_tank.tank_shape column.
type Shape string

const (
	// CylindricalCone is a vertical cylinder with a partial cone at the
	// bottom, as calculated by core/volume/usecases/cylindrical_cone.
	CylindricalCone Shape = "cylindrical_cone"
)
