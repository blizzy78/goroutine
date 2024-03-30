package goroutine

import (
	"context"
	"time"
)

func Example() {
	// This function does the actual work.
	// In this example, we're not using the Context, but you really always should.
	worker := func(_ context.Context) {
		time.Sleep(100 * time.Millisecond)
	}

	goroutines := New()

	// Start a new goroutine.
	goroutines.Go(context.Background(), worker)

	// Cancel all goroutines' contexts, and wait for them to finish.
	_ = goroutines.CancelAll(context.Background(), true)
}
