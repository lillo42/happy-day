package applications

import (
	"happyday/abstract"
	"happyday/common"
	"happyday/product/infrastructure"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProviderCreateOperation,
	wire.Bind(new(abstract.Operation[CreateRequest, CreateResponse]), new(*CreateOperation)),

	ProviderChangeOperation,
	wire.Bind(new(abstract.Operation[ChangeRequest, common.VoidResponse]), new(*ChangeOperation)),

	ProviderDeleteOperation,
	wire.Bind(new(abstract.Operation[DeleteRequest, common.VoidResponse]), new(*DeleteOperation)),
)

func ProviderCreateOperation(repository infrastructure.Repository,
	readOnlyRepository infrastructure.ReadOnlyRepository) *CreateOperation {
	return &CreateOperation{
		repository:         repository,
		readOnlyRepository: readOnlyRepository,
	}
}

func ProviderChangeOperation(repository infrastructure.Repository,
	readOnlyRepository infrastructure.ReadOnlyRepository) *ChangeOperation {
	return &ChangeOperation{
		repository:         repository,
		readOnlyRepository: readOnlyRepository,
	}
}

func ProviderDeleteOperation(repository infrastructure.Repository) *DeleteOperation {
	return &DeleteOperation{
		repository: repository,
	}
}
