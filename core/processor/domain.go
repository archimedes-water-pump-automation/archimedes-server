package processor

import "time"

type WaterTankEvent struct {
	TankID    string    `json:"tank_id"`
	EventType string    `json:"event_type"`
	Volume    float64   `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}

type PumpEvent struct {
	PumpID     string    `json:"pump_id"`
	EventType  string    `json:"event_type"`
	Timestamp  time.Time `json:"timestamp"`
	StopReason string    `json:"stop_reason,omitempty"`
}
