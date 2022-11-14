package product

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

func ProvideGetAllHandler(repository infrastructure.ProductRepository) GetAllHandler {
	return GetAllHandler{repository: repository}
}

func ProvideGetByIdHandler(repository infrastructure.ProductRepository) GetByIdHandler {
	return GetByIdHandler{repository: repository}
}

func ProvideChangeOrCreateHandler(repository infrastructure.ProductRepository) ChangeOrCreateHandler {
	return ChangeOrCreateHandler{repository: repository}
}

func ProvideDeleteHandler(repository infrastructure.ProductRepository) DeleteHandler {
	return DeleteHandler{repository: repository}
}
