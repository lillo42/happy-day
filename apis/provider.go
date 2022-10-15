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
	getAllHandler application.GetAllReservationsHandler,
	getByIdHandler application.GetReservationByIdHandler,
	createHandler application.CreateReservationHandler,
	changeHandler application.ChangeReservationHandler,
	deleteHandler application.DeleteReservationHandler,
	quoteHandler application.QuoteReservationHandler,
) ReservationController {
	return ReservationController{
		getAllHandler:  getAllHandler,
		getByIdHandler: getByIdHandler,
		createHandler:  createHandler,
		changeHandler:  changeHandler,
		deleteHandler:  deleteHandler,
		quoteHandler:   quoteHandler,
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
