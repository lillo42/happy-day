package reservation

import (
	"context"

	"happy_day/application/customer"
	"happy_day/domain/reservation"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type (
	ChangeRequest struct {
		Id                  uuid.UUID                        `json:"id"`
		Discount            float64                          `json:"discount"`
		Delivery            reservation.DeliveryOrPickUp     `json:"delivery"`
		PickUp              reservation.DeliveryOrPickUp     `json:"pickUp"`
		PaymentInstallments []reservation.PaymentInstallment `json:"paymentInstallments,omitempty"`
		Comment             string                           `json:"comment,omitempty"`
		Customer            reservation.Customer             `json:"customer"`
		Address             reservation.Address              `json:"address"`
	}

	ChangeHandler struct {
		repository infrastructure.ReservationRepository
	}
)

func (handler ChangeHandler) Handle(ctx context.Context, req ChangeRequest) (reservation.State, error) {
	for _, item := range req.PaymentInstallments {
		if item.Amount <= 0 {
			return reservation.State{}, infrastructure.ErrReservationPaymentInstallmentAmount
		}
	}

	err := customer.Validate(req.Customer.State)
	if err != nil {
		return reservation.State{}, err
	}

	err = validateAddress(req.Address)
	if err != nil {
		return reservation.State{}, err
	}

	state, err := handler.repository.GetById(ctx, req.Id)
	if err != nil {
		return reservation.State{}, err
	}

	state.Delivery = req.Delivery
	state.PickUp = req.PickUp
	state.PaymentInstallments = req.PaymentInstallments
	state.Comment = req.Comment
	state.Customer = req.Customer
	state.Address = req.Address
	state.Discount = req.Discount
	state.FinalPrice = state.Price - state.Discount

	return handler.repository.Save(ctx, state)
}

func validateAddress(state reservation.Address) error {
	if len(state.City) == 0 {
		return infrastructure.ErrReservationAddressCityIsEmpty
	}

	if len(state.Street) == 0 {
		return infrastructure.ErrReservationAddressStreetIsEmpty
	}

	if len(state.Number) == 0 {
		return infrastructure.ErrReservationAddressNumberIsInvalid
	}

	if len(state.PostalCode) == 0 {
		return infrastructure.ErrReservationAddressPostalCodeIsEmpty
	}

	return nil
}
