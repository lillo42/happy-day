package product

import (
	"time"

	"github.com/google/uuid"
)

type (
	State struct {
		Id         uuid.UUID `bson:"id" json:"id,omitempty"`
		Name       string    `bson:"name" json:"name"`
		Price      float64   `bson:"price" json:"price"`
		Products   []Product `bson:"products,omitempty" json:"products"`
		CreatedAt  time.Time `bson:"createdAt" json:"createdAt"`
		ModifiedAt time.Time `bson:"modifiedAt" json:"modifiedAt"`
	}

	Product struct {
		Id       uuid.UUID `bson:"id" json:"id"`
		Quantity int64     `bson:"quantity" json:"quantity"`
	}
)
