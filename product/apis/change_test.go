package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"happyday/abstract"
	"happyday/common"
	"happyday/middlewares"
	"happyday/product/applications"
	"happyday/product/infrastructure"
	"happyday/product/test"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestChangeEndpoint_Should_ReturnError_When_ErrorToParse(t *testing.T) {
	engine := gin.New()
	controller := Controller{}
	controller.MapEndpoint(engine)

	req, _ := http.NewRequest(http.MethodPut, "/api/products/1234", nil)

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestChangeEndpoint_Should_ReturnInvalidBody_When_BodyIsNotJson(t *testing.T) {
	engine := gin.New()
	controller := Controller{}
	controller.MapEndpoint(engine)

	id := uuid.New()
	req, _ := http.NewRequest(http.MethodPut, "/api/products/"+id.String(), bytes.NewBuffer([]byte(`"jsonStr"`)))

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, common.Problem, recorder.Header().Get(common.ContentType))
}

func TestChangeEndpoint_Should_ReturnError_When_ErrorToExecuteChange(t *testing.T) {
	err := errors.New(common.RandString(10))
	middlewares.AddErrors(map[error]common.ProblemDetails{
		err: {
			Status: http.StatusUnprocessableEntity,
		},
	})

	id := uuid.New()
	operation := &abstract.MockOperation[applications.ChangeRequest, common.VoidResponse]{}
	operation.
		On("Execute", mock.Anything, mock.Anything).
		Return(common.VoidResponse{}, err)

	engine := gin.New()
	controller := Controller{changeOperation: operation}
	controller.MapEndpoint(engine)

	b, _ := json.Marshal(ChangeRequest{Name: common.RandString(100), Price: 100.0})
	req, _ := http.NewRequest(http.MethodPut, "/api/products/"+id.String(), bytes.NewBuffer(b))

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	assert.Equal(t, common.Problem, recorder.Header().Get(common.ContentType))
}

func TestChangeEndpoint_Should_ReturnError_When_ErrorToGetById(t *testing.T) {
	err := errors.New(common.RandString(10))
	middlewares.AddErrors(map[error]common.ProblemDetails{
		err: {
			Status: http.StatusNotFound,
		},
	})

	id := uuid.New()
	operation := &abstract.MockOperation[applications.ChangeRequest, common.VoidResponse]{}
	operation.
		On("Execute", mock.Anything, mock.Anything).
		Return(common.VoidResponse{}, nil)

	repository := &test.MockRepository{}
	repository.
		On("GetById", mock.Anything, id).
		Return(infrastructure.DetailsViewModel{}, err)

	engine := gin.New()
	controller := Controller{changeOperation: operation, readOnlyRepository: repository}
	controller.MapEndpoint(engine)

	b, _ := json.Marshal(ChangeRequest{Name: common.RandString(100), Price: 10.1})
	req, _ := http.NewRequest(http.MethodPut, "/api/products/"+id.String(), bytes.NewBuffer(b))

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, common.Problem, recorder.Header().Get(common.ContentType))
}

func TestChangeEndpoint(t *testing.T) {
	err := errors.New(common.RandString(10))
	middlewares.AddErrors(map[error]common.ProblemDetails{
		err: {
			Status: http.StatusNotFound,
		},
	})

	operation := &abstract.MockOperation[applications.ChangeRequest, common.VoidResponse]{}
	operation.
		On("Execute", mock.Anything, mock.Anything).
		Return(common.VoidResponse{}, nil)

	id := uuid.New()
	name := common.RandString(100)
	price := 100.0

	repository := &test.MockRepository{}
	repository.
		On("GetById", mock.Anything, id).
		Return(infrastructure.DetailsViewModel{
			ID:    id,
			Name:  name,
			Price: price,
		}, nil)

	engine := gin.New()
	controller := Controller{changeOperation: operation, readOnlyRepository: repository}
	controller.MapEndpoint(engine)

	b, _ := json.Marshal(ChangeRequest{Name: name, Price: price})
	req, _ := http.NewRequest(http.MethodPut, "/api/products/"+id.String(), bytes.NewBuffer(b))

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, ProductV1, recorder.Header().Get(common.ContentType))
}
