package common

import (
	"context"
)

type ShutdownHandler interface {
	Shutdown(ctx context.Context)
}
