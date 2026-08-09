package tank

import "time"

type Tank struct {
	Name      string    `db:"name" json:"name"`
	Capacity  float64   `db:"capacity" json:"capacity"`
	Volume    float64   `db:"current_volume" json:"current_volume"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type TankSummary struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
