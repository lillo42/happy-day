package customer

import (
	"context"

	"happy_day/domain/customer"
	"happy_day/infrastructure"
)

type GetAllHandler struct {
	repository infrastructure.CustomerRepository
}

func (handler GetAllHandler) Handle(ctx context.Context, req infrastructure.CustomerFilter) (infrastructure.Page[customer.State], error) {
	if req.Size <= 0 {
		req.Size = 50
	}

	if req.Page <= 1 {
		req.Page = 1
	}

	return handler.repository.GetAll(ctx, req)
}
