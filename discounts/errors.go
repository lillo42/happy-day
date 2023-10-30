package discounts

import "errors"

var (
	ErrDiscountNotFound  = errors.New("discount not found")
	ErrConcurrencyUpdate = errors.New("concurrency issue")
	ErrNameIsEmpty       = errors.New("name is missing")
	ErrNameIsTooLarge    = errors.New("name is to large")
	ErrPriceIsInvalid    = errors.New("price is invalid")
	ErrProductsIsMissing = errors.New("products is missing")
	ErrProductNotFound   = errors.New("product not found")
)
