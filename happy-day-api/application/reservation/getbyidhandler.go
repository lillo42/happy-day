package reservation

import (
	"context"

	"happy_day/domain/reservation"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type GetByIdHandler struct {
	repository infrastructure.ReservationRepository
}

func (handler GetByIdHandler) Handle(ctx context.Context, req uuid.UUID) (reservation.State, error) {
	return handler.repository.GetById(ctx, req)
}
