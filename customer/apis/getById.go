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
	Comment  string    `json:"comment,omitempty"`
	Phones   []string  `json:"phones,omitempty"`
	CreateAt time.Time `json:"createAt"`
	ModifyAt time.Time `json:"modifyAt"`
}

func (controller Controller) getById(ctx context.Context, id uuid.UUID) (DetailsResponse, error) {
	viewModel, err := controller.readOnlyRepository.GetById(ctx, id)
	if err != nil {
		return DetailsResponse{}, err
	}

	res := DetailsResponse{
		ID:       viewModel.ID,
		Name:     viewModel.Name,
		Comment:  viewModel.Comment,
		Phones:   viewModel.Phones,
		CreateAt: viewModel.CreateAt,
		ModifyAt: viewModel.ModifyAt,
		Hateoas: common.Hateoas{Links: []common.Link{
			{
				Href:   "/api/customer/" + viewModel.ID.String(),
				Method: http.MethodGet,
				Rel:    "read:customer",
			},
			{
				Href:   "/api/customer/" + viewModel.ID.String(),
				Method: http.MethodPut,
				Rel:    "update:customer",
			},
			{
				Href:   "/api/customer/" + viewModel.ID.String(),
				Method: http.MethodDelete,
				Rel:    "delete:customer",
			},
		}},
	}

	return res, nil
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

	context.Header(common.ContentType, CustomerV1)
	context.JSON(http.StatusOK, res)
}
