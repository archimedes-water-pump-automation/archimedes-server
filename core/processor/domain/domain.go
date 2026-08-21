// Package domain defines the payloads carried over the MQTT streams that
// the processor use cases unmarshal and act on.
package domain

import "time"

// WaterTankEvent is a tank-level sensor reading published to the water tank
// MQTT topic.
type WaterTankEvent struct {
	TankID    string `json:"tank_id"`
	EventType string `json:"event_type"`
	// FluidDistance is the sensor-reported distance down to the fluid
	// surface, consumed by IVolumeCalculator.Calculate.
	FluidDistance float64   `json:"distance"`
	Timestamp     time.Time `json:"timestamp"`
}

// PumpEvent is a pump start/stop notification published to the pump status
// MQTT topic. EventType is expected to be "start" or "stop"; StopReason is
// only meaningful for "stop" events.
type PumpEvent struct {
	PumpID     string    `json:"pump_id"`
	EventType  string    `json:"event_type"`
	Timestamp  time.Time `json:"timestamp"`
	StopReason string    `json:"stop_reason,omitempty"`
}
