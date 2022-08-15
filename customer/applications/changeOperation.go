package applications

import (
	"context"

	"happyday/abstract"
	"happyday/common"
	"happyday/customer/domain"
	"happyday/customer/infrastructure"

	"github.com/google/uuid"
)

type (
	ChangeRequest struct {
		ID      uuid.UUID
		Name    string
		Comment string
		Phones  []domain.Phone
	}

	ChangeOperation struct {
		repository infrastructure.Repository
	}
)

var _ abstract.Operation[ChangeRequest, common.VoidResponse] = (*ChangeOperation)(nil)

func (operation ChangeOperation) Execute(ctx context.Context, req ChangeRequest) (common.VoidResponse, error) {
	res := common.VoidResponse{}

	if root, err := operation.repository.Get(ctx, req.ID); err != nil {
		return res, err
	} else if err := root.ChangeName(req.Name); err != nil {
		return res, err
	} else if err := root.ChangeComment(req.Comment); err != nil {
		return res, err
	} else if err := root.ChangePhones(req.Phones); err != nil {
		return res, err
	} else if err := operation.repository.Save(ctx, root); err != nil {
		return res, err
	}

	return res, nil
}
