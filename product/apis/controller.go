package apis

import (
	"net/http"

	"happyday/abstract"
	"happyday/common"
	"happyday/product/applications"
	"happyday/product/domain"
	"happyday/product/infrastructure"

	"github.com/gin-gonic/gin"
)

var NotFound = common.ProblemDetails{
	Type:    "/api/products/not-found",
	Title:   "PRC001",
	Details: "Product not found",
	Status:  http.StatusNotFound,
}

type Controller struct {
	createOperation abstract.Operation[applications.CreateRequest, applications.CreateResponse]
	changeOperation abstract.Operation[applications.ChangeRequest, common.VoidResponse]
	deleteOperation abstract.Operation[applications.DeleteRequest, common.VoidResponse]

	readOnlyRepository infrastructure.ReadOnlyRepository
}

func (controller Controller) MapEndpoint(engine *gin.Engine) {
	group := engine.Group("/api/products")
	group.POST("", controller.createEndpoint)
	group.GET("", controller.getAllEndpoint)
	group.GET("/:id", controller.getByIdEndpoint)
	group.PUT("/:id", controller.changeEndpoint)
	group.DELETE(":id", controller.deleteEndpoint)
}

func (controller Controller) ErrorMapping() map[error]common.ProblemDetails {
	return map[error]common.ProblemDetails{
		domain.NameIsEmpty: {
			Type:    "/api/products/name-is-missing",
			Title:   domain.NameIsEmpty.Code(),
			Details: domain.NameIsEmpty.Message(),
			Status:  http.StatusBadRequest,
		},

		domain.NameIsTooLarge: {
			Type:    "/api/products/name-is-too-larger",
			Title:   domain.NameIsTooLarge.Code(),
			Details: domain.NameIsTooLarge.Message(),
			Status:  http.StatusUnprocessableEntity,
		},

		domain.PriceIsInvalid: {
			Type:    "/api/products/invalid-price",
			Title:   domain.PriceIsInvalid.Code(),
			Details: domain.PriceIsInvalid.Message(),
			Status:  http.StatusUnprocessableEntity,
		},

		infrastructure.NotFound: NotFound,
		infrastructure.ConcurrencyIssue: {
			Type:    "/api/products/conflict",
			Title:   infrastructure.ConcurrencyIssue.Code(),
			Details: infrastructure.ConcurrencyIssue.Message(),
			Status:  http.StatusConflict,
		},
	}
}
