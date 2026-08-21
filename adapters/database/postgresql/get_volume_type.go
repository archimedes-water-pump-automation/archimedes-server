package postgresql

import (
	"archimedes-server/core/volume/domain"
	"archimedes-server/core/volume/interfaces"
	"archimedes-server/core/volume/usecases/cylindrical_cone"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// getVolumeType implements
// archimedes-server/core/volume/interfaces.IGetVolumeType against
// PostgreSQL.
type getVolumeType struct {
	pool *archimedesPool
}

// NewGetVolumeType builds an IGetVolumeType backed by pool.
func NewGetVolumeType(pool *pgxpool.Pool) *getVolumeType {
	return &getVolumeType{
		pool: NewPool(pool),
	}
}

// GetVolumeFromShape reads tankID's shape and dimensions and returns the
// matching volume.IVolumeCalculator. It returns an error if the tank does
// not exist or its shape is not one of the shapes registered below.
func (repository *getVolumeType) GetVolumeFromShape(
	ctx context.Context,
	tankID string,
) (interfaces.IVolumeCalculator, error) {
	tankShapes := []struct {
		Shape      domain.Shape   `db:"tank_shape"`
		Dimensions map[string]any `db:"dimensions"`
	}{}

	err := repository.pool.ReadArchimedes(ctx, &tankShapes, `
		SELECT
			tank_shape,
			dimensions
		FROM
			archimedes.water_tank
		WHERE
			id = $1;
	`, tankID)
	if err != nil {
		return nil, err
	}
	if len(tankShapes) == 0 {
		return nil, errors.New("tank not found")
	}

	tankShape := tankShapes[0]

	switch tankShape.Shape {
	case domain.CylindricalCone:
		return cylindrical_cone.NewCylindricalPartialConeShape(tankShape.Dimensions), nil
	default:
		return nil, fmt.Errorf("unsupported tank shape: %q", tankShape.Shape)
	}
}
