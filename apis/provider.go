package apis

import (
	"happy_day/application"
	"happy_day/infrastructure"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	ProvideReservationController,

	ProvideProductController,
)

func ProvideReservationController(
	quoteHandler application.QuoteReservationHandler,
) ReservationController {
	return ReservationController{
		quoteHandler: quoteHandler,
	}
}

func ProvideProductController(
	createOrChangeHandler application.CreateOrChangeProductHandler,
	deleteHandler application.DeleteProductHandler,
	getAllHandler application.GetAllProductsHandler,
	repository infrastructure.ProductRepository,
) ProductController {
	return ProductController{
		createOrChangeHandler: createOrChangeHandler,
		deleteHandler:         deleteHandler,
		getAllHandler:         getAllHandler,
		repository:            repository,
	}
}
