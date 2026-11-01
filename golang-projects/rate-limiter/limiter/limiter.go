// Package limiter implements an in-memory token bucket rate limiter.
package limiter

import (
	"math"
	"sync"
	"time"
)

// TokenBucket allows bursts up to capacity and replenishes tokens at rate
// tokens per second. It is safe for concurrent use.
type TokenBucket struct {
	mu       sync.Mutex
	capacity float64
	tokens   float64
	rate     float64
	last     time.Time
	now      func() time.Time
}

// New creates a token bucket. Capacity and rate must be positive.
func New(capacity int, rate float64) *TokenBucket {
	if capacity < 1 {
		capacity = 1
	}
	if rate <= 0 || math.IsNaN(rate) || math.IsInf(rate, 0) {
		rate = 1
	}
	now := time.Now
	return &TokenBucket{
		capacity: float64(capacity),
		tokens:   float64(capacity),
		rate:     rate,
		last:     now(),
		now:      now,
	}
}

// Allow reports whether one token can be consumed immediately.
func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	current := b.now()
	elapsed := current.Sub(b.last).Seconds()
	if elapsed > 0 {
		b.tokens = math.Min(b.capacity, b.tokens+elapsed*b.rate)
		b.last = current
	}
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}
