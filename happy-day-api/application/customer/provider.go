package customer

import (
	"happy_day/infrastructure"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	ProvideGetAllHandler,
	ProvideGetByIdHandler,
	ProvideChangeOrCreateHandler,
	ProvideDeleteHandler,
)

func ProvideGetAllHandler(repository infrastructure.CustomerRepository) GetAllHandler {
	return GetAllHandler{repository: repository}
}

func ProvideGetByIdHandler(repository infrastructure.CustomerRepository) GetByIdHandler {
	return GetByIdHandler{repository: repository}
}

func ProvideChangeOrCreateHandler(repository infrastructure.CustomerRepository) ChangeOrCreateHandler {
	return ChangeOrCreateHandler{repository: repository}
}

func ProvideDeleteHandler(repository infrastructure.CustomerRepository) DeleteHandler {
	return DeleteHandler{repository: repository}
}
