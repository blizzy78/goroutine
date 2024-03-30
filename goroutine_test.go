package goroutine

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/matryer/is"
)

func TestGoroutines_AwaitTermination(t *testing.T) {
	is := is.New(t)

	goroutines := New()

	start := time.Now()

	goroutines.Go(context.Background(), func(_ context.Context) {
		time.Sleep(10 * time.Millisecond)
	})

	err := goroutines.AwaitTermination(context.Background())
	is.NoErr(err)

	is.True(time.Since(start) >= 10*time.Millisecond)
	is.True(time.Since(start) < 20*time.Millisecond)
}

func TestGoroutines_AwaitTermination_Multiple(t *testing.T) {
	is := is.New(t)

	goroutines := New()

	start := time.Now()

	for range 10 {
		goroutines.Go(context.Background(), func(_ context.Context) {
			time.Sleep(10 * time.Millisecond)
		})
	}

	err := goroutines.AwaitTermination(context.Background())
	is.NoErr(err)

	is.True(time.Since(start) >= 10*time.Millisecond)
	is.True(time.Since(start) < 20*time.Millisecond)
}

func TestGoroutines_AwaitTermination_Timeout(t *testing.T) {
	is := is.New(t)

	goroutines := New()

	start := time.Now()

	goroutines.Go(context.Background(), func(_ context.Context) {
		time.Sleep(100 * time.Millisecond)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := goroutines.AwaitTermination(ctx)
	is.True(errors.Is(err, context.DeadlineExceeded))

	is.True(time.Since(start) >= 10*time.Millisecond)
	is.True(time.Since(start) < 20*time.Millisecond)
}

func TestGoroutines_CancelAll_NoWait(t *testing.T) {
	is := is.New(t)

	goroutines := New()

	start := time.Now()

	goroutines.Go(context.Background(), func(_ context.Context) {
		time.Sleep(100 * time.Millisecond)
	})

	err := goroutines.CancelAll(context.Background(), false)
	is.NoErr(err)

	is.True(time.Since(start) < 10*time.Millisecond)
}

func TestGoroutines_CancelAll_AwaitTermination(t *testing.T) {
	is := is.New(t)

	goroutines := New()

	start := time.Now()

	goroutines.Go(context.Background(), func(_ context.Context) {
		time.Sleep(10 * time.Millisecond)
	})

	err := goroutines.CancelAll(context.Background(), true)
	is.NoErr(err)

	is.True(time.Since(start) >= 10*time.Millisecond)
	is.True(time.Since(start) < 20*time.Millisecond)
}

func TestGoroutines_CancelAll_AwaitTermination_Timeout(t *testing.T) {
	is := is.New(t)

	goroutines := New()

	start := time.Now()

	goroutines.Go(context.Background(), func(_ context.Context) {
		time.Sleep(100 * time.Millisecond)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := goroutines.CancelAll(ctx, true)
	is.True(errors.Is(err, context.DeadlineExceeded))

	is.True(time.Since(start) >= 10*time.Millisecond)
	is.True(time.Since(start) < 20*time.Millisecond)
}
