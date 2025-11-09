package ports

import (
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, topic string, payload []byte) error
	Close() error
}
