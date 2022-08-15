package applications

import (
	"context"

	"happyday/abstract"
	"happyday/common"
	"happyday/customer/infrastructure"

	"github.com/google/uuid"
)

type (
	DeleteRequest struct {
		ID uuid.UUID
	}

	DeleteOperation struct {
		repository infrastructure.Repository
	}
)

var _ abstract.Operation[DeleteRequest, common.VoidResponse] = (*DeleteOperation)(nil)

func (operation DeleteOperation) Execute(ctx context.Context, req DeleteRequest) (common.VoidResponse, error) {
	res := common.VoidResponse{}

	if root, err := operation.repository.Get(ctx, req.ID); err != nil {
		return res, err
	} else if err := operation.repository.Delete(ctx, root); err != nil {
		return res, err
	}

	return res, nil
}
