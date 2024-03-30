[![GoDoc](https://pkg.go.dev/badge/github.com/blizzy78/goroutine)](https://pkg.go.dev/github.com/blizzy78/goroutine)


goroutine
=========

A Go package that provides a simple way to manage goroutines and facilitate their graceful shutdown.

```go
import "github.com/blizzy78/goroutine"
```


Code example
------------

```go
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
```


License
-------

This package is licensed under the MIT license.
