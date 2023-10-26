package products

import "errors"

var (
	ErrProductNotExists    = errors.New("product not exist")
	ErrNameIsEmpty         = errors.New("name is empty")
	ErrNameTooLarge        = errors.New("name is too large")
	ErrPriceIsInvalid      = errors.New("price is invalid")
	ErrBoxProductNotExists = errors.New("box product not exist")
)
