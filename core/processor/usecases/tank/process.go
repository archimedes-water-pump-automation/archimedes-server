// Package tank implements the IProcessStream use case that turns water
// tank level MQTT events into a computed volume and a repository write.
package tank

import (
	"archimedes-server/core/log"
	"archimedes-server/core/processor/domain"
	tankInterfaces "archimedes-server/core/tank/interfaces"
	volumeInterfaces "archimedes-server/core/volume/interfaces"
	"context"
	"encoding/json"
	"fmt"
)

// ProcessTankUpdate implements
// archimedes-server/core/processor/interfaces.IProcessStream for the water
// tank MQTT topic.
type ProcessTankUpdate struct {
	repository    tankInterfaces.IUpdateTank
	getVolumeType volumeInterfaces.IGetVolumeType
}

// NewProcessTankUpdate builds a ProcessTankUpdate that resolves each tank's
// volume calculator through getVolumeType and persists results through
// repository.
func NewProcessTankUpdate(
	repository tankInterfaces.IUpdateTank,
	getVolumeType volumeInterfaces.IGetVolumeType,
) *ProcessTankUpdate {
	return &ProcessTankUpdate{
		repository:    repository,
		getVolumeType: getVolumeType,
	}
}

// Process unmarshals data as a processor domain.WaterTankEvent, converts
// the event's sensor distance into a volume using the tank's registered
// shape, and writes the result back to the tank repository.
func (u *ProcessTankUpdate) Process(ctx context.Context, data []byte) error {
	var event domain.WaterTankEvent

	err := json.Unmarshal(data, &event)
	if err != nil {
		log.Log(fmt.Sprintf("error on event unmarshal: %q", err.Error()))
		return err
	}

	log.Log(fmt.Sprintf("%s event received: %q", event.EventType, data))

	calculator, err := u.getVolumeType.GetVolumeFromShape(ctx, event.TankID)
	if err != nil {
		log.Log(fmt.Sprintf("error on getting volume calculator: %q", err.Error()))
		return err
	}

	volume, err := calculator.Calculate(ctx, event.FluidDistance)
	if err != nil {
		log.Log(fmt.Sprintf("error on calculating tank volume: %q", err.Error()))
		return err
	}

	err = u.repository.UpdateVolume(ctx, event.TankID, volume, event.Timestamp)
	if err != nil {
		log.Log(fmt.Sprintf("error on updating tank volume: %q", err.Error()))
		return err
	}

	return nil
}
