package application

import (
	"context"
	"happy_day/domain/customer"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type GetCustomerByIdHandler struct {
	repository infrastructure.CustomerRepository
}

func (handler GetCustomerByIdHandler) Handle(ctx context.Context, req uuid.UUID) (customer.State, error) {
	return handler.repository.GetById(ctx, req)
}
