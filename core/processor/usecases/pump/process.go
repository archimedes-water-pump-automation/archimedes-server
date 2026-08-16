package pump

import (
	"archimedes-server/core/log"
	"archimedes-server/core/processor/domain"
	"archimedes-server/core/pump/interfaces"
	"context"
	"encoding/json"
)

type ProcessPumpUpdate struct {
	repository interfaces.IUpdatePumpStatus
}

func NewProcessPumpStatusUpdate(repository interfaces.IUpdatePumpStatus) *ProcessPumpUpdate {
	return &ProcessPumpUpdate{
		repository: repository,
	}
}

func (u *ProcessPumpUpdate) Process(ctx context.Context, data []byte) error {
	var event domain.PumpEvent

	err := json.Unmarshal(data, &event)

	if err != nil {
		log.Log("error on event unmarshal: " + err.Error())
		return err
	}

	log.Log(event.EventType + " event received: " + string(data))

	switch event.EventType {
	case "start":
		err = u.repository.StartPump(ctx, event.PumpID, event.Timestamp)

		if err != nil {
			log.Log("error on starting pump: " + err.Error())
			return err
		}
	case "stop":
		err = u.repository.StopPump(ctx, event.PumpID, event.Timestamp, event.StopReason)

		if err != nil {
			log.Log("error on stopping pump: " + err.Error())
			return err
		}
	default:
		log.Log("unknown event type: " + event.EventType)
		return nil
	}

	return nil
}
