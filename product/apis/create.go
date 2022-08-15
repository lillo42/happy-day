package apis

import (
	"net/http"

	"happyday/common"
	"happyday/middlewares"
	"happyday/product/applications"
	"happyday/product/domain"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Name     string           `json:"name,omitempty"`
	Price    float64          `json:"price"`
	Priority int64            `json:"priority"`
	Products []ProductRequest `json:"products"`
}

func (controller Controller) createEndpoint(context *gin.Context) {
	var req CreateRequest
	if err := context.BindJSON(&req); err != nil {
		middlewares.HandleProblem(context, middlewares.InvalidBody)
		return
	}

	products := make([]domain.Product, len(req.Products))
	for i, product := range req.Products {
		products[i] = domain.NewProduct(product.ID, product.Amount)
	}

	res, err := controller.createOperation.Execute(context.Request.Context(), applications.CreateRequest{
		Name:     req.Name,
		Price:    req.Price,
		Priority: req.Priority,
		Products: products,
	})

	if err != nil {
		middlewares.HandleError(context, err)
		return
	}

	viewModel, err := controller.getById(context.Request.Context(), res.ID)
	if err != nil {
		middlewares.HandleError(context, err)
		return
	}

	context.Header(common.ContentType, ProductV1)
	context.JSON(http.StatusCreated, viewModel)
}
