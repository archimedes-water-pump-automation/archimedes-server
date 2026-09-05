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
	"time"
)

// timeNow is the clock used to stamp events that reached this server
// without a timestamp of their own. A variable so tests can pin it.
var timeNow = time.Now

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

// Process unmarshals data as a processor domain.PumpEvent and opens or
// closes the pump's run according to the state it reports, storing the
// event's reason as the stop reason.
//
// A state of "unknown" is the controller's last will and stores nothing:
// it says the controller became unreachable, not that the pump stopped,
// and closing a run on it would put a fabricated stop time in the
// history. Other unrecognized states and non-pump events are logged and
// ignored for the same reason, since the topic may carry event types
// this processor doesn't yet handle.
func (u *ProcessPumpUpdate) Process(ctx context.Context, data []byte) error {
	var event domain.PumpEvent

	err := json.Unmarshal(data, &event)
	if err != nil {
		log.Log(fmt.Sprintf("error on event unmarshal: %q", err.Error()))
		return err
	}

	log.Log(fmt.Sprintf("%s event received: %q", event.Event, data))

	if event.Event != domain.EventPump {
		log.Log(fmt.Sprintf("unknown event type: %q", event.Event))
		return nil
	}

	timestamp := event.At(timeNow().UTC())

	switch event.State {
	case domain.PumpStateOn:
		err = u.repository.StartPump(ctx, event.Device, timestamp)
		if err != nil {
			log.Log(fmt.Sprintf("error on starting pump: %q", err.Error()))
			return err
		}
	case domain.PumpStateOff:
		err = u.repository.StopPump(ctx, event.Device, timestamp, event.Reason)
		if err != nil {
			log.Log(fmt.Sprintf("error on stopping pump: %q", err.Error()))
			return err
		}
	default:
		log.Log(fmt.Sprintf("unknown pump state: %q (%q)", event.State, event.Reason))
		return nil
	}

	return nil
}
