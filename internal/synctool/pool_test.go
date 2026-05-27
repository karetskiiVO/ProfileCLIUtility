package synctool

import (
	"context"
	"errors"
	"testing"
	"time"

	require "github.com/stretchr/testify/require"
)

type funcTask[T any] struct {
	fn func(context.Context) (T, error)
}

func (ft funcTask[T]) Run(ctx context.Context) (T, error) {
	return ft.fn(ctx)
}

func recvWithin[T any](t *testing.T, ch <-chan T, d time.Duration) T {
	t.Helper()
	select {
	case v := <-ch:
		return v
	case <-time.After(d):
		require.FailNow(t, "timeout waiting for value")
		var zero T
		return zero
	}
}

func TestPoolWait_NoTasks_ReturnsInitialValue(t *testing.T) {
	pool := New(42, func(a, b int) int { return a + b })

	got, err := pool.Wait()
	require.NoError(t, err)
	require.Equal(t, 42, got)
}

func TestPoolWait_AccumulatesResults(t *testing.T) {
	pool := New(10, func(a, b int) int { return a + b })
	pool.SetLimit(4)

	pool.Go(funcTask[int]{fn: func(context.Context) (int, error) { return 1, nil }})
	pool.Go(funcTask[int]{fn: func(context.Context) (int, error) { return 2, nil }})
	pool.Go(funcTask[int]{fn: func(context.Context) (int, error) { return 3, nil }})

	got, err := pool.Wait()
	require.NoError(t, err)
	require.Equal(t, 16, got)
}

func TestPoolTryGo_RespectsLimit(t *testing.T) {
	pool := New(0, func(a, b int) int { return a + b })
	pool.SetLimit(1)

	started := make(chan struct{})
	release := make(chan struct{})

	ok := pool.TryGo(funcTask[int]{fn: func(context.Context) (int, error) {
		close(started)
		<-release
		return 1, nil
	}})
	require.True(t, ok)

	_ = recvWithin(t, started, 500*time.Millisecond)

	ok = pool.TryGo(funcTask[int]{fn: func(context.Context) (int, error) { return 2, nil }})
	require.False(t, ok)

	close(release)
	got, err := pool.Wait()
	require.NoError(t, err)
	require.Equal(t, 1, got)
}

func TestPoolWait_ErrorReturnsZeroValue(t *testing.T) {
	type result struct {
		Sum int
	}

	boom := errors.New("boom")
	pool := New(result{Sum: 7}, func(a, b result) result { return result{Sum: a.Sum + b.Sum} })

	pool.Go(funcTask[result]{fn: func(context.Context) (result, error) { return result{Sum: 5}, nil }})
	pool.Go(funcTask[result]{fn: func(context.Context) (result, error) { return result{}, boom }})

	got, err := pool.Wait()
	require.Error(t, err)
	require.ErrorIs(t, err, boom)
	require.Equal(t, result{}, got)
}

func TestPoolWithContext_CancelsTaskContextOnError(t *testing.T) {
	parent, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pool, ctx := WithContext(parent, 0, func(a, b int) int { return a + b })
	pool.SetLimit(2)

	boom := errors.New("boom")
	canceled := make(chan error, 1)

	pool.Go(funcTask[int]{fn: func(taskCtx context.Context) (int, error) {
		<-taskCtx.Done()
		err := taskCtx.Err()
		canceled <- err
		return 0, err
	}})

	pool.Go(funcTask[int]{fn: func(context.Context) (int, error) {
		return 0, boom
	}})

	_, err := pool.Wait()
	require.Error(t, err)
	require.ErrorIs(t, err, boom)

	cancelErr := recvWithin(t, canceled, 500*time.Millisecond)
	require.ErrorIs(t, cancelErr, context.Canceled)
	require.ErrorIs(t, ctx.Err(), context.Canceled)
}
