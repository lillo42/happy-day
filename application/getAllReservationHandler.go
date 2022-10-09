package application

import (
	"context"
	"happy_day/infrastructure"
)

type (
	GetAllReservationHandler struct {
		repository infrastructure.ReservationRepository
	}
)

func (handler GetAllReservationHandler) Handler(ctx context.Context, req infrastructure.ReservationFilter) (infrastructure.Page[infrastructure.ReservationViewModel], error) {
	if req.Size <= 0 {
		req.Size = 50
	}

	if req.Page <= 1 {
		req.Page = 1
	}

	return handler.repository.GetAll(ctx, req)
}
