package products

import (
	"github.com/google/uuid"
	"time"
)

type (
	Product struct {
		ID       uuid.UUID    `json:"id"`
		Name     string       `json:"name"`
		Price    float64      `json:"price"`
		Products []BoxProduct `json:"products"`
		CreateAt time.Time    `json:"createAt"`
		UpdateAt time.Time    `json:"updateAt"`
		Version  uint         `json:"-"`
	}

	BoxProduct struct {
		ID       uuid.UUID `json:"id"`
		Quantity uint      `json:"quantity"`
	}
)
