package apis

import "errors"

var (
	ErrInvalidBody error = errors.New("invalid json body")
)
