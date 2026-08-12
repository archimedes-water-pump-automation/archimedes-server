package tank

import (
	"archimedes-server/core/log"
	"context"
	"math"
)

type cylindricalShape struct{}

func (*cylindricalShape) Calculate(ctx context.Context, tank *TankShape, fluidHeight float64) (float64, error) {
	var radius float64

	if radiusKey, ok := tank.Dimensions["radius"]; ok {
		tank.Dimensions["radius"] = radiusKey

		if r, ok := radiusKey.(float64); ok {
			radius = r
		} else {
			log.Log("radius is not a float64, skipping...")
		}
		return math.Pi * (radius * radius) * fluidHeight, nil
	} else {
		log.Log("radius not found, skipping...")
	}

	return 0, nil
}
