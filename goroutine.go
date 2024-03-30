package goroutine

import (
	"context"
	"slices"
	"sync"
)

// Goroutines manages a collection of goroutines.
type Goroutines struct {
	sync.Mutex
	goroutines []*goroutine
}

type goroutine struct {
	index  int
	cancel context.CancelFunc
	done   <-chan struct{}
}

// New creates a new Goroutines.
func New() *Goroutines {
	return &Goroutines{}
}

// Go starts fun in a new goroutine. This works basically the same as the language's go keyword.
func (gs *Goroutines) Go(ctx context.Context, fun func(ctx context.Context)) {
	ctx, cancel := context.WithCancel(ctx)
	done := make(chan struct{})

	gor := goroutine{
		cancel: cancel,
		done:   done,
	}

	gs.Lock()

	gs.goroutines = append(gs.goroutines, &gor)
	gor.index = len(gs.goroutines) - 1

	gs.Unlock()

	go func() {
		defer close(done)
		defer gs.remove(&gor)
		defer cancel()

		fun(ctx)
	}()
}

func (gs *Goroutines) remove(gor *goroutine) {
	gs.Lock()
	defer gs.Unlock()

	gs.goroutines = slices.Delete(gs.goroutines, gor.index, gor.index+1)

	for _, other := range gs.goroutines[gor.index:] {
		other.index--
	}
}

// CancelAll cancels all goroutines' contexts. If awaitTermination==true, waits for all goroutines to finish.
// Returns ctx.Err() if ctx is canceled while waiting, or nil otherwise.
func (gs *Goroutines) CancelAll(ctx context.Context, awaitTermination bool) error {
	gors := gs.getGoroutines()

	for _, gor := range gors {
		gor.cancel()
	}

	if !awaitTermination {
		return nil
	}

	return gs.awaitTermination(ctx, gors)
}

// AwaitTermination waits for all goroutines to finish.
// Returns ctx.Err() if ctx is canceled while waiting, or nil otherwise.
func (gs *Goroutines) AwaitTermination(ctx context.Context) error {
	return gs.awaitTermination(ctx, gs.getGoroutines())
}

func (gs *Goroutines) awaitTermination(ctx context.Context, gors []*goroutine) error {
	for _, gor := range gors {
		select {
		case <-gor.done:

		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func (gs *Goroutines) getGoroutines() []*goroutine {
	gs.Lock()

	gors := make([]*goroutine, len(gs.goroutines))
	copy(gors, gs.goroutines)

	gs.Unlock()

	return gors
}
