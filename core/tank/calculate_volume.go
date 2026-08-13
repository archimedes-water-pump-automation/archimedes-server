package tank

import (
	"archimedes-server/core/log"
	"context"
	"errors"
	"math"
)

type calculateVolume struct {
	repository IReadTankShape
}

func NewCalculateVolume(repository IReadTankShape) *calculateVolume {
	return &calculateVolume{
		repository: repository,
	}
}

func (c *calculateVolume) CalculateVolume(ctx context.Context, tankID string, fluidHeight float64) (float64, error) {
	tankShape, err := c.repository.GetTankShape(ctx, tankID)
	if err != nil {
		return 0, err
	}
	if tankShape == nil {
		err := errors.New("no tank shape found for tankID: " + tankID)
		log.Log(err.Error())
		return 0, err
	}

	switch tankShape.Shape {
	case Cylindrical:
		return NewCylindricalShape().Calculate(ctx, tankShape, fluidHeight)
	default:
		log.Log("error on calculating tank volume: unknown tank shape")
		return 0, nil
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////

type cylindricalShape struct{}

func NewCylindricalShape() *cylindricalShape {
	return &cylindricalShape{}
}

func (c *cylindricalShape) Calculate(ctx context.Context, tank *TankShape, fluidHeight float64) (float64, error) {
	var tankRadius float64
	var tankHeight float64

	if radiusKey, ok := tank.Dimensions["radius"]; ok {
		tank.Dimensions["radius"] = radiusKey

		if r, ok := radiusKey.(float64); ok {
			tankRadius = r
		} else {
			err := errors.New("radius is not a float64, skipping...")
			log.Log(err.Error())
			return 0, err
		}

	} else {
		err := errors.New("radius not found, skipping...")
		log.Log(err.Error())
		return 0, err
	}

	if heightKey, ok := tank.Dimensions["height"]; ok {
		tank.Dimensions["height"] = heightKey

		if r, ok := heightKey.(float64); ok {
			tankHeight = r
		} else {
			err := errors.New("height is not a float64, skipping...")
			log.Log(err.Error())
			return 0, err
		}

	} else {
		err := errors.New("height not found, skipping...")
		log.Log(err.Error())
		return 0, err
	}

	return math.Pi * (tankRadius * tankRadius) * (tankHeight - fluidHeight), nil

}
