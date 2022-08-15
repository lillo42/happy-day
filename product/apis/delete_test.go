package apis

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"happyday/abstract"
	"happyday/common"
	"happyday/middlewares"
	"happyday/product/applications"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteEndpoint_Should_ReturnError_When_ErrorToParse(t *testing.T) {
	engine := gin.New()
	controller := Controller{}
	controller.MapEndpoint(engine)

	req, _ := http.NewRequest(http.MethodDelete, "/api/products/1234", nil)

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteEndpoint_Should_ReturnError_When_ErrorToExecute(t *testing.T) {
	id := uuid.New()

	err := errors.New(common.RandString(10))
	middlewares.AddErrors(map[error]common.ProblemDetails{
		err: {
			Status: http.StatusUnprocessableEntity,
		},
	})

	operation := &abstract.MockOperation[applications.DeleteRequest, common.VoidResponse]{}
	operation.
		On("Execute", mock.Anything, applications.DeleteRequest{ID: id}).
		Return(common.VoidResponse{}, err)

	engine := gin.New()
	controller := Controller{deleteOperation: operation}
	controller.MapEndpoint(engine)

	req, _ := http.NewRequest(http.MethodDelete, "/api/products/"+id.String(), nil)

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestDeleteEndpoint(t *testing.T) {
	id := uuid.New()

	operation := &abstract.MockOperation[applications.DeleteRequest, common.VoidResponse]{}
	operation.
		On("Execute", mock.Anything, applications.DeleteRequest{ID: id}).
		Return(common.VoidResponse{}, nil)

	engine := gin.New()
	controller := Controller{deleteOperation: operation}
	controller.MapEndpoint(engine)

	req, _ := http.NewRequest(http.MethodDelete, "/api/products/"+id.String(), nil)

	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNoContent, rec.Code)
}
