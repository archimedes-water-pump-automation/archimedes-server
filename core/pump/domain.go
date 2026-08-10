package pump

import "time"

type Pump struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type PumpStatus struct {
	ID         string     `db:"id" json:"id"`
	PumpID     string     `db:"pump_id" json:"pump_id"`
	StartedAt  time.Time  `db:"started_at" json:"started_at"`
	StoppedAt  *time.Time `db:"stopped_at" json:"stopped_at,omitempty"`
	StopReason *string    `db:"stop_reason" json:"stop_reason,omitempty"`
}
