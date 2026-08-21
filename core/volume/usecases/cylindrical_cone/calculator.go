// Package cylindrical_cone implements IVolumeCalculator for tanks shaped as
// a vertical cylinder with a partial cone at the bottom: the fluid volume
// is the sum of the cylindrical section above the cone plus whatever
// portion of the cone is submerged.
package cylindrical_cone

import (
	"archimedes-server/core/log"
	"context"
	"fmt"
	"math"
)

// tankDimensions are the physical measurements of a cylindrical-cone tank,
// as loaded from its stored dimensions map.
type tankDimensions struct {
	// biggerRadius is the cylinder's (and the cone's widest) radius.
	biggerRadius float64
	// inclineAngle is the cone wall's angle from vertical, in degrees; it
	// determines how quickly the cone's radius shrinks with depth.
	inclineAngle      float64
	cylindricalHeight float64
	conicalHeight     float64
}

// cylindricalPartialConeShape is the domain.CylindricalCone
// IVolumeCalculator implementation. dimensions holds the raw, unvalidated
// shape parameters as read from the database; they are validated and typed
// lazily on each Calculate call.
type cylindricalPartialConeShape struct {
	dimensions map[string]any
}

// NewCylindricalPartialConeShape builds an IVolumeCalculator for a
// cylindrical-cone tank from its raw dimensions map. The map is expected to
// contain bigger_radius, incline_angle, cylindrical_height, and
// conical_height as numbers; validity is checked on Calculate, not here.
func NewCylindricalPartialConeShape(dimensions map[string]any) *cylindricalPartialConeShape {
	return &cylindricalPartialConeShape{
		dimensions: dimensions,
	}
}

// Calculate returns the fluid volume for a sensor reading of fluidDistance,
// the distance from the sensor (mounted at the tank's top) down to the
// fluid surface. A fluidDistance at or beyond the tank's total height
// (cylindrical + conical) yields zero, rather than a negative volume, so
// an over-range sensor reading is reported as an empty tank. It returns an
// error if the tank's dimensions fail schema validation.
func (c *cylindricalPartialConeShape) Calculate(ctx context.Context, fluidDistance float64) (float64, error) {
	var cylindricalVolume, partialConicalVolume float64

	if err := validate(c.dimensions); err != nil {
		log.Log(err.Error())
		return 0.0, err
	}

	biggerRadius, ok := c.dimensions["bigger_radius"].(float64)
	if !ok {
		return 0.0, fmt.Errorf("invalid dimension type")
	}
	inclineAngle, ok := c.dimensions["incline_angle"].(float64)
	if !ok {
		return 0.0, fmt.Errorf("invalid dimension type")
	}
	cylindricalHeight, ok := c.dimensions["cylindrical_height"].(float64)
	if !ok {
		return 0.0, fmt.Errorf("invalid dimension type")
	}
	conicalHeight, ok := c.dimensions["conical_height"].(float64)
	if !ok {
		return 0.0, fmt.Errorf("invalid dimension type")
	}

	dims := tankDimensions{
		biggerRadius:      biggerRadius,
		inclineAngle:      inclineAngle,
		cylindricalHeight: cylindricalHeight,
		conicalHeight:     conicalHeight,
	}

	fluidHeight := (dims.cylindricalHeight + dims.conicalHeight) - fluidDistance
	if fluidHeight < 0 {
		fluidHeight = 0
	}

	cylindricalVolume = c.cylindricalVolume(fluidHeight, dims)

	if fluidHeight > dims.cylindricalHeight {
		partialConicalVolume = c.conicalVolume(fluidHeight, dims)
	}

	return cylindricalVolume + partialConicalVolume, nil
}

// cylindricalVolume returns the volume of fluid held in the cylindrical
// section. If fluidHeight exceeds the cylinder's height, the fluid
// surface is inside the cone below, so the cylinder is reported as full.
func (c *cylindricalPartialConeShape) cylindricalVolume(fluidHeight float64, dims tankDimensions) float64 {
	if fluidHeight > dims.cylindricalHeight {
		return math.Pi * (dims.biggerRadius * dims.biggerRadius) * dims.cylindricalHeight
	}
	return math.Pi * (dims.biggerRadius * dims.biggerRadius) * fluidHeight
}

// conicalVolume returns the volume of fluid held in the submerged part of
// the cone using the frustum volume formula, treating the fluid surface as
// the frustum's smaller face.
func (c *cylindricalPartialConeShape) conicalVolume(fluidHeight float64, dims tankDimensions) float64 {
	fluidHeight -= dims.cylindricalHeight
	smallerRadius := dims.biggerRadius - (fluidHeight * math.Tan(dims.inclineAngle*math.Pi/180.0))

	return (fluidHeight * math.Pi / 3.0) *
		(dims.biggerRadius*dims.biggerRadius + dims.biggerRadius*smallerRadius + smallerRadius*smallerRadius)
}
