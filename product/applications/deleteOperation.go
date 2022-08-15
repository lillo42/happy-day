package applications

import (
	"context"

	"happyday/abstract"
	"happyday/common"
	"happyday/product/infrastructure"

	"github.com/google/uuid"
)

type (
	DeleteRequest struct{ ID uuid.UUID }

	DeleteOperation struct {
		repository infrastructure.Repository
	}
)

var _ abstract.Operation[DeleteRequest, common.VoidResponse] = (*DeleteOperation)(nil)

func (o DeleteOperation) Execute(ctx context.Context, req DeleteRequest) (common.VoidResponse, error) {
	res := common.VoidResponse{}
	root, err := o.repository.Get(ctx, req.ID)

	if err != nil {
		return res, err
	}

	if err = o.repository.Delete(ctx, root); err != nil {
		return res, err
	}

	return res, nil
}
