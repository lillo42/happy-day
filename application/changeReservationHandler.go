package application

import (
	"context"
	"errors"
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
		PaymentInstallments []reservation.PaymentInstallment `json:"paymentInstallments,omitempty"`
		Comment             string                           `json:"comment,omitempty"`
		Customer            reservation.Customer             `json:"customer"`
		Address             reservation.Address              `json:"address"`
	}

	ChangeReservationHandler struct {
		repository infrastructure.ReservationRepository
	}
)

var (
	ErrReservationPaymentInstallmentAmount = errors.New("payment installment amount cannot be less or equal to zero")
	ErrReservationAddressCityIsEmpty       = errors.New("address city cannot be empty")
	ErrReservationAddressStreetIsEmpty     = errors.New("address street cannot be empty")
	ErrReservationAddressNumberIsInvalid   = errors.New("address number cannot be empty")
	ErrReservationAddressPostalCodeIsEmpty = errors.New("address postal code cannot be empty")
)

func (handler ChangeReservationHandler) Handle(ctx context.Context, req ChangeReservationRequest) (reservation.State, error) {
	for _, item := range req.PaymentInstallments {
		if item.Amount <= 0 {
			return reservation.State{}, ErrReservationPaymentInstallmentAmount
		}
	}

	err := validateCustomer(req.Customer.State)
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
		return ErrReservationAddressCityIsEmpty
	}

	if len(state.Street) == 0 {
		return ErrReservationAddressStreetIsEmpty
	}

	if len(state.Number) == 0 {
		return ErrReservationAddressNumberIsInvalid
	}

	if len(state.PostalCode) == 0 {
		return ErrReservationAddressPostalCodeIsEmpty
	}

	return nil
}
