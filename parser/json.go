package parser

import "encoding/json"

type JSON[T any] struct{}

func (j JSON[T]) Parse(b []byte) (T, error) {
	var t T
	err := json.Unmarshal(b, &t)
	return t, err
}
