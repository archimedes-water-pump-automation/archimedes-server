package tank

import (
	"archimedes-server/core/log"
	"context"
	"errors"
	"math"
)

var (
	errRadiusNotFloat64 = errors.New("radius is not a float64, skipping...")
	errRadiusNotFound   = errors.New("radius not found, skipping...")
	errHeightNotFloat64 = errors.New("height is not a float64, skipping...")
	errHeightNotFound   = errors.New("height not found, skipping...")
)

type cylindricalShape struct{}

func (*cylindricalShape) Calculate(ctx context.Context, tank *TankShape, fluidHeight float64) (float64, error) {
	var tankRadius float64
	var tankHeight float64

	if radiusKey, ok := tank.Dimensions["radius"]; ok {
		tank.Dimensions["radius"] = radiusKey

		if r, ok := radiusKey.(float64); ok {
			tankRadius = r
		} else {
			log.Log(errRadiusNotFloat64.Error())
			return 0, errRadiusNotFloat64
		}

	} else {
		log.Log(errRadiusNotFound.Error())
		return 0, errRadiusNotFound
	}

	if heightKey, ok := tank.Dimensions["height"]; ok {
		tank.Dimensions["height"] = heightKey

		if r, ok := heightKey.(float64); ok {
			tankHeight = r
		} else {
			log.Log(errHeightNotFloat64.Error())
			return 0, errHeightNotFloat64
		}

	} else {
		log.Log(errHeightNotFound.Error())
		return 0, errHeightNotFound
	}

	return math.Pi * (tankRadius * tankRadius) * (tankHeight - fluidHeight), nil

}
