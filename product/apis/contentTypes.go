package apis

import "happyday/common"

const (
	Product = "application/vnd.happy-day.product+json"

	ProductV1 = Product + "; " + common.Version1
)
