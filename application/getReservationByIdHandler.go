package application

import (
	"context"
	"happy_day/domain/reservation"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type GetReservationByIdHandler struct {
	repository infrastructure.ReservationRepository
}

func (handler GetReservationByIdHandler) Handler(ctx context.Context, req uuid.UUID) (reservation.State, error) {
	return handler.repository.Get(ctx, req)
}
