package eventbus

import (
	"sync/atomic"
	"testing"
	"time"

	"hei-gin/sdk/contracts"
)

func TestPublishUsesSubscriberSnapshot(t *testing.T) {
	b := newBus()
	var calls atomic.Int32

	unsubscribe := b.Subscribe("topic", func(event contracts.Event) {
		calls.Add(1)
	})
	unsubscribe()

	b.Publish("topic", "payload")
	b.wg.Wait()

	if calls.Load() != 0 {
		t.Fatalf("unsubscribed handler was called %d times", calls.Load())
	}
}

func TestPublishLimitsConcurrentSubscribers(t *testing.T) {
	b := newBus()
	const subscribers = maxConcurrentSubscribers + 32

	started := make(chan struct{}, subscribers)
	release := make(chan struct{})
	var running atomic.Int32
	var peak atomic.Int32

	for i := 0; i < subscribers; i++ {
		b.Subscribe("topic", func(event contracts.Event) {
			current := running.Add(1)
			for {
				old := peak.Load()
				if current <= old || peak.CompareAndSwap(old, current) {
					break
				}
			}
			started <- struct{}{}
			<-release
			running.Add(-1)
		})
	}

	done := make(chan struct{})
	go func() {
		b.Publish("topic", nil)
		close(done)
	}()

	for i := 0; i < maxConcurrentSubscribers; i++ {
		select {
		case <-started:
		case <-time.After(time.Second):
			close(release)
			t.Fatalf("timed out waiting for subscriber %d", i)
		}
	}

	select {
	case <-done:
		close(release)
		t.Fatal("Publish returned before backpressure was released")
	case <-time.After(20 * time.Millisecond):
	}

	close(release)
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Publish did not finish after subscribers were released")
	}
	b.wg.Wait()

	if got := peak.Load(); got > maxConcurrentSubscribers {
		t.Fatalf("peak concurrency = %d, want <= %d", got, maxConcurrentSubscribers)
	}
}
