package apis

import (
	"happy_day/application"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	ProvideReservationController,
)

func ProvideReservationController(
	quoteHandler application.QuoteReservationHandler,
) ReservationController {
	return ReservationController{
		quoteHandler: quoteHandler,
	}
}
