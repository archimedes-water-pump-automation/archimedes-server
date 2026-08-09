package processor

import "time"

type Event struct {
	TankID    string    `json:"tank_id"`
	EventType string    `json:"event_type"`
	Volume    float64   `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}
