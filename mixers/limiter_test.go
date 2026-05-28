package mixers

import (
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	vrl := NewRateLimiter(10, 5)
	if vrl == nil {
		t.Fatal("expected non-nil VodkaRateLimiter")
	}
	vrl.Stop()
}

func TestRateLimiterStop(t *testing.T) {
	vrl := NewRateLimiter(10, 5)

	done := make(chan struct{})
	go func() {
		vrl.Stop()
		close(done)
	}()

	select {
	case <-done:
		// Stop() returned without blocking — goroutine has exited or will exit
	case <-time.After(2 * time.Second):
		t.Fatal("Stop() did not return in time; goroutine may be leaked")
	}
}

func TestRateLimiterAllowsUnderBurst(t *testing.T) {
	vrl := NewRateLimiter(10, 3)
	defer vrl.Stop()

	l := vrl.getVisitor("127.0.0.1")
	for i := 0; i < 3; i++ {
		if !l.allow() {
			t.Fatalf("expected request %d to be allowed within burst", i+1)
		}
	}
}

func TestRateLimiterBlocksOverBurst(t *testing.T) {
	vrl := NewRateLimiter(1, 2)
	defer vrl.Stop()

	l := vrl.getVisitor("10.0.0.1")
	l.allow()
	l.allow()
	if l.allow() {
		t.Fatal("expected request to be blocked after burst exhausted")
	}
}
