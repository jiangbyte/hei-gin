package eventbus

import (
	"testing"
	"time"

	"hei-gin/sdk/shared/contracts"
)

func TestEventBusPublishSubscribe(t *testing.T) {
	b := newBus()
	ch := make(chan contracts.Event, 1)

	unsub := b.Subscribe("demo", func(event contracts.Event) {
		ch <- event
	})
	defer unsub()

	b.Publish("demo", "hello")

	select {
	case event := <-ch:
		if event.Topic != "demo" {
			t.Fatalf("unexpected topic: %s", event.Topic)
		}
		if event.Data != "hello" {
			t.Fatalf("unexpected data: %v", event.Data)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for event")
	}
}
