package apis

import (
	"net/http"
	"strconv"

	"happyday/common"
	"happyday/customer/infrastructure"
	"happyday/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Response struct {
	common.Hateoas
	ID      uuid.UUID `json:"id,omitempty"`
	Name    string    `json:"name,omitempty"`
	Comment string    `json:"comment,omitempty"`
	Phones  []string  `json:"phones,omitempty"`
}

func (controller Controller) getAll(context *gin.Context) {
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

	for index, customer := range page.Items {
		res.Items[index] = Response{
			Hateoas: common.Hateoas{
				Links: []common.Link{
					{
						Href:   "/api/customers/" + customer.ID.String(),
						Method: http.MethodGet,
						Rel:    "read:customer",
					},
				},
			},
			ID:      customer.ID,
			Name:    customer.Name,
			Comment: customer.Comment,
			Phones:  customer.Phones,
		}
	}

	context.Header(common.ContentType, CustomerV1)
	context.JSON(http.StatusOK, res)
}
