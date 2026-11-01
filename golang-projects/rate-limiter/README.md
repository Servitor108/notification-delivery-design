# Token Bucket Rate Limiter

A small, concurrency-safe in-memory token bucket. It permits a configured burst and replenishes tokens at a steady rate.

## Run

```bash
go test -race ./...
```

## Example

```go
bucket := limiter.New(10, 2) // burst of 10, refill at 2 tokens/second
if !bucket.Allow() {
    // Return HTTP 429 or defer the operation.
}
```

## Concepts

- Token bucket burst capacity and refill rate.
- Mutex-protected state for concurrent goroutines.
- Monotonic elapsed-time calculation via Go's time package.
- Race-detector-friendly tests.

This implementation is process-local. In a multi-instance service, each instance has its own bucket; a distributed limit would require shared coordination (for example Redis or an API gateway).
