package synctool

import (
	"context"

	errgroup "golang.org/x/sync/errgroup"
)

type Pool[T any] struct {
	eg  *errgroup.Group
	ctx context.Context
	res chan T
	acc func(T, T) T
}

func WithContext[T any](ctx context.Context, zero T, acc func(T, T) T) (*Pool[T], context.Context) {
	eg, resCtx := errgroup.WithContext(ctx)
	res := make(chan T, 1)
	res <- zero

	return &Pool[T]{
		eg:  eg,
		ctx: ctx,
		res: res,
		acc: acc,
	}, resCtx
}

func New[T any](zero T, acc func(T, T) T) *Pool[T] {
	res, _ := WithContext(context.Background(), zero, acc)
	return res
}

func (p *Pool[T]) Wait() (T, error) {
	err := p.eg.Wait()
	if err != nil {
		var zero T
		return zero, err
	}
	return <-p.res, nil
}

func (p *Pool[T]) SetLimit(n int) {
	p.eg.SetLimit(n)
}

func (p *Pool[T]) createGoroutine(task ResultTask[T]) func() error {
	return func() error {
		res, err := task.Run(p.ctx)
		if err != nil {
			return err
		}

		accRes := <-p.res
		p.res <- p.acc(accRes, res)
		return nil
	}
}

func (p *Pool[T]) TryGo(task ResultTask[T]) bool {
	return p.eg.TryGo(p.createGoroutine(task))
}

func (p *Pool[T]) Go(task ResultTask[T]) {
	p.eg.Go(p.createGoroutine(task))
}
