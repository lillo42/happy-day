package domain

import "happyday/common"

var (
	NameIsEmpty    = common.NewError("PRD000", "Missing name")
	NameIsTooLarge = common.NewError("PRD001", "Name length is too large")

	PriceIsInvalid = common.NewError("PRD002", "Invalid price amount")

	ProductNotExist = common.NewError("PRD003", "Product not exist")
)
