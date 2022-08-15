package apis

import (
	"happyday/abstract"
	"happyday/common"
	"happyday/customer/applications"
	"happyday/customer/infrastructure"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(ProviderController)

func ProviderController(
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
