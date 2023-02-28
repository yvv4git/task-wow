package transport

import "context"

type Transport interface {
	Listen(ctx context.Context) error
	Close() error
}
