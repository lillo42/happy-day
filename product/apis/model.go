package apis

import "github.com/google/uuid"

type ProductRequest struct {
	ID     uuid.UUID `json:"id,omitempty"`
	Amount int64     `json:"amount,omitempty"`
}
