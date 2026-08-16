package cylindrical_cone

import (
	"archimedes-server/core/log"
	"context"
	"math"
)

type cylindricalPartialConeShape struct {
	dimensions map[string]any
}

func NewCylindricalPartialConeShape(dimensions map[string]any) *cylindricalPartialConeShape {
	return &cylindricalPartialConeShape{
		dimensions: dimensions,
	}
}

func (c *cylindricalPartialConeShape) Calculate(ctx context.Context, fluidDistance float64) (float64, error) {
	var biggerRadius, inclineAngle, tankCylindricalHeight, tankConicalHeight float64
	var cylindricalVolume, partialConicalVolume float64

	if err := validate(c.dimensions); err != nil {
		log.Log(err.Error())
		return 0.0, err
	}

	biggerRadius = c.dimensions["bigger_radius"].(float64)
	inclineAngle = c.dimensions["incline_angle"].(float64)
	tankCylindricalHeight = c.dimensions["cylindrical_height"].(float64)
	tankConicalHeight = c.dimensions["conical_height"].(float64)

	cylindricalVolume = c.cylindricalVolume(fluidDistance, biggerRadius, tankCylindricalHeight)
	partialConicalVolume = c.conicalVolume(fluidDistance, biggerRadius, inclineAngle, tankConicalHeight, tankCylindricalHeight)

	return cylindricalVolume + partialConicalVolume, nil
}

func (c *cylindricalPartialConeShape) cylindricalVolume(fluidDistance, tankRadius, tankCylindricalHeight float64) float64 {
	if fluidDistance > tankCylindricalHeight {
		return math.Pi * (tankRadius * tankRadius) * tankCylindricalHeight
	}
	fluidHeight := tankCylindricalHeight - fluidDistance
	return math.Pi * (tankRadius * tankRadius) * fluidHeight
}

func (c *cylindricalPartialConeShape) conicalVolume(fluidDistance, biggerRadius, inclineAngle, tankConicalHeight, tankCylindricalHeight float64) float64 {
	fluidHeight := (tankCylindricalHeight + tankConicalHeight) - fluidDistance

	if fluidHeight > tankCylindricalHeight {

		fluidHeight -= tankCylindricalHeight
		smallerRadius := biggerRadius - (fluidDistance * math.Tan(inclineAngle*math.Pi/180.0))

		return (fluidHeight * math.Pi / 3.0) *
			(biggerRadius*biggerRadius + biggerRadius*smallerRadius + smallerRadius*smallerRadius)
	}
	return 0.0
}
