package abstract

import "time"

type (
	Event interface{ OccurredOn() time.Time }

	State interface{ On(event Event) }

	AggregateRoot[T State] interface {
		State() T
		Version() int64
		Events() []Event
	}

	DefaultAggregateRoot[T State] struct {
		state   T
		version int64
		events  []Event
	}
)

func (root *DefaultAggregateRoot[T]) State() T {
	return root.state
}

func (root *DefaultAggregateRoot[T]) Version() int64 {
	return root.version
}

func (root *DefaultAggregateRoot[T]) Events() []Event {
	return root.events
}

func (root *DefaultAggregateRoot[T]) On(event Event) {
	root.state.On(event)
	root.events = append(root.events, event)
}

func NewAggregateRoot[T State](state T, version int64) DefaultAggregateRoot[T] {
	return DefaultAggregateRoot[T]{
		state:   state,
		version: version,
		events:  make([]Event, 0),
	}
}
