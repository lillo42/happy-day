package customers

import "errors"

var (
	ErrNotFound             = errors.New("customer not found")
	ErrConcurrencyUpdate    = errors.New("concurrency update")
	ErrNameIsEmpty          = errors.New("name is empty")
	ErrNameIsTooLarge       = errors.New("name is too large")
	ErrPhoneNumberIsInvalid = errors.New("phone number is invalid")
)
