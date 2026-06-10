package eventbus

import (
	"sync"

	"hei-gin/api"
	"hei-gin/sdk/middleware"
)

// DefaultBus is the global event bus instance.
var DefaultBus api.EventBus = newBus()

type subscriberEntry struct {
	id   uint64
	fn   api.EventSubscriber
}

type bus struct {
	mu          sync.RWMutex
	subscribers map[string][]subscriberEntry
	closed      chan struct{}
	wg          sync.WaitGroup
	subIDSeq    uint64
}

func newBus() *bus {
	return &bus{
		subscribers: make(map[string][]subscriberEntry),
		closed:      make(chan struct{}),
	}
}

func (b *bus) Publish(topic string, data any) {
	b.mu.RLock()
	subs, ok := b.subscribers[topic]
	b.mu.RUnlock()
	if !ok {
		return
	}
	event := api.Event{Topic: topic, Data: data}
	for _, entry := range subs {
		entry := entry
		b.wg.Add(1)
		middleware.GoSafe(func() {
			defer b.wg.Done()
			select {
			case <-b.closed:
				return
			default:
				entry.fn(event)
			}
		})
	}
}

func (b *bus) Subscribe(topic string, sub api.EventSubscriber) func() {
	b.mu.Lock()
	b.subIDSeq++
	id := b.subIDSeq
	b.subscribers[topic] = append(b.subscribers[topic], subscriberEntry{id: id, fn: sub})
	b.mu.Unlock()

	return func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		subs := b.subscribers[topic]
		for i := range subs {
			if subs[i].id == id {
				b.subscribers[topic] = append(subs[:i], subs[i+1:]...)
				return
			}
		}
	}
}

// Topics used across the system.
const (
	TopicUserConnected    = "user:connected"
	TopicUserDisconnected = "user:disconnected"
	TopicMessageSent      = "message:sent"
	TopicMessageRead      = "message:read"
)
