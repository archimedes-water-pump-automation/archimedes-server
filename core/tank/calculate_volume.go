package tank

import (
	"archimedes-server/core/log"
	"context"
	"errors"
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
		return c.CylindricalPrototype().Calculate(ctx, tankShape, fluidHeight)
	default:
		log.Log("error on calculating tank volume: unknown tank shape")
		return 0, nil
	}
}

func (c *calculateVolume) CylindricalPrototype() IShape {
	return &cylindricalShape{}
}
