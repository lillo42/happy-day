package applications

import (
	"context"

	"happyday/abstract"
	"happyday/product/domain"
	"happyday/product/infrastructure"

	"github.com/google/uuid"
)

type (
	CreateOperation struct {
		repository         infrastructure.Repository
		readOnlyRepository infrastructure.ReadOnlyRepository
	}

	CreateRequest struct {
		Name     string
		Price    float64
		Priority int64
		Products []domain.Product
	}

	CreateResponse struct {
		ID uuid.UUID
	}
)

var _ abstract.Operation[CreateRequest, CreateResponse] = (*CreateOperation)(nil)

func (operation CreateOperation) Execute(ctx context.Context, req CreateRequest) (CreateResponse, error) {
	res := CreateResponse{}
	root := operation.repository.Create(uuid.New())

	if err := root.Create(req.Name, req.Price, req.Priority, req.Products, func(product domain.Product) (bool, error) {
		return operation.readOnlyRepository.Exists(ctx, product.ID())
	}); err != nil {
		return res, err
	}

	if err := operation.repository.Save(ctx, root); err != nil {
		return res, err
	}

	res.ID = root.State().ID()
	return res, nil
}
