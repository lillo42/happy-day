package customer

import (
	"context"

	"happy_day/domain/customer"
	"happy_day/infrastructure"
)

type GetByIdHandler struct {
	repository infrastructure.CustomerRepository
}

func (handler GetByIdHandler) Handle(ctx context.Context, req uint) (customer.State, error) {
	return handler.repository.GetById(ctx, req)
}
