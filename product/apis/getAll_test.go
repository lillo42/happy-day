package apis

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"happyday/common"
	"happyday/middlewares"
	"happyday/product/infrastructure"
	"happyday/product/test"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllEndpoint_Should_ReturnError_When_GetAllReturnError(t *testing.T) {
	err := errors.New(common.RandString(10))
	middlewares.AddErrors(map[error]common.ProblemDetails{
		err: {
			Status: http.StatusNotFound,
		},
	})

	repository := &test.MockRepository{}
	repository.
		On("GetAll", mock.Anything, mock.Anything).
		Return(common.Page[infrastructure.ViewModel]{}, err)

	engine := gin.New()
	controller := Controller{readOnlyRepository: repository}
	controller.MapEndpoint(engine)

	req, _ := http.NewRequest(http.MethodGet, "/api/products?name=test&id=info&page=-1&size=-1", nil)

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, common.Problem, recorder.Header().Get(common.ContentType))
}

func TestGetAllEndpoint(t *testing.T) {
	err := errors.New(common.RandString(10))
	middlewares.AddErrors(map[error]common.ProblemDetails{
		err: {
			Status: http.StatusNotFound,
		},
	})

	repository := &test.MockRepository{}
	repository.
		On("GetAll", mock.Anything, mock.Anything).
		Return(common.Page[infrastructure.ViewModel]{
			TotalElements: 1,
			TotalPages:    1,
			Items: []infrastructure.ViewModel{
				{
					ID:       uuid.New(),
					Name:     common.RandString(10),
					IsEnable: false,
				},
			},
		}, nil)

	engine := gin.New()
	controller := Controller{readOnlyRepository: repository}
	controller.MapEndpoint(engine)

	req, _ := http.NewRequest(http.MethodGet, "/api/products?name=test&id=info&page=-1&size=-1", nil)

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, ProductV1, recorder.Header().Get(common.ContentType))
}
