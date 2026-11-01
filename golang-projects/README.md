# Go Backend Mini-Projects

A collection of small, runnable Go projects demonstrating practical backend engineering fundamentals. Each project is intentionally self-contained and uses the Go standard library so it is easy to clone, inspect, test, and extend.

## Projects

| Project | What it demonstrates | Run |
|---|---|---|
| [Worker Pool](worker-pool/README.md) | Goroutines, channels, bounded concurrency, cancellation, graceful shutdown | `cd worker-pool && go test ./...` |
| [Token Bucket Rate Limiter](rate-limiter/README.md) | Thread safety, mutexes, time-based refill, concurrent tests | `cd rate-limiter && go test ./...` |
| [URL Shortener API](url-shortener/README.md) | REST API design, HTTP handlers, validation, crypto/rand, concurrency-safe storage | `cd url-shortener && go test ./...` |

## Engineering choices

- Prefer the standard library to keep the examples focused on Go fundamentals.
- Include tests for core behavior and important edge cases.
- Use context cancellation for work that should stop when its caller is done.
- Keep examples small enough to understand, but structured so they can be extended.

## Suggested next improvements

- Add persistence to the URL shortener (PostgreSQL), migrations, and repository interfaces.
- Add metrics and structured logging to the worker pool.
- Add configurable per-key limits and distributed storage to the rate limiter.
- Add Docker Compose, CI, benchmarks, and load tests as the projects evolve.

These are learning/portfolio implementations, not production-ready services. Review limits and operational requirements before deploying them.
