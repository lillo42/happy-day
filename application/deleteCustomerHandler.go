package application

import (
	"context"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type (
	DeleteCustomerRequest struct {
		Id uuid.UUID `json:"id"`
	}

	DeleteCustomerHandler struct {
		repository infrastructure.CustomerRepository
	}
)

func (handler DeleteCustomerHandler) Handle(ctx context.Context, req DeleteCustomerRequest) error {
	return handler.repository.Delete(ctx, req.Id)
}
