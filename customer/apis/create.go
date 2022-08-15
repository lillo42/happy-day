package apis

import (
	"net/http"

	"happyday/common"
	"happyday/customer/applications"
	"happyday/customer/domain"
	"happyday/middlewares"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Name    string   `json:"name,omitempty"`
	Comment string   `json:"comment,omitempty"`
	Phones  []string `json:"phones,omitempty"`
}

func (controller Controller) createEndpoint(context *gin.Context) {
	var req CreateRequest
	if err := context.BindJSON(&req); err != nil {
		middlewares.HandleProblem(context, middlewares.InvalidBody)
		return
	}

	phones := make([]domain.Phone, len(req.Phones))

	for index, phone := range req.Phones {
		phones[index] = domain.NewPhone(phone)
	}

	res, err := controller.createOperation.Execute(context.Request.Context(), applications.CreateRequest{
		Name:    req.Name,
		Comment: req.Comment,
		Phones:  phones,
	})
	if err != nil {
		middlewares.HandleError(context, err)
		return
	}

	body, err := controller.getById(context.Request.Context(), res.ID)
	if err != nil {
		middlewares.HandleError(context, err)
		return
	}

	context.Header(common.ContentType, CustomerV1)
	context.JSON(http.StatusCreated, body)
}
