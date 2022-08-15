package domain

import "time"

type (
	Created struct {
		name       string
		comment    string
		phones     []Phone
		occurredOn time.Time
	}

	NameChanged struct {
		name       string
		occurredOn time.Time
	}

	CommentChanged struct {
		comment    string
		occurredOn time.Time
	}

	PhonesChanges struct {
		phones     []Phone
		occurredOn time.Time
	}
)

func (event CommentChanged) Comment() string       { return event.comment }
func (event CommentChanged) OccurredOn() time.Time { return event.occurredOn }

func (event PhonesChanges) Phones() []Phone       { return event.phones }
func (event PhonesChanges) OccurredOn() time.Time { return event.occurredOn }

func (event NameChanged) Name() string          { return event.name }
func (event NameChanged) OccurredOn() time.Time { return event.occurredOn }

func (event Created) Name() string          { return event.name }
func (event Created) Comment() string       { return event.comment }
func (event Created) Phones() []Phone       { return event.phones }
func (event Created) OccurredOn() time.Time { return event.occurredOn }
