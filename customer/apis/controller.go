package apis

import (
	"net/http"

	"happyday/abstract"
	"happyday/common"
	"happyday/customer/applications"
	"happyday/customer/domain"
	"happyday/customer/infrastructure"

	"github.com/gin-gonic/gin"
)

var (
	NotFound = common.ProblemDetails{
		Type:    "/api/customers/errors/not-found",
		Title:   "PDC000",
		Details: "Product not found",
		Status:  http.StatusNotFound,
	}
)

type Controller struct {
	createOperation abstract.Operation[applications.CreateRequest, applications.CreateResponse]
	changeOperation abstract.Operation[applications.ChangeRequest, common.VoidResponse]
	deleteOperation abstract.Operation[applications.DeleteRequest, common.VoidResponse]

	readOnlyRepository infrastructure.ReadOnlyRepository
}

func (controller Controller) MapEndpoints(engine *gin.Engine) {
	group := engine.Group("/api/customers")
	group.POST("", controller.createEndpoint)
	group.GET("", controller.getAll)
	group.GET("/:id", controller.getByIdEndpoint)
	group.PUT("/:id", controller.changeEndpoint)
	group.DELETE("/:id", controller.deleteEndpoint)
}

func (controller Controller) ErrorMapping() map[error]common.ProblemDetails {
	return map[error]common.ProblemDetails{
		domain.NameIsEmpty: {
			Type:    "/api/customers/name-is-missing",
			Title:   domain.NameIsEmpty.Code(),
			Details: domain.NameIsEmpty.Message(),
			Status:  http.StatusBadRequest,
		},

		domain.NameIsTooLarge: {
			Type:    "/api/customer/name-is-too-larger",
			Title:   domain.NameIsTooLarge.Code(),
			Details: domain.NameIsTooLarge.Message(),
			Status:  http.StatusUnprocessableEntity,
		},

		domain.CommentTooLarge: {
			Type:    "/api/customer/comment-is-too-larger",
			Title:   domain.CommentTooLarge.Code(),
			Details: domain.CommentTooLarge.Message(),
			Status:  http.StatusUnprocessableEntity,
		},

		domain.PhonesIsEmpty: {
			Type:    "/api/customer/phones-is-too-larger",
			Title:   domain.PhonesIsEmpty.Code(),
			Details: domain.PhonesIsEmpty.Message(),
			Status:  http.StatusBadRequest,
		},

		domain.PhoneLengthIsInvalid: {
			Type:    "/api/customer/phones-number-length-is-invalid",
			Title:   domain.PhoneLengthIsInvalid.Code(),
			Details: domain.PhoneLengthIsInvalid.Message(),
			Status:  http.StatusUnprocessableEntity,
		},

		domain.PhoneNumberIsInvalid: {
			Type:    "/api/customer/phones-number-is-invalid",
			Title:   domain.PhoneNumberIsInvalid.Code(),
			Details: domain.PhoneNumberIsInvalid.Message(),
			Status:  http.StatusUnprocessableEntity,
		},

		infrastructure.NotFound: NotFound,
		infrastructure.ConcurrencyIssue: {
			Type:    "/api/customer/conflict",
			Title:   infrastructure.ConcurrencyIssue.Code(),
			Details: infrastructure.ConcurrencyIssue.Message(),
			Status:  http.StatusConflict,
		},
	}
}
