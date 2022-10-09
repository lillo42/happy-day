package application

import (
	"context"
	"happy_day/domain/reservation"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type (
	ChangeReservationRequest struct {
		Id                  uuid.UUID                        `json:"id"`
		Discount            float64                          `json:"discount"`
		Delivery            reservation.DeliveryOrPickUp     `json:"delivery"`
		PickUp              reservation.DeliveryOrPickUp     `json:"pickUp"`
		PaymentInstallments []reservation.PaymentInstallment `json:"paymentInstallments"`
		Comment             string                           `json:"comment"`
	}

	ChangeReservationHandler struct {
		reservationRepository infrastructure.ReservationRepository
	}
)

func (handler ChangeReservationHandler) Handler(ctx context.Context, req ChangeReservationRequest) (reservation.State, error) {
	state, err := handler.reservationRepository.Get(ctx, req.Id)
	if err != nil {
		return state, err
	}

	state.Discount = req.Discount
	state.Delivery = req.Delivery
	state.PickUp = req.PickUp
	state.PaymentInstallments = req.PaymentInstallments
	state.Comment = req.Comment

	return handler.reservationRepository.Save(ctx, state)
}
