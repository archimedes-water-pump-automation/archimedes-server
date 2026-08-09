package processor

import "context"

type IProcessTankStream interface {
	Process(ctx context.Context, data []byte) error
}
