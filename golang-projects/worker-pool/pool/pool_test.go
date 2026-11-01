package pool

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestPoolExecutesSubmittedJobs(t *testing.T) {
	p := New(context.Background(), 3, 4)
	var completed atomic.Int32
	for i := 0; i < 10; i++ {
		if err := p.Submit(context.Background(), func(context.Context) error {
			completed.Add(1)
			return nil
		}); err != nil {
			t.Fatalf("Submit() error = %v", err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := p.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}
	if got := completed.Load(); got != 10 {
		t.Fatalf("completed = %d, want 10", got)
	}
}

func TestSubmitAfterShutdown(t *testing.T) {
	p := New(context.Background(), 1, 0)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := p.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
	err := p.Submit(context.Background(), func(context.Context) error { return nil })
	if !errors.Is(err, ErrClosed) {
		t.Fatalf("Submit() error = %v, want ErrClosed", err)
	}
}
