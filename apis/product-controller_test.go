package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"happy_day/application"
	"happy_day/common"
	"happy_day/domain/product"
	"happy_day/infrastructure"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductControllerCreateWhenIsNotJson(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(common.RandString(10)))

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/products")

	controller := ProductController{}
	err := controller.create(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidBody, err)
}

func TestProductControllerCreateWhenErrOnHandler(t *testing.T) {
	b, _ := json.Marshal(application.CreateOrChangeProductRequest{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/products")

	controller := ProductController{createOrChangeHandler: application.CreateOrChangeProductHandler{}}
	err := controller.create(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, application.ErrProductNameIsEmpty, err)
}

func TestProductControllerCreate(t *testing.T) {
	b, _ := json.Marshal(application.CreateOrChangeProductRequest{
		State: product.State{
			Name:     common.RandString(10),
			Price:    2.5,
			Products: []product.Product{},
		},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/products")

	repo := &infrastructure.MockProductRepository{}
	repo.
		On("Save", mock.Anything, mock.Anything).
		Return(product.State{}, nil)

	controller := ProductController{createOrChangeHandler: application.ProvideCreateOrChangeProductHandler(repo)}

	err := controller.create(ctx)
	assert.Nil(t, err)
}

func TestProductControllerUpdateWhenIsNotJson(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(common.RandString(10)))

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/products/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	controller := ProductController{}
	err := controller.update(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidBody, err)
}

func TestProductControllerUpdateWhenIdIsNotUuid(t *testing.T) {
	b, _ := json.Marshal(application.CreateOrChangeProductRequest{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/products/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("123")

	controller := ProductController{createOrChangeHandler: application.CreateOrChangeProductHandler{}}
	err := controller.update(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, ErrProductNotFound, err)
}

func TestProductControllerUpdateWhenErrOnHandler(t *testing.T) {
	b, _ := json.Marshal(application.CreateOrChangeProductRequest{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/products/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	expectedErr := errors.New(common.RandString(10))
	repo := &infrastructure.MockProductRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(product.State{}, expectedErr)

	controller := ProductController{createOrChangeHandler: application.ProvideCreateOrChangeProductHandler(repo)}
	err := controller.update(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestProductControllerUpdate(t *testing.T) {
	b, _ := json.Marshal(application.CreateOrChangeProductRequest{
		State: product.State{
			Name:  common.RandString(10),
			Price: 2.5,
		},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/products/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	repo := &infrastructure.MockProductRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(product.State{
			Id:    uuid.New(),
			Name:  common.RandString(10),
			Price: 10,
		}, nil)

	repo.
		On("Save", mock.Anything, mock.Anything).
		Return(product.State{
			Id:    uuid.New(),
			Name:  common.RandString(10),
			Price: 3,
		}, nil)

	controller := ProductController{createOrChangeHandler: application.ProvideCreateOrChangeProductHandler(repo)}
	err := controller.update(ctx)

	assert.Nil(t, err)
}

func TestProductControllerDeleteWhenIdIsNotUuid(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/products/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("123")

	controller := ProductController{}
	err := controller.delete(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, ErrProductNotFound, err)
}

func TestProductControllerDeleteWhenErrOnHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/products/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	repo := &infrastructure.MockProductRepository{}
	repo.
		On("ExistAnyWithProduct", mock.Anything, mock.Anything).
		Return(true, nil)

	controller := ProductController{deleteHandler: application.ProvideDeleteProductHandler(repo)}
	err := controller.delete(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, application.ErrExistOtherProductWithThisProduct, err)
}

func TestProductControllerDelete(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/products/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	repo := &infrastructure.MockProductRepository{}
	repo.
		On("ExistAnyWithProduct", mock.Anything, mock.Anything).
		Return(false, nil)

	repo.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil)

	controller := ProductController{deleteHandler: application.ProvideDeleteProductHandler(repo)}
	err := controller.delete(ctx)

	assert.Nil(t, err)
}

func TestProductControllerGetWhenIdIsNotUuid(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/products/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("123")

	controller := ProductController{}
	err := controller.get(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, ErrProductNotFound, err)
}

func TestProductControllerGetWhenErrToGetById(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/products/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	expectedErr := errors.New(common.RandString(10))
	repo := &infrastructure.MockProductRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(product.State{}, expectedErr)

	controller := ProductController{repository: repo}
	err := controller.get(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestProductControllerGet(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/products/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	repo := &infrastructure.MockProductRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(product.State{}, nil)

	controller := ProductController{repository: repo}
	err := controller.get(ctx)

	assert.Nil(t, err)
}

func TestProductControllerGetAllWhenErrOnHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/products")

	expectedErr := errors.New(common.RandString(10))
	repo := &infrastructure.MockProductRepository{}
	repo.
		On("GetAll", mock.Anything, mock.Anything).
		Return(infrastructure.Page[product.State]{}, expectedErr)

	controller := ProductController{getAllHandler: application.ProvideGetAllProductsHandler(repo)}
	err := controller.getAll(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestProductControllerGetAll(t *testing.T) {
	q := make(url.Values)
	q.Set("text", common.RandString(10))
	q.Set("page", strconv.Itoa(10))
	q.Set("size", strconv.Itoa(10))
	q.Set("sort", "id_asc")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)

	repo := &infrastructure.MockProductRepository{}
	repo.
		On("GetAll", mock.Anything, mock.Anything).
		Return(infrastructure.Page[product.State]{}, nil)

	controller := ProductController{getAllHandler: application.ProvideGetAllProductsHandler(repo)}
	err := controller.getAll(ctx)

	assert.Nil(t, err)
}
