package application

import (
	"context"
	"happy_day/domain/reservation"
	"happy_day/infrastructure"
)

type (
	GetAllReservationsHandler struct {
		repository infrastructure.ReservationRepository
	}
)

func (handler GetAllReservationsHandler) Handle(ctx context.Context, req infrastructure.ReservationFilter) (infrastructure.Page[reservation.State], error) {
	if req.Size <= 0 {
		req.Size = 50
	}

	if req.Page <= 1 {
		req.Page = 1
	}

	return handler.repository.GetAll(ctx, req)
}
