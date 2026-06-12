package eventbus

import (
	"sync"

	"hei-gin/sdk/shared/contracts"
	"hei-gin/sdk/shared/safe"
)

var DefaultBus contracts.EventBus = newBus()

const maxConcurrentSubscribers = 256

type subscriberEntry struct {
	id uint64
	fn contracts.EventSubscriber
}

type bus struct {
	mu          sync.RWMutex
	subscribers map[string][]subscriberEntry
	closed      chan struct{}
	wg          sync.WaitGroup
	subIDSeq    uint64
	sem         chan struct{}
}

func newBus() *bus {
	return &bus{
		subscribers: make(map[string][]subscriberEntry),
		closed:      make(chan struct{}),
		sem:         make(chan struct{}, maxConcurrentSubscribers),
	}
}

func (b *bus) Publish(topic string, data any) {
	b.mu.RLock()
	subs, ok := b.subscribers[topic]
	if ok {
		subs = append([]subscriberEntry(nil), subs...)
	}
	b.mu.RUnlock()
	if !ok {
		return
	}
	event := contracts.Event{Topic: topic, Data: data}
	for _, entry := range subs {
		entry := entry
		select {
		case <-b.closed:
			return
		case b.sem <- struct{}{}:
		}
		b.wg.Add(1)
		safe.GoSafe(func() {
			defer func() {
				<-b.sem
				b.wg.Done()
			}()
			select {
			case <-b.closed:
				return
			default:
				entry.fn(event)
			}
		})
	}
}

func (b *bus) Subscribe(topic string, sub contracts.EventSubscriber) func() {
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

const (
	TopicUserConnected    = "user:connected"
	TopicUserDisconnected = "user:disconnected"
	TopicMessageSent      = "message:sent"
	TopicMessageRead      = "message:read"
)
