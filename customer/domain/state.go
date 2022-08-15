package domain

import (
	"happyday/abstract"

	"github.com/google/uuid"
)

type State struct {
	id      uuid.UUID
	name    string
	comment string
	phones  []Phone
}

func NewState(id uuid.UUID, name string, comment string, phones []Phone) *State {
	return &State{
		id:      id,
		name:    name,
		comment: comment,
		phones:  phones,
	}
}

func NewStateWithID(id uuid.UUID) *State {
	return &State{id: id}
}

func (state *State) ID() uuid.UUID {
	return state.id
}

func (state *State) Name() string {
	return state.name
}

func (state *State) Phones() []Phone {
	return state.phones
}

func (state *State) Comment() string {
	return state.comment
}

func (state *State) On(event abstract.Event) {
	switch e := event.(type) {
	case Created:
		state.name = e.Name()
		state.comment = e.Comment()
		state.phones = e.Phones()
	case NameChanged:
		state.name = e.Name()
	case CommentChanged:
		state.comment = e.Comment()
	case PhonesChanges:
		state.phones = e.Phones()
	}
}
