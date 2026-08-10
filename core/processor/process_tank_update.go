package processor

import (
	"archimedes-server/core/log"
	"archimedes-server/core/tank"
	"context"
	"encoding/json"
)

type ProcessTankUpdate struct {
	repository tank.IUpdateTank
}

func NewProcessTankUpdate(repository tank.IUpdateTank) *ProcessTankUpdate {
	return &ProcessTankUpdate{
		repository: repository,
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

	err = u.repository.UpdateVolume(ctx, event.TankID, event.Volume, event.Timestamp)

	if err != nil {
		log.Log("error on updating tank volume: " + err.Error())
		return err
	}

	return nil
}
