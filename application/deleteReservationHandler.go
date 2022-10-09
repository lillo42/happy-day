package application

import (
	"context"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type DeleteReservationHandler struct {
	reservationRepository infrastructure.ReservationRepository
}

func (handler DeleteReservationHandler) Handler(ctx context.Context, req uuid.UUID) error {
	return handler.reservationRepository.Delete(ctx, req)
}
