package limiter

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestBucketAllowsInitialBurst(t *testing.T) {
	b := New(3, 1)
	for i := 0; i < 3; i++ {
		if !b.Allow() {
			t.Fatalf("Allow() denied token %d in initial burst", i+1)
		}
	}
	if b.Allow() {
		t.Fatal("Allow() accepted a fourth token without refill")
	}
}

func TestBucketIsSafeForConcurrentUse(t *testing.T) {
	b := New(100, 0.001)
	var allowed atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if b.Allow() {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()
	if got := allowed.Load(); got != 100 {
		t.Fatalf("allowed = %d, want exactly 100", got)
	}
}
