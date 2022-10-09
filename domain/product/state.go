package product

import "github.com/google/uuid"

type (
	State struct {
		Id       uuid.UUID `bson:"id" json:"id,omitempty"`
		Name     string    `bson:"name" json:"name"`
		Price    float64   `bson:"price" json:"price"`
		Products []Product `bson:"products,omitempty" json:"products"`
	}

	Product struct {
		Id     uuid.UUID `bson:"id" json:"id"`
		Amount int64     `bson:"amount" json:"amount"`
	}
)
