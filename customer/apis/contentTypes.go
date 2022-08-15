package apis

import "happyday/common"

const (
	Customer   = "application/vnd.happy-day.customer+json"
	CustomerV1 = Customer + "; " + common.Version1
)
