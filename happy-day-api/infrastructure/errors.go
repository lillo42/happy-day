package infrastructure

import "errors"

var (
	ErrCustomerNameIsEmpty    = errors.New("customer name is empty")
	ErrCustomerPhonesIsEmpty  = errors.New("customer phones is empty")
	ErrCustomerPhoneIsInvalid = errors.New("customer phone is invalid")

	ErrProductNameIsEmpty               = errors.New("product name is empty")
	ErrProductPriceIsLessThanZero       = errors.New("product price is less than zero")
	ErrProductAmountIsInvalid           = errors.New("product amount is invalid")
	ErrExistOtherProductWithThisProduct = errors.New("exist other product with this product")

	ErrReservationPaymentInstallmentAmount = errors.New("payment installment amount cannot be less or equal to zero")
	ErrReservationAddressCityIsEmpty       = errors.New("address city cannot be empty")
	ErrReservationAddressStreetIsEmpty     = errors.New("address street cannot be empty")
	ErrReservationAddressNumberIsInvalid   = errors.New("address number cannot be empty")
	ErrReservationAddressPostalCodeIsEmpty = errors.New("address postal code cannot be empty")
	ErrProductListIsEmpty                  = errors.New("product list cannot be empty")
)
