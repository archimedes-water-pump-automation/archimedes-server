// Package domain holds the core pump entities, independent of how they are
// stored or transported.
package domain

import "time"

// Pump identifies a physical water pump registered in the system.
type Pump struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// PumpStatus is one run of a pump: when it started and, once known, when it
// stopped and why. StoppedAt and StopReason are nil while the pump is still
// running.
type PumpStatus struct {
	ID         string     `db:"id" json:"id"`
	PumpID     string     `db:"pump_id" json:"pump_id"`
	StartedAt  time.Time  `db:"started_at" json:"started_at"`
	StoppedAt  *time.Time `db:"stopped_at" json:"stopped_at,omitempty"`
	StopReason *string    `db:"stop_reason" json:"stop_reason,omitempty"`
}
