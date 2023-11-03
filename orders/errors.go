package orders

import "errors"

var (
	ErrOrderNotFound         = errors.New("order not found")
	ErrConcurrencyUpdate     = errors.New("concurrency issue")
	ErrAddressIsEmpty        = errors.New("address is empty")
	ErrAddressIsTooLarge     = errors.New("address is too large")
	ErrDeliveryAtIsInvalid   = errors.New("delivery at is invalid")
	ErrTotalPriceIsInvalid   = errors.New("total price is invalid")
	ErrDiscountIsInvalid     = errors.New("discount is invalid")
	ErrFinalPriceIsInvalid   = errors.New("final price is invalid")
	ErrCustomerNotFound      = errors.New("customer not found")
	ErrProductNotFound       = errors.New("product not found")
	ErrPaymentValueIsInvalid = errors.New("payment value is invalid")
	ErrPaymentInfoIsEmpty    = errors.New("payment info is empty")
	ErrPaymentInfoIsTooLarge = errors.New("payment info is too larger")
)
