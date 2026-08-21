// Package pump implements the IProcessStream use case that turns pump
// status MQTT events into repository writes.
package pump

import (
	"archimedes-server/core/log"
	"archimedes-server/core/processor/domain"
	"archimedes-server/core/pump/interfaces"
	"context"
	"encoding/json"
	"fmt"
)

// ProcessPumpUpdate implements
// archimedes-server/core/processor/interfaces.IProcessStream for the pump
// status MQTT topic.
type ProcessPumpUpdate struct {
	repository interfaces.IUpdatePumpStatus
}

// NewProcessPumpStatusUpdate builds a ProcessPumpUpdate that persists pump
// events through repository.
func NewProcessPumpStatusUpdate(repository interfaces.IUpdatePumpStatus) *ProcessPumpUpdate {
	return &ProcessPumpUpdate{
		repository: repository,
	}
}

// Process unmarshals data as a processor domain.PumpEvent and starts or
// stops the pump's run accordingly. Unrecognized EventType values are
// logged and ignored rather than treated as an error, since the topic may
// carry event types this processor doesn't yet handle.
func (u *ProcessPumpUpdate) Process(ctx context.Context, data []byte) error {
	var event domain.PumpEvent

	err := json.Unmarshal(data, &event)
	if err != nil {
		log.Log(fmt.Sprintf("error on event unmarshal: %q", err.Error()))
		return err
	}

	log.Log(fmt.Sprintf("%s event received: %q", event.EventType, data))

	switch event.EventType {
	case "start":
		err = u.repository.StartPump(ctx, event.PumpID, event.Timestamp)
		if err != nil {
			log.Log(fmt.Sprintf("error on starting pump: %q", err.Error()))
			return err
		}
	case "stop":
		err = u.repository.StopPump(ctx, event.PumpID, event.Timestamp, event.StopReason)
		if err != nil {
			log.Log(fmt.Sprintf("error on stopping pump: %q", err.Error()))
			return err
		}
	default:
		log.Log(fmt.Sprintf("unknown event type: %q", event.EventType))
		return nil
	}

	return nil
}
