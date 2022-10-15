package middlewares

import (
	"happy_day/apis"
	"happy_day/application"
	"happy_day/infrastructure"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ErrorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		err := next(ctx)
		if err != nil {
			problem, exist := errors[err]
			if exist {
				return ctx.JSON(problem.Status, problem)
			}

			ctx.Logger().Error("unexpected error during message processing", err)
			return ctx.JSON(unexpectedErr.Status, unexpectedErr)
		}

		return err
	}
}

type ProblemDetails struct {
	Type    string `json:"type"`
	Title   string `json:"title"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

var (
	unexpectedErr = ProblemDetails{
		Type:    "/api/v1/unexpected-error",
		Title:   "APP000",
		Message: "Unexpected error",
		Status:  http.StatusInternalServerError,
	}

	errors = map[error]ProblemDetails{
		apis.ErrInvalidBody: {
			Type:    "/api/v1/invalid-body",
			Title:   "APP000",
			Message: apis.ErrInvalidBody.Error(),
			Status:  http.StatusBadRequest,
		},

		// Products
		infrastructure.ErrProductConcurrencyIssue: {
			Type:    "/api/v1/products/concurrency-issue",
			Title:   "PRD000",
			Message: infrastructure.ErrProductConcurrencyIssue.Error(),
			Status:  http.StatusConflict,
		},
		infrastructure.ErrProductNotFound: {
			Type:    "/api/v1/products/not-found",
			Title:   "PROD001",
			Message: infrastructure.ErrProductNotFound.Error(),
			Status:  http.StatusNotFound,
		},
		application.ErrProductAmountIsInvalid: {
			Type:    "/api/v1/products/amount-is-invalid",
			Title:   "PROD002",
			Message: application.ErrProductAmountIsInvalid.Error(),
			Status:  http.StatusUnprocessableEntity,
		},
		application.ErrProductNameIsEmpty: {
			Type:    "/api/v1/products/name-is-empty",
			Title:   "PROD003",
			Message: application.ErrProductNameIsEmpty.Error(),
			Status:  http.StatusBadRequest,
		},
		application.ErrProductPriceIsLessThanZero: {
			Type:    "/api/v1/products/price-is-less-than-zero",
			Title:   "PROD004",
			Message: application.ErrProductPriceIsLessThanZero.Error(),
			Status:  http.StatusUnprocessableEntity,
		},

		// Customers
		infrastructure.ErrCustomerConcurrencyIssue: {
			Type:    "/api/v1/customers/concurrency-issue",
			Title:   "CUS000",
			Message: infrastructure.ErrCustomerConcurrencyIssue.Error(),
			Status:  http.StatusConflict,
		},
		infrastructure.ErrCustomerNotFound: {
			Type:    "/api/v1/customers/not-found",
			Title:   "CUS001",
			Message: infrastructure.ErrCustomerNotFound.Error(),
			Status:  http.StatusNotFound,
		},
		application.ErrCustomerNameIsEmpty: {
			Type:    "/api/v1/customers/name-is-empty",
			Title:   "CUS002",
			Message: application.ErrCustomerNameIsEmpty.Error(),
			Status:  http.StatusBadRequest,
		},
		application.ErrCustomerPhonesIsEmpty: {
			Type:    "/api/v1/customers/phones-is-empty",
			Title:   "CUS003",
			Message: application.ErrCustomerPhonesIsEmpty.Error(),
			Status:  http.StatusBadRequest,
		},
		application.ErrCustomerPhoneIsInvalid: {
			Type:    "/api/v1/customers/phone-is-invalid",
			Title:   "CUS004",
			Message: application.ErrCustomerPhoneIsInvalid.Error(),
			Status:  http.StatusUnprocessableEntity,
		},

		// Reservations
		application.ErrProductListIsEmpty: {
			Type:    "/api/v1/reservations/product-list-is-empty",
			Title:   "RSV000",
			Message: application.ErrProductListIsEmpty.Error(),
			Status:  http.StatusBadRequest,
		},

		infrastructure.ErrOneProductNotFound: {
			Type:    "/api/v1/reservations/one-product-not-found",
			Title:   "RSV001",
			Message: infrastructure.ErrOneProductNotFound.Error(),
			Status:  http.StatusUnprocessableEntity,
		},
		application.ErrReservationAddressNumberIsInvalid: {
			Type:    "/api/v1/reservations/invalid-address-number",
			Title:   "RSV002",
			Message: application.ErrReservationAddressNumberIsInvalid.Error(),
			Status:  http.StatusUnprocessableEntity,
		},
		application.ErrReservationAddressStreetIsEmpty: {
			Type:    "/api/v1/reservations/address-street-is-empty",
			Title:   "RSV003",
			Message: application.ErrReservationAddressStreetIsEmpty.Error(),
			Status:  http.StatusBadRequest,
		},
		application.ErrReservationAddressPostalCodeIsEmpty: {
			Type:    "/api/v1/reservations/address-postal-code-is-empty",
			Title:   "RSV004",
			Message: application.ErrReservationAddressPostalCodeIsEmpty.Error(),
			Status:  http.StatusBadRequest,
		},
	}
)
