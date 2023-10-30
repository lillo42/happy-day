package discounts

import (
	"github.com/google/uuid"
	"time"
)

type (
	Discount struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		Price    float64   `json:"price"`
		Products []Product `json:"products"`
		CreateAt time.Time `json:"createAt"`
		UpdateAt time.Time `json:"updateAt"`
		Version  uint      `json:"-"`
	}

	Product struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		Quantity uint      `json:"quantity"`
	}
)
