package processor

import (
	"archimedes-worker/core/log"
	"archimedes-worker/core/tank"
	"context"
	"encoding/json"
)

type ProcessTankStreamUseCase struct {
	repository tank.ITankRepository
}

func NewProcessTankStreamUseCase(repository tank.ITankRepository) *ProcessTankStreamUseCase {
	return &ProcessTankStreamUseCase{
		repository: repository,
	}
}

func (u *ProcessTankStreamUseCase) Process(ctx context.Context, data []byte) error {
	var event Event

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
