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
	"time"
)

// timeNow is the clock used to stamp events that reached this server
// without a timestamp of their own. A variable so tests can pin it.
var timeNow = time.Now

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
// shape, and writes the result back to the tank repository against the
// event's device id.
//
// Readings the tank node flagged invalid store nothing: the distance is
// null, and treating that as zero centimetres would record a tank filled
// to the sensor. Events that are not level events are ignored the same
// way, since the topic may carry event types this processor doesn't yet
// handle.
func (u *ProcessTankUpdate) Process(ctx context.Context, data []byte) error {
	var event domain.WaterTankEvent

	err := json.Unmarshal(data, &event)
	if err != nil {
		log.Log(fmt.Sprintf("error on event unmarshal: %q", err.Error()))
		return err
	}

	log.Log(fmt.Sprintf("%s event received: %q", event.Event, data))

	if event.Event != domain.EventLevel {
		log.Log(fmt.Sprintf("unknown event type: %q", event.Event))
		return nil
	}

	fluidDistance, ok := event.Reading()
	if !ok {
		log.Log(fmt.Sprintf(
			"tank %q reported no usable reading (%q), volume left unchanged",
			event.Device, event.Reason,
		))
		return nil
	}

	calculator, err := u.getVolumeType.GetVolumeFromShape(ctx, event.Device)
	if err != nil {
		log.Log(fmt.Sprintf("error on getting volume calculator: %q", err.Error()))
		return err
	}

	volume, err := calculator.Calculate(ctx, fluidDistance)
	if err != nil {
		log.Log(fmt.Sprintf("error on calculating tank volume: %q", err.Error()))
		return err
	}

	err = u.repository.UpdateVolume(ctx, event.Device, volume, event.At(timeNow().UTC()))
	if err != nil {
		log.Log(fmt.Sprintf("error on updating tank volume: %q", err.Error()))
		return err
	}

	return nil
}
