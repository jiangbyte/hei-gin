package contracts

type Event struct {
	Topic string
	Data  any
}

type EventSubscriber func(event Event)

type EventBus interface {
	Publish(topic string, data any)
	Subscribe(topic string, sub EventSubscriber) func()
}
