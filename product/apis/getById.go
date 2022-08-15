package apis

import (
	"context"
	"net/http"
	"time"

	"happyday/common"
	"happyday/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DetailsResponse struct {
	common.Hateoas
	ID       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Price    float64   `json:"price"`
	IsEnable bool      `json:"isEnable,omitempty"`
	CreateAt time.Time `json:"createAt"`
	ModifyAt time.Time `json:"modifyAt"`
}

func (controller Controller) getById(context context.Context, id uuid.UUID) (*DetailsResponse, error) {
	res, err := controller.readOnlyRepository.GetById(context, id)
	if err != nil {
		return nil, err
	}

	hateoas := common.Hateoas{Links: []common.Link{
		{
			Href:   "/api/products/" + res.ID.String(),
			Method: http.MethodGet,
			Rel:    "read:products",
		},
		{
			Href:   "/api/products/" + res.ID.String(),
			Method: http.MethodPut,
			Rel:    "admin:update:products",
		},
		{
			Href:   "/api/products/" + res.ID.String(),
			Method: http.MethodDelete,
			Rel:    "admin:delete:products",
		},
	}}

	return &DetailsResponse{
		Hateoas:  hateoas,
		ID:       res.ID,
		Name:     res.Name,
		Price:    res.Price,
		IsEnable: res.IsEnable,
		CreateAt: res.CreateAt,
		ModifyAt: res.ModifyAt,
	}, nil
}

func (controller Controller) getByIdEndpoint(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		middlewares.HandleProblem(context, NotFound)
		return
	}

	res, err := controller.getById(context.Request.Context(), id)
	if err != nil {
		middlewares.HandleError(context, err)
		return
	}

	context.Header(common.ContentType, ProductV1)
	context.JSON(http.StatusOK, res)
}
