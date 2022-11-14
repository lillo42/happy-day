package customer

import (
	"context"

	"happy_day/domain/customer"
	"happy_day/infrastructure"

	"github.com/google/uuid"
)

type GetByIdHandler struct {
	repository infrastructure.CustomerRepository
}

func (handler GetByIdHandler) Handle(ctx context.Context, req uuid.UUID) (customer.State, error) {
	return handler.repository.GetById(ctx, req)
}
