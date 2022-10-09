package application

import (
	"happy_day/infrastructure"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	ProvideQuoteReservationHandler,
)

func ProvideQuoteReservationHandler(repository infrastructure.ProductRepository) QuoteReservationHandler {
	return QuoteReservationHandler{
		repository: repository,
	}
}
