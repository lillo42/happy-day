package domain

import (
	"time"
)

type (
	Created struct {
		name       string
		price      float64
		isEnable   bool
		priority   int64
		products   []Product
		occurredOn time.Time
	}

	NameChanged struct {
		name       string
		occurredOn time.Time
	}

	PriceChanged struct {
		price      float64
		occurredOn time.Time
	}

	Disabled struct{ occurredOn time.Time }
	Enabled  struct{ occurredOn time.Time }

	ProductsChanged struct {
		products   []Product
		occurredOn time.Time
	}

	PriorityChanged struct {
		priority   int64
		occurredOn time.Time
	}
)

func (event Created) Name() string          { return event.name }
func (event Created) Price() float64        { return event.price }
func (event Created) IsEnable() bool        { return event.isEnable }
func (event Created) Priority() int64       { return event.priority }
func (event Created) Products() []Product   { return event.products }
func (event Created) OccurredOn() time.Time { return event.occurredOn }

func (event NameChanged) Name() string          { return event.name }
func (event NameChanged) OccurredOn() time.Time { return event.occurredOn }

func (event PriceChanged) Price() float64        { return event.price }
func (event PriceChanged) OccurredOn() time.Time { return event.occurredOn }

func (event Disabled) OccurredOn() time.Time { return event.occurredOn }
func (event Enabled) OccurredOn() time.Time  { return event.occurredOn }

func (event ProductsChanged) Products() []Product   { return event.products }
func (event ProductsChanged) OccurredOn() time.Time { return event.occurredOn }

func (event PriorityChanged) Priority() int64       { return event.priority }
func (event PriorityChanged) OccurredOn() time.Time { return event.occurredOn }
