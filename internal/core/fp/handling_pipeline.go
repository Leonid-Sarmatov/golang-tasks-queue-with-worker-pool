package fp

import (
	"context"
)

type Handler[I, O any] func(ctx context.Context, x I) (O, error)

type Middleware[I, O any] func(h Handler[I, O]) Handler[I, O]

func Chain[I, O any](h Handler[I, O], middles ...Middleware[I, O]) Handler[I, O] {
	for i := len(middles) - 1; i >= 0; i -= 1 {
		h = middles[i](h)
	}
	return h
}
