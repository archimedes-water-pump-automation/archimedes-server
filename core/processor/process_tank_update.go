package processor

import (
	"archimedes-server/core/log"
	"archimedes-server/core/tank"
	"context"
	"encoding/json"
)

type ProcessTankUpdate struct {
	repository tank.IUpdateTank
	calculateVolume tank.ICalculateVolume
}

func NewProcessTankUpdate(repository tank.IUpdateTank, calculateVolume tank.ICalculateVolume) *ProcessTankUpdate {
	return &ProcessTankUpdate{
		repository: repository,
		calculateVolume: calculateVolume,
	}
}

func (u *ProcessTankUpdate) Process(ctx context.Context, data []byte) error {
	var event WaterTankEvent

	err := json.Unmarshal(data, &event)

	if err != nil {
		log.Log("error on event unmarshal: " + err.Error())
		return err
	}

	log.Log(event.EventType + " event received: " + string(data))

	volume, err := u.calculateVolume.CalculateVolume(ctx, event.TankID, event.FluidHeight)
	if err != nil {
		log.Log("error on calculating tank volume: " + err.Error())
		return err
	}

	err = u.repository.UpdateVolume(ctx, event.TankID, volume, event.Timestamp)

	if err != nil {
		log.Log("error on updating tank volume: " + err.Error())
		return err
	}

	return nil
}
