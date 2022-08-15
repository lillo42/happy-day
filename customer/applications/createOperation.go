package applications

import (
	"context"

	"happyday/abstract"
	"happyday/customer/domain"
	"happyday/customer/infrastructure"

	"github.com/google/uuid"
)

type (
	CreateRequest struct {
		Name    string
		Comment string
		Phones  []domain.Phone
	}

	CreateResponse struct {
		ID uuid.UUID
	}

	CreateOperation struct {
		repository infrastructure.Repository
	}
)

var _ abstract.Operation[CreateRequest, CreateResponse] = (*CreateOperation)(nil)

func (operation CreateOperation) Execute(ctx context.Context, req CreateRequest) (CreateResponse, error) {
	res := CreateResponse{
		ID: uuid.New(),
	}

	root := operation.repository.Create(res.ID)
	if err := root.Create(req.Name, req.Comment, req.Phones); err != nil {
		return res, err
	}

	if err := operation.repository.Save(ctx, root); err != nil {
		return res, err
	}

	return res, nil
}
