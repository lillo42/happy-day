package domain

import (
	"happyday/abstract"

	"github.com/google/uuid"
)

type State struct {
	id       uuid.UUID
	name     string
	price    float64
	priority int64
	products []Product
	isEnable bool
}

func NewState(id uuid.UUID, name string, price float64, isEnable bool, priority int64, products []Product) *State {
	return &State{
		id:       id,
		name:     name,
		price:    price,
		isEnable: isEnable,
		priority: priority,
		products: products,
	}
}

func NewStateWithID(id uuid.UUID) *State {
	return &State{id: id}
}

func (state *State) ID() uuid.UUID       { return state.id }
func (state *State) Name() string        { return state.name }
func (state *State) Price() float64      { return state.price }
func (state *State) IsEnable() bool      { return state.isEnable }
func (state *State) Products() []Product { return state.products }
func (state *State) Priority() int64     { return state.priority }

func (state *State) On(event abstract.Event) {
	switch e := event.(type) {
	case Created:
		state.name = e.Name()
		state.price = e.Price()
		state.priority = e.Priority()
		state.products = e.Products()
	case NameChanged:
		state.name = e.Name()
	case PriceChanged:
		state.price = e.Price()
	case Enabled:
		state.isEnable = true
	case Disabled:
		state.isEnable = false
	case PriorityChanged:
		state.priority = e.Priority()
	case ProductsChanged:
		state.products = e.Products()
	}
}
