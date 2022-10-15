package application

import (
	"happy_day/infrastructure"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	ProvideGetAllReservationHandler,
	ProvideGetReservationByIdHandler,
	ProvideCreateReservationHandler,
	ProvideChangeReservationHandler,
	ProvideQuoteReservationHandler,
	ProvideDeleteReservationHandler,

	ProvideGetAllProductsHandler,
	ProvideGetProductByIdHandler,
	ProvideCreateOrChangeProductHandler,
	ProvideDeleteProductHandler,

	ProvideGetAllCustomersHandler,
	ProvideGetCustomerByIdHandler,
	ProvideCreateOrChangeCustomerHandler,
	ProvideDeleteCustomerHandler,
)

// Reservation

func ProvideGetAllReservationHandler(repository infrastructure.ReservationRepository) GetAllReservationsHandler {
	return GetAllReservationsHandler{repository: repository}
}

func ProvideGetReservationByIdHandler(repository infrastructure.ReservationRepository) GetReservationByIdHandler {
	return GetReservationByIdHandler{repository: repository}
}

func ProvideCreateReservationHandler(
	productRepository infrastructure.ProductRepository,
	reservationRepository infrastructure.ReservationRepository) CreateReservationHandler {
	return CreateReservationHandler{
		productRepository:     productRepository,
		reservationRepository: reservationRepository,
	}
}

func ProvideChangeReservationHandler(repository infrastructure.ReservationRepository) ChangeReservationHandler {
	return ChangeReservationHandler{repository: repository}
}

func ProvideDeleteReservationHandler(repository infrastructure.ReservationRepository) DeleteReservationHandler {
	return DeleteReservationHandler{repository: repository}
}

func ProvideQuoteReservationHandler(repository infrastructure.ProductRepository) QuoteReservationHandler {
	return QuoteReservationHandler{
		repository: repository,
	}
}

// Product

func ProvideGetAllProductsHandler(repository infrastructure.ProductRepository) GetAllProductsHandler {
	return GetAllProductsHandler{repository: repository}
}

func ProvideGetProductByIdHandler(repository infrastructure.ProductRepository) GetProductByIdHandler {
	return GetProductByIdHandler{repository: repository}
}

func ProvideCreateOrChangeProductHandler(repository infrastructure.ProductRepository) CreateOrChangeProductHandler {
	return CreateOrChangeProductHandler{repository: repository}
}

func ProvideDeleteProductHandler(repository infrastructure.ProductRepository) DeleteProductHandler {
	return DeleteProductHandler{repository: repository}
}

// Customer

func ProvideGetAllCustomersHandler(repository infrastructure.CustomerRepository) GetAllCustomersHandler {
	return GetAllCustomersHandler{repository: repository}
}

func ProvideGetCustomerByIdHandler(repository infrastructure.CustomerRepository) GetCustomerByIdHandler {
	return GetCustomerByIdHandler{repository: repository}
}

func ProvideCreateOrChangeCustomerHandler(repository infrastructure.CustomerRepository) CreateOrChangeCustomerHandler {
	return CreateOrChangeCustomerHandler{
		repository: repository,
	}
}

func ProvideDeleteCustomerHandler(repository infrastructure.CustomerRepository) DeleteCustomerHandler {
	return DeleteCustomerHandler{repository: repository}
}
