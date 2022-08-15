package apis

import (
	"net/http"

	"happyday/middlewares"
	"happyday/product/applications"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (controller Controller) deleteEndpoint(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))

	if err != nil {
		middlewares.HandleProblem(context, NotFound)
		return
	}

	_, err = controller.deleteOperation.Execute(context.Request.Context(), applications.DeleteRequest{ID: id})
	if err != nil {
		middlewares.HandleError(context, err)
		return
	}

	context.Status(http.StatusNoContent)
}
