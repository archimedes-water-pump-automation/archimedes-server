package tank

import "time"

type Tank struct {
	ID   string `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type TankStatus struct {
	Capacity  float64   `db:"capacity" json:"capacity"`
	Volume    float64   `db:"current_volume" json:"current_volume"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Shape string

const (
	Cylindrical Shape = "cylindrical"
)

type TankShape struct {
	Shape      Shape                  `db:"tank_shape" json:"tank_shape"`
	Dimensions map[string]interface{} `db:"dimensions" json:"dimensions"`
}
