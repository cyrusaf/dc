package source

import "context"

type Local[T any] struct {
	Data T
}

func (l Local[T]) Read(ctx context.Context) (T, error) {
	return l.Data, nil
}
