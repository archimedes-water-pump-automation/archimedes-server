package tank

import (
	"context"
	"time"
)

type ITankRepository interface {
	UpdateVolume(ctx context.Context, tankID int64, newVolume float64, updatedAt time.Time) error
}
