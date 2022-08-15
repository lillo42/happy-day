package apis

import (
	"net/http"

	"happyday/common"
	"happyday/customer/applications"
	"happyday/customer/domain"
	"happyday/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChangeRequest struct {
	Name    string
	Comment string
	Phones  []string
}

func (controller Controller) changeEndpoint(context *gin.Context) {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		middlewares.HandleProblem(context, NotFound)
		return
	}

	var req ChangeRequest
	if err = context.BindJSON(&req); err != nil {
		middlewares.HandleProblem(context, middlewares.InvalidBody)
		return
	}

	phones := make([]domain.Phone, len(req.Phones))
	for index, phone := range req.Phones {
		phones[index] = domain.NewPhone(phone)
	}

	_, err = controller.changeOperation.Execute(context.Request.Context(), applications.ChangeRequest{
		ID:      id,
		Name:    req.Name,
		Comment: req.Comment,
		Phones:  phones,
	})

	if err != nil {
		middlewares.HandleError(context, err)
		return
	}

	res, err := controller.getById(context.Request.Context(), id)
	if err != nil {
		middlewares.HandleError(context, err)
		return
	}

	context.Header(common.ContentType, CustomerV1)
	context.JSON(http.StatusOK, res)
}
