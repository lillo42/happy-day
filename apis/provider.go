package apis

import (
	"happy_day/application"

	"github.com/google/wire"
)

var ProvideSet = wire.NewSet(
	ProvideReservationController,

	ProvideProductController,

	ProvideCustomerController,
)

func ProvideReservationController(
	createHandler application.CreateReservationHandler,
	quoteHandler application.QuoteReservationHandler,
) ReservationController {
	return ReservationController{
		createHandler: createHandler,
		quoteHandler:  quoteHandler,
	}
}

func ProvideProductController(
	getAllHandler application.GetAllProductsHandler,
	getByIdHandler application.GetProductByIdHandler,
	createOrChangeHandler application.CreateOrChangeProductHandler,
	deleteHandler application.DeleteProductHandler,
) ProductController {
	return ProductController{
		getAllHandler:         getAllHandler,
		getByIdHandler:        getByIdHandler,
		createOrChangeHandler: createOrChangeHandler,
		deleteHandler:         deleteHandler,
	}
}

func ProvideCustomerController(
	getAllHandler application.GetAllCustomersHandler,
	getByIdHandler application.GetCustomerByIdHandler,
	createOrChangeHandler application.CreateOrChangeCustomerHandler,
	deleteHandler application.DeleteCustomerHandler,
) CustomerController {
	return CustomerController{
		getAllHandler:         getAllHandler,
		getByIdHandler:        getByIdHandler,
		createOrChangeHandler: createOrChangeHandler,
		deleteHandler:         deleteHandler,
	}
}
