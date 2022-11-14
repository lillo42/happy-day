package reservation

import (
	"context"

	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type DeleteHandler struct {
	repository infrastructure.ReservationRepository
}

func (handler DeleteHandler) Handle(ctx context.Context, req uuid.UUID) error {
	return handler.repository.Delete(ctx, req)
}
