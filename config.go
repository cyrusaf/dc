package dc

import (
	"context"
	"fmt"
	"sync"
)

type Source[T any] interface {
	Read(ctx context.Context) (T, error)
}

type Config[T any] struct {
	Source Source[T]

	mu   sync.RWMutex
	data T
}

func NewConfig[T any](ctx context.Context, source Source[T], defaults T) (*Config[T], error) {
	config := &Config[T]{
		Source: source,
		data:   defaults,
	}

	err := config.Update(ctx)
	if err != nil {
		return config, fmt.Errorf("error fetching initial config values from source, using defaults specified: %w", err)
	}
	return config, nil
}

func (c *Config[T]) Update(ctx context.Context) error {
	data, err := c.Source.Read(ctx)
	if err != nil {
		return fmt.Errorf("reading source: %w", err)
	}

	c.mu.Lock()
	c.data = data
	c.mu.Unlock()

	return nil
}

func (c *Config[T]) Get() T {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data
}
