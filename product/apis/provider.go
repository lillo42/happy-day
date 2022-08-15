package apis

import (
	"happyday/abstract"
	"happyday/common"
	"happyday/product/applications"
	"happyday/product/infrastructure"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(ProviderProductController)

func ProviderProductController(
	createOperation abstract.Operation[applications.CreateRequest, applications.CreateResponse],
	changeOperation abstract.Operation[applications.ChangeRequest, common.VoidResponse],
	deleteOperation abstract.Operation[applications.DeleteRequest, common.VoidResponse],
	readOnlyRepository infrastructure.ReadOnlyRepository) Controller {
	return Controller{
		createOperation:    createOperation,
		changeOperation:    changeOperation,
		deleteOperation:    deleteOperation,
		readOnlyRepository: readOnlyRepository,
	}
}
