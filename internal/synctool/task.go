package synctool

import (
	"context"
)

type ResultTask[T any] interface {
	Run(context.Context) (T, error)
}

type sliceTask[InT, OutT any] struct {
	slice    []InT
	taskFunc func(context.Context, InT) (OutT, error)
}

func (st sliceTask[InT, OutT]) Run(ctx context.Context) ([]OutT, error) {
	res := make([]OutT, len(st.slice))
	for i, item := range st.slice {
		out, err := st.taskFunc(ctx, item)
		if err != nil {
			return nil, err
		}
		res[i] = out
	}
	return res, nil
}

func SliceTask[InT, OutT any](slice []InT, taskFunc func(context.Context, InT) (OutT, error)) ResultTask[[]OutT] {
	return sliceTask[InT, OutT]{
		slice:    slice,
		taskFunc: taskFunc,
	}
}
