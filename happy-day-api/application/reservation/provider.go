package reservation

import (
	"happy_day/infrastructure"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	ProvideGetAllHandler,
	ProvideGetByIdHandler,
	ProvideCreateHandler,
	ProvideChangeHandler,
	ProvideQuoteHandler,
	ProvideDeleteHandler,
)

func ProvideGetAllHandler(repository infrastructure.ReservationRepository) GetAllHandler {
	return GetAllHandler{repository: repository}
}

func ProvideGetByIdHandler(repository infrastructure.ReservationRepository) GetByIdHandler {
	return GetByIdHandler{repository: repository}
}

func ProvideCreateHandler(
	productRepository infrastructure.ProductRepository,
	reservationRepository infrastructure.ReservationRepository) CreateHandler {
	return CreateHandler{
		productRepository:     productRepository,
		reservationRepository: reservationRepository,
	}
}

func ProvideChangeHandler(repository infrastructure.ReservationRepository) ChangeHandler {
	return ChangeHandler{repository: repository}
}

func ProvideDeleteHandler(repository infrastructure.ReservationRepository) DeleteHandler {
	return DeleteHandler{repository: repository}
}

func ProvideQuoteHandler(repository infrastructure.ProductRepository) QuoteHandler {
	return QuoteHandler{
		repository: repository,
	}
}
