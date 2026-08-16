package tank

import (
	"archimedes-server/core/log"
	"archimedes-server/core/processor/domain"
	tankInterfaces "archimedes-server/core/tank/interfaces"
	volumeInterfaces "archimedes-server/core/volume/interfaces"
	"context"
	"encoding/json"
)

type ProcessTankUpdate struct {
	repository    tankInterfaces.IUpdateTank
	getVolumeType volumeInterfaces.IGetVolumeType
}

func NewProcessTankUpdate(repository tankInterfaces.IUpdateTank, getVolumeType volumeInterfaces.IGetVolumeType) *ProcessTankUpdate {
	return &ProcessTankUpdate{
		repository:    repository,
		getVolumeType: getVolumeType,
	}
}

func (u *ProcessTankUpdate) Process(ctx context.Context, data []byte) error {
	var event domain.WaterTankEvent

	err := json.Unmarshal(data, &event)

	if err != nil {
		log.Log("error on event unmarshal: " + err.Error())
		return err
	}

	log.Log(event.EventType + " event received: " + string(data))

	calculator, err := u.getVolumeType.GetVolumeFromShape(ctx, event.TankID)
	if err != nil {
		log.Log("error on getting volume calculator: " + err.Error())
		return err
	}

	volume, err := calculator.Calculate(ctx, event.FluidDistance)
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
