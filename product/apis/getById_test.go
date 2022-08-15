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

func TestGetByIdEndpoint_Should_ReturnError_When_ErrorToParse(t *testing.T) {
	engine := gin.New()
	controller := Controller{}
	controller.MapEndpoint(engine)

	req, _ := http.NewRequest(http.MethodGet, "/api/products/1234", nil)

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetByIdEndpoint_Should_ReturnError_When_ErrorToGetById(t *testing.T) {
	err := errors.New(common.RandString(10))
	middlewares.AddErrors(map[error]common.ProblemDetails{
		err: {
			Status: http.StatusNotFound,
		},
	})

	id := uuid.New()

	repository := &test.MockRepository{}
	repository.
		On("GetById", mock.Anything, id).
		Return(infrastructure.DetailsViewModel{}, err)

	engine := gin.New()
	controller := Controller{readOnlyRepository: repository}
	controller.MapEndpoint(engine)

	req, _ := http.NewRequest(http.MethodGet, "/api/products/"+id.String(), nil)

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, common.Problem, recorder.Header().Get(common.ContentType))
}

func TestGetByIdEndpoint(t *testing.T) {
	id := uuid.New()

	repository := &test.MockRepository{}
	repository.
		On("GetById", mock.Anything, id).
		Return(infrastructure.DetailsViewModel{
			ID:    id,
			Name:  common.RandString(10),
			Price: 100,
		}, nil)

	engine := gin.New()
	controller := Controller{readOnlyRepository: repository}
	controller.MapEndpoint(engine)

	req, _ := http.NewRequest(http.MethodGet, "/api/products/"+id.String(), nil)

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, ProductV1, recorder.Header().Get(common.ContentType))
}
