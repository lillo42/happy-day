package application

import (
	"context"
	"happy_day/domain/customer"
	"happy_day/infrastructure"
)

type GetAllCustomersHandler struct {
	repository infrastructure.CustomerRepository
}

func (handler GetAllCustomersHandler) Handle(ctx context.Context, req infrastructure.CustomerFilter) (infrastructure.Page[customer.State], error) {
	if req.Size <= 0 {
		req.Size = 50
	}

	if req.Page <= 1 {
		req.Page = 1
	}

	return handler.repository.GetAll(ctx, req)
}
