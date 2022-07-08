package transport

import (
	"context"
)

type Server interface {
	Name() string
	ListenAndServe() error
	Close(ctx context.Context) error
}
