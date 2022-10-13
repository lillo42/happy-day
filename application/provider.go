package application

import (
	"happy_day/infrastructure"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	ProvideQuoteReservationHandler,

	ProvideCreateOrChangeProductHandler,
	ProvideDeleteProductHandler,
	ProvideGetAllProductsHandler,

	ProvideDeleteCustomerHandler,
	ProvideGetAllCustomersHandler,
)

func ProvideQuoteReservationHandler(repository infrastructure.ProductRepository) QuoteReservationHandler {
	return QuoteReservationHandler{
		repository: repository,
	}
}

func ProvideCreateOrChangeProductHandler(repository infrastructure.ProductRepository) CreateOrChangeProductHandler {
	return CreateOrChangeProductHandler{
		repository: repository,
	}
}

func ProvideDeleteProductHandler(repository infrastructure.ProductRepository) DeleteProductHandler {
	return DeleteProductHandler{
		repository: repository,
	}
}

func ProvideGetAllProductsHandler(repository infrastructure.ProductRepository) GetAllProductsHandler {
	return GetAllProductsHandler{
		repository: repository,
	}
}

func ProvideCreateOrChangeCustomerHandler(repository infrastructure.CustomerRepository) CreateOrChangeCustomerHandler {
	return CreateOrChangeCustomerHandler{
		repository: repository,
	}
}

func ProvideDeleteCustomerHandler(repository infrastructure.CustomerRepository) DeleteCustomerHandler {
	return DeleteCustomerHandler{
		repository: repository,
	}
}

func ProvideGetAllCustomersHandler(repository infrastructure.CustomerRepository) GetAllCustomersHandler {
	return GetAllCustomersHandler{
		repository: repository,
	}
}
