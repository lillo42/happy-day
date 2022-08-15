package domain

import (
	"strconv"
	"time"

	"happyday/abstract"
)

type (
	AggregateRoot interface {
		abstract.AggregateRoot[*State]
		Create(name, comment string, phones []Phone) error
		ChangeName(name string) error
		ChangeComment(comment string) error
		ChangePhones(phones []Phone) error
	}

	aggregateRoot struct {
		abstract.DefaultAggregateRoot[*State]
	}
)

func NewAggregateRoot(state *State, version int64) *aggregateRoot {
	return &aggregateRoot{
		abstract.NewAggregateRoot(state, version),
	}
}

func (root *aggregateRoot) Create(name, comment string, phones []Phone) error {
	if err := checkName(name); err != nil {
		return err
	}

	if err := checkComment(comment); err != nil {
		return err
	}

	if err := checkPhones(phones); err != nil {
		return err
	}

	root.On(Created{
		name:       name,
		comment:    comment,
		phones:     phones,
		occurredOn: time.Now().UTC(),
	})
	return nil
}

func checkName(name string) error {
	if len(name) == 0 {
		return NameIsEmpty
	}

	if len(name) > 255 {
		return NameIsTooLarge
	}

	return nil
}
func (root *aggregateRoot) ChangeName(name string) error {
	if err := checkName(name); err != nil {
		return err
	}

	root.On(NameChanged{
		name:       name,
		occurredOn: time.Now().UTC(),
	})

	return nil
}

func checkComment(comment string) error {
	if len(comment) > 1_000 {
		return CommentTooLarge
	}

	return nil
}
func (root *aggregateRoot) ChangeComment(comment string) error {
	if err := checkComment(comment); err != nil {
		return err
	}

	root.On(CommentChanged{
		comment:    comment,
		occurredOn: time.Now().UTC(),
	})

	return nil
}

func checkPhones(phones []Phone) error {
	if len(phones) == 0 {
		return PhonesIsEmpty
	}

	for _, phone := range phones {
		length := len(phone.Number())
		if length < 9 || length > 11 {
			return PhoneLengthIsInvalid
		}

		if _, err := strconv.Atoi(phone.Number()); err != nil {
			return PhoneNumberIsInvalid
		}
	}

	return nil
}
func (root *aggregateRoot) ChangePhones(phones []Phone) error {
	if err := checkPhones(phones); err != nil {
		return err
	}

	root.On(PhonesChanges{
		phones:     phones,
		occurredOn: time.Now().UTC(),
	})
	return nil
}
