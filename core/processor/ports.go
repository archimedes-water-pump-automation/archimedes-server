package processor

import "context"

type IProcessStream interface {
	Process(ctx context.Context, data []byte) error
}
