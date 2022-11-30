//go:build wireinject
// +build wireinject

package apis

import (
	"happy_day/application/product"
	"happy_day/application/reservation"
	"happy_day/infrastructure"

	"github.com/google/wire"

	"happy_day/application/customer"
)

var (
	ProviderSet = wire.NewSet(
		customer.ProviderSet,
		product.ProviderSet,
		reservation.ProvideSet,
		infrastructure.ProviderSet,
	)
)

// Customer

func initializeGetAllCustomerHandler() customer.GetAllHandler {
	wire.Build(ProviderSet)
	return customer.GetAllHandler{}
}

func initializeGetCustomerByIdHandler() customer.GetByIdHandler {
	wire.Build(ProviderSet)
	return customer.GetByIdHandler{}
}

func initializeChangeOrCreateCustomerHandler() customer.ChangeOrCreateHandler {
	wire.Build(ProviderSet)
	return customer.ChangeOrCreateHandler{}

}

func initializeDeleteCustomerHandler() customer.DeleteHandler {
	wire.Build(ProviderSet)
	return customer.DeleteHandler{}
}

// Product
func initializeGetAllProductHandler() product.GetAllHandler {
	wire.Build(ProviderSet)
	return product.GetAllHandler{}
}

func initalizeGetProductByIdHandler() product.GetByIdHandler {
	wire.Build(ProviderSet)
	return product.GetByIdHandler{}
}

func initializeChangeOrCreateProductHandler() product.ChangeOrCreateHandler {
	wire.Build(ProviderSet)
	return product.ChangeOrCreateHandler{}
}

func initializeDeleteProductHandler() product.DeleteHandler {
	wire.Build(ProviderSet)
	return product.DeleteHandler{}
}

// Reservation
func initializeGetAllReservationHandler() reservation.GetAllHandler {
	wire.Build(ProviderSet)
	return reservation.GetAllHandler{}
}

func initializeGetReservationByIdHandler() reservation.GetByIdHandler {
	wire.Build(ProviderSet)
	return reservation.GetByIdHandler{}
}

func initializeCreateReservationHandler() reservation.CreateHandler {
	wire.Build(ProviderSet)
	return reservation.CreateHandler{}
}

func initializeChangeReservationHandler() reservation.ChangeHandler {
	wire.Build(ProviderSet)
	return reservation.ChangeHandler{}
}

func initializeDeleteReservationHandler() reservation.DeleteHandler {
	wire.Build(ProviderSet)
	return reservation.DeleteHandler{}
}

func initializeQuoteReservationHandler() reservation.QuoteHandler {
	wire.Build(ProviderSet)
	return reservation.QuoteHandler{}
}
