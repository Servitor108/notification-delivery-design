# Context-Aware Worker Pool

A bounded worker pool built with goroutines and channels. It limits concurrent work, applies backpressure when its queue is full, and supports graceful shutdown.

## Run

```bash
go test ./...
```

## Example

```go
p := pool.New(context.Background(), 4, 20)
err := p.Submit(context.Background(), func(ctx context.Context) error {
    // Perform work; respect ctx cancellation.
    return nil
})
if err != nil {
    log.Fatal(err)
}

shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
if err := p.Shutdown(shutdownCtx); err != nil {
    log.Printf("shutdown: %v", err)
}
```

## Concepts

- Fixed number of goroutines bounds concurrent execution.
- A buffered channel provides a bounded queue and backpressure.
- Submit supports caller cancellation.
- Shutdown stops submissions, drains accepted jobs, and can cancel work on timeout.

The pool intentionally leaves job error reporting to the caller's job wrapper; a production implementation would usually add result collection, panic recovery, metrics, and explicit error handling.
