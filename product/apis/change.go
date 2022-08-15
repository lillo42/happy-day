package apis

import (
	"net/http"

	"happyday/common"
	"happyday/middlewares"
	"happyday/product/applications"
	"happyday/product/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChangeRequest struct {
	Name     string           `json:"name,omitempty"`
	Price    float64          `json:"price,omitempty"`
	IsEnable bool             `json:"isEnable"`
	Priority int64            `json:"priority"`
	Products []ProductRequest `json:"products"`
}

func (controller Controller) changeEndpoint(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		middlewares.HandleProblem(context, NotFound)
		return
	}

	var req ChangeRequest
	if err := context.BindJSON(&req); err != nil {
		middlewares.HandleProblem(context, middlewares.InvalidBody)
		return
	}

	products := make([]domain.Product, len(req.Products))
	for i, product := range req.Products {
		products[i] = domain.NewProduct(product.ID, product.Amount)
	}
	_, err = controller.changeOperation.Execute(context.Request.Context(), applications.ChangeRequest{
		ID:       id,
		Name:     req.Name,
		Price:    req.Price,
		IsEnable: req.IsEnable,
		Priority: req.Priority,
		Products: products,
	})

	if err != nil {
		middlewares.HandleError(context, err)
		return
	}

	viewModel, err := controller.getById(context.Request.Context(), id)
	if err != nil {
		middlewares.HandleError(context, err)
		return
	}

	context.Header(common.ContentType, ProductV1)
	context.JSON(http.StatusOK, viewModel)
}
