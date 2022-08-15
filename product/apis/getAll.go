package apis

import (
	"net/http"
	"strconv"

	"happyday/common"
	"happyday/middlewares"
	"happyday/product/infrastructure"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	common.Hateoas
	ID       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Price    float64   `json:"price"`
	IsEnable bool      `json:"isEnable,omitempty"`
}

func (controller Controller) getAllEndpoint(context *gin.Context) {
	filter := infrastructure.Filter{}
	filter.Text = context.Query("text")

	var err error
	filter.Page, err = strconv.ParseInt(context.Query("page"), 10, 64)
	if err != nil || filter.Page < 1 {
		filter.Page = 1
	}

	filter.Size, err = strconv.ParseInt(context.Query("size"), 10, 64)
	if err != nil || filter.Size < 1 {
		filter.Size = 50
	}

	filter.OrderBy = infrastructure.NameAsc
	orderBy := context.Query("orderBy")
	if len(orderBy) > 0 {
		filter.OrderBy = infrastructure.OrderBy(orderBy)
	}

	page, err := controller.readOnlyRepository.GetAll(context.Request.Context(), filter)
	if err != nil {
		middlewares.HandleError(context, err)
		return
	}

	res := common.Page[Response]{
		TotalElements: page.TotalElements,
		TotalPages:    page.TotalPages,
		Items:         make([]Response, len(page.Items)),
	}

	for index, product := range page.Items {
		res.Items[index] = Response{
			Hateoas: common.Hateoas{
				Links: []common.Link{
					{
						Href:   "/api/products/" + product.ID.String(),
						Method: http.MethodGet,
						Rel:    "read:product",
					},
				},
			},
			ID:       product.ID,
			Name:     product.Name,
			Price:    product.Price,
			IsEnable: product.IsEnable,
		}
	}

	context.Header(common.ContentType, ProductV1)
	context.JSON(http.StatusOK, res)
}
