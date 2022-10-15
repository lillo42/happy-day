package application

import (
	"context"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type DeleteReservationHandler struct {
	repository infrastructure.ReservationRepository
}

func (handler DeleteReservationHandler) Handle(ctx context.Context, req uuid.UUID) error {
	return handler.repository.Delete(ctx, req)
}
