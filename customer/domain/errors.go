package domain

import "happyday/common"

var (
	NameIsEmpty    = common.NewError("CST000", "Name is empty")
	NameIsTooLarge = common.NewError("CST001", "Name length is to large")

	CommentTooLarge = common.NewError("CST002", "Comment length is to large")

	PhonesIsEmpty        = common.NewError("CST003", "Phones is missing")
	PhoneLengthIsInvalid = common.NewError("CST004", "invalid phone length")
	PhoneNumberIsInvalid = common.NewError("CST005", "invalid phone")
)
