// Package domain holds the core water tank entities, independent of how
// they are stored or transported.
package domain

import "time"

// Tank identifies a physical water tank registered in the system.
type Tank struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// TankStatus is the latest known volume reading for a tank.
type TankStatus struct {
	Capacity  float64   `db:"capacity" json:"capacity"`
	Volume    float64   `db:"current_volume" json:"current_volume"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
