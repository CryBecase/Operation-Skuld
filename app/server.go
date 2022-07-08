package app

import "context"

type Server interface {
	Run() error
	Close(ctx context.Context) error
}
