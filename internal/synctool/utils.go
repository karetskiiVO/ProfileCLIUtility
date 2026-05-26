package synctool

import (
	"context"
	"iter"
)

func MapPoolSlice[InT, OutT any](ctx context.Context, workers int, slice []InT, fn func(InT) (OutT, error)) ([]OutT, error) {
	if workers <= 0 {
		workers = 1
	}
	if slice == nil {
		return nil, nil
	}

	if workers == 1 {
		res := make([]OutT, len(slice))
		for i, item := range slice {
			out, err := fn(item)
			if err != nil {
				return nil, err
			}
			res[i] = out
		}
		return res, nil
	}

	pool := New(nil, func(a, b []OutT) []OutT {
		return append(a, b...)
	})
	pool.SetLimit(workers)

	for chunk := range chunkSlice(slice, workers) {
		pool.Go(SliceTask(chunk, addContextToFunc(fn)))
	}

	return pool.Wait()
}

func addContextToFunc[InT, OutT any](f func(InT) (OutT, error)) func(context.Context, InT) (OutT, error) {
	return func(ctx context.Context, item InT) (OutT, error) {
		select {
		case <-ctx.Done():
			var zero OutT
			return zero, ctx.Err()
		default:
			return f(item)
		}
	}
}

func chunkSlice[T any](slice []T, chunkNum int) iter.Seq[[]T] {
	chunkSize := (len(slice) + chunkNum - 1) / chunkNum
	return func(yield func([]T) bool) {
		if chunkSize <= 0 {
			return
		}
		for i := 0; i < len(slice); i += chunkSize {
			end := i + chunkSize
			if end > len(slice) {
				end = len(slice)
			}
			if !yield(slice[i:end]) {
				return
			}
		}
	}
}
