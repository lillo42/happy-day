package domain

import "github.com/google/uuid"

type (
	Product struct {
		id     uuid.UUID
		amount int64
	}
)

func NewProduct(id uuid.UUID, amount int64) Product {
	return Product{id: id, amount: amount}
}

func (product Product) ID() uuid.UUID { return product.id }
func (product Product) Amount() int64 { return product.amount }
