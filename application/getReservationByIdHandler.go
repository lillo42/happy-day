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

func (handler GetReservationByIdHandler) Handle(ctx context.Context, req uuid.UUID) (reservation.State, error) {
	return handler.repository.GetById(ctx, req)
}
