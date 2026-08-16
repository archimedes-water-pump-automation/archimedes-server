package interfaces

import "context"

type IGetVolumeType interface {
	GetVolumeFromShape(ctx context.Context, tankID string) (IVolumeCalculator, error)
}
