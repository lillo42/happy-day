package apis

import (
	"happy_day/application"
	"happy_day/infrastructure"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	ProvideReservationController,

	ProvideProductController,

	ProvideCustomerController,
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

func ProvideCustomerController(
	createOrChangeHandler application.CreateOrChangeCustomerHandler,
	deleteHandler application.DeleteCustomerHandler,
	getAllHandler application.GetAllCustomersHandler,
	repository infrastructure.CustomerRepository,
) CustomerController {
	return CustomerController{
		createOrChangeHandler: createOrChangeHandler,
		deleteHandler:         deleteHandler,
		getAllHandler:         getAllHandler,
		repository:            repository,
	}
}
