package applications

import (
	"context"

	"happyday/abstract"
	"happyday/common"
	"happyday/product/domain"
	"happyday/product/infrastructure"

	"github.com/google/uuid"
)

type (
	ChangeRequest struct {
		ID       uuid.UUID
		Name     string
		Price    float64
		IsEnable bool
		Priority int64
		Products []domain.Product
	}

	ChangeOperation struct {
		repository         infrastructure.Repository
		readOnlyRepository infrastructure.ReadOnlyRepository
	}
)

var _ abstract.Operation[ChangeRequest, common.VoidResponse] = (*ChangeOperation)(nil)

func (operation ChangeOperation) Execute(ctx context.Context, req ChangeRequest) (common.VoidResponse, error) {
	root, err := operation.repository.Get(ctx, req.ID)

	if err == nil {
		err = root.ChangeName(req.Name)
	}

	if err == nil {
		err = root.ChangePrice(req.Price)
	}

	if err == nil {
		err = root.ChangePriority(req.Priority)
	}

	if err == nil {
		err = root.ChangeProducts(req.Products, func(product domain.Product) (bool, error) {
			return operation.readOnlyRepository.Exists(ctx, product.ID())
		})
	}

	if err == nil {
		if req.IsEnable {
			err = root.Enable()
		} else {
			err = root.Disable()
		}
	}

	if err != nil {
		return common.VoidResponse{}, err
	}

	if err = operation.repository.Save(ctx, root); err != nil {
		return common.VoidResponse{}, err
	}

	return common.VoidResponse{}, nil
}
