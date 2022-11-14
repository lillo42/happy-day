package customer

import (
	"context"

	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type (
	DeleteRequest struct {
		Id uuid.UUID `json:"id"`
	}

	DeleteHandler struct {
		repository infrastructure.CustomerRepository
	}
)

func (handler DeleteHandler) Handle(ctx context.Context, req DeleteRequest) error {
	return handler.repository.Delete(ctx, req.Id)
}
