package customer

import (
	"context"
	"happy_day/infrastructure"
)

type (
	DeleteRequest struct {
		Id uint `json:"id"`
	}

	DeleteHandler struct {
		repository infrastructure.CustomerRepository
	}
)

func (handler DeleteHandler) Handle(ctx context.Context, req DeleteRequest) error {
	return handler.repository.Delete(ctx, req.Id)
}
