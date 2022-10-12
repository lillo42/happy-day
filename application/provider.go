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
