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

func TestCreateEndpoint_Should_ReturnInvalidBody_When_BodyIsNotJson(t *testing.T) {
	engine := gin.New()
	controller := Controller{}
	controller.MapEndpoint(engine)

	req, _ := http.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer([]byte(`"jsonStr"`)))

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, common.Problem, recorder.Header().Get(common.ContentType))
}

func TestCreateEndpoint_Should_ReturnError_When_ErrorToExecuteCreate(t *testing.T) {
	err := errors.New(common.RandString(10))
	middlewares.AddErrors(map[error]common.ProblemDetails{
		err: {
			Status: http.StatusUnprocessableEntity,
		},
	})

	name := common.RandString(100)
	price := 100.0
	priority := int64(1)

	operation := &abstract.MockOperation[applications.CreateRequest, applications.CreateResponse]{}
	operation.
		On("Execute", mock.Anything, applications.CreateRequest{Name: name, Price: price, Priority: priority}).
		Return(applications.CreateResponse{}, err)

	engine := gin.New()
	controller := Controller{createOperation: operation}
	controller.MapEndpoint(engine)

	b, _ := json.Marshal(CreateRequest{Name: name, Price: price})
	req, _ := http.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(b))

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	assert.Equal(t, common.Problem, recorder.Header().Get(common.ContentType))
}

func TestCreateEndpoint_Should_ReturnError_When_ErrorToGetById(t *testing.T) {
	err := errors.New(common.RandString(10))
	middlewares.AddErrors(map[error]common.ProblemDetails{
		err: {
			Status: http.StatusNotFound,
		},
	})

	name := common.RandString(100)
	price := 100.0

	operation := &abstract.MockOperation[applications.CreateRequest, applications.CreateResponse]{}
	operation.
		On("Execute", mock.Anything, applications.CreateRequest{Name: name, Price: price}).
		Return(applications.CreateResponse{ID: uuid.New()}, nil)

	repository := &test.MockRepository{}
	repository.
		On("GetById", mock.Anything, mock.Anything).
		Return(infrastructure.DetailsViewModel{}, err)

	engine := gin.New()
	controller := Controller{createOperation: operation, readOnlyRepository: repository}
	controller.MapEndpoint(engine)

	b, _ := json.Marshal(CreateRequest{Name: name, Price: price})
	req, _ := http.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(b))

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, common.Problem, recorder.Header().Get(common.ContentType))
}

func TestCreateEndpoint(t *testing.T) {
	err := errors.New(common.RandString(10))
	middlewares.AddErrors(map[error]common.ProblemDetails{
		err: {
			Status: http.StatusNotFound,
		},
	})

	name := common.RandString(100)
	price := 100.0

	operation := &abstract.MockOperation[applications.CreateRequest, applications.CreateResponse]{}
	operation.
		On("Execute", mock.Anything, applications.CreateRequest{Name: name, Price: price}).
		Return(applications.CreateResponse{ID: uuid.New()}, nil)

	viewModel := infrastructure.DetailsViewModel{
		ID:       uuid.New(),
		Name:     name,
		Price:    price,
		IsEnable: false,
	}

	repository := &test.MockRepository{}
	repository.
		On("GetById", mock.Anything, mock.Anything).
		Return(viewModel, nil)

	engine := gin.New()
	controller := Controller{createOperation: operation, readOnlyRepository: repository}
	controller.MapEndpoint(engine)

	b, _ := json.Marshal(CreateRequest{Name: name, Price: price})
	req, _ := http.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(b))

	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)
	assert.Equal(t, ProductV1, recorder.Header().Get(common.ContentType))
}
