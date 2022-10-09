package infrastructure

import (
	"context"
	"happy_day/domain/reservation"

	"github.com/google/uuid"
)

type (
	ReservationOrderBy string
	ReservationFilter  struct {
		Text    string
		Page    int64
		Size    int64
		OrderBy ReservationOrderBy
	}

	ReservationViewModel struct {
	}

	ReservationRepository interface {
		Get(ctx context.Context, id uuid.UUID) (reservation.State, error)
		Save(ctx context.Context, state reservation.State) (reservation.State, error)
		Delete(ctx context.Context, id uuid.UUID) error

		GetAll(ctx context.Context, filter ReservationFilter) (Page[ReservationViewModel], error)
	}
)
