package applications

import (
	"happyday/abstract"
	"happyday/common"
	"happyday/customer/infrastructure"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideCreateOperation,
	wire.Bind(new(abstract.Operation[CreateRequest, CreateResponse]), new(*CreateOperation)),

	ProvideDeleteOperation,
	wire.Bind(new(abstract.Operation[DeleteRequest, common.VoidResponse]), new(*DeleteOperation)),

	ProvideChangeOperation,
	wire.Bind(new(abstract.Operation[ChangeRequest, common.VoidResponse]), new(*ChangeOperation)),
)

func ProvideCreateOperation(repository infrastructure.Repository) *CreateOperation {
	return &CreateOperation{
		repository: repository,
	}
}

func ProvideChangeOperation(repository infrastructure.Repository) *ChangeOperation {
	return &ChangeOperation{
		repository: repository,
	}
}

func ProvideDeleteOperation(repository infrastructure.Repository) *DeleteOperation {
	return &DeleteOperation{
		repository: repository,
	}
}
