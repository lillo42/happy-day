package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"happy_day/application"
	"happy_day/common"
	"happy_day/domain/customer"
	"happy_day/infrastructure"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCustomerControllerCreateWhenIsNotJson(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(common.RandString(10)))

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/customers")

	controller := CustomerController{}
	err := controller.create(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidBody, err)
}

func TestCustomerControllerCreateWhenErrOnHandler(t *testing.T) {
	b, _ := json.Marshal(application.CreateOrChangeCustomerRequest{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/customers")

	controller := CustomerController{createOrChangeHandler: application.CreateOrChangeCustomerHandler{}}
	err := controller.create(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, application.ErrCustomerNameIsEmpty, err)
}

func TestCustomerControllerCreate(t *testing.T) {
	b, _ := json.Marshal(application.CreateOrChangeCustomerRequest{
		State: customer.State{
			Name:    common.RandString(10),
			Comment: common.RandString(10),
			Phones: []customer.Phone{
				{
					Number: "123456789",
				},
			},
		},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/customers")

	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("Save", mock.Anything, mock.Anything).
		Return(customer.State{}, nil)

	controller := CustomerController{createOrChangeHandler: application.ProvideCreateOrChangeCustomerHandler(repo)}

	err := controller.create(ctx)
	assert.Nil(t, err)
}

func TestCustomerControllerUpdateWhenIsNotJson(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(common.RandString(10)))

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/customers/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	controller := CustomerController{}
	err := controller.update(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidBody, err)
}

func TestCustomerControllerUpdateWhenIdIsNotUuid(t *testing.T) {
	b, _ := json.Marshal(application.CreateOrChangeCustomerRequest{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/customers/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("123")

	controller := CustomerController{createOrChangeHandler: application.CreateOrChangeCustomerHandler{}}
	err := controller.update(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrCustomerNotFound, err)
}

func TestCustomerControllerUpdateWhenErrOnHandler(t *testing.T) {
	b, _ := json.Marshal(application.CreateOrChangeCustomerRequest{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/customers/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	repo := &infrastructure.MockCustomerRepository{}
	controller := CustomerController{createOrChangeHandler: application.ProvideCreateOrChangeCustomerHandler(repo)}
	err := controller.update(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, application.ErrCustomerNameIsEmpty, err)
}

func TestCustomerControllerUpdate(t *testing.T) {
	b, _ := json.Marshal(application.CreateOrChangeCustomerRequest{
		State: customer.State{
			Name:    common.RandString(10),
			Comment: common.RandString(10),
			Phones: []customer.Phone{
				{
					Number: "123456789",
				},
			},
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

	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(customer.State{
			Name:    common.RandString(10),
			Comment: common.RandString(10),
			Phones: []customer.Phone{
				{
					Number: "123456789",
				},
			},
		}, nil)

	repo.
		On("Save", mock.Anything, mock.Anything).
		Return(customer.State{
			Name:    common.RandString(10),
			Comment: common.RandString(10),
			Phones: []customer.Phone{
				{
					Number: "123456789",
				},
			},
		}, nil)

	controller := CustomerController{createOrChangeHandler: application.ProvideCreateOrChangeCustomerHandler(repo)}
	err := controller.update(ctx)

	assert.Nil(t, err)
}

func TestCustomerControllerDeleteWhenIdIsNotUuid(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/customers/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("123")

	controller := CustomerController{}
	err := controller.delete(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrCustomerNotFound, err)
}

func TestCustomerControllerDeleteWhenErrOnHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/customers/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	expectedErr := errors.New(common.RandString(10))
	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("Delete", mock.Anything, mock.Anything).
		Return(expectedErr)

	controller := CustomerController{deleteHandler: application.ProvideDeleteCustomerHandler(repo)}
	err := controller.delete(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCustomerControllerDelete(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/customers/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil)

	controller := CustomerController{deleteHandler: application.ProvideDeleteCustomerHandler(repo)}
	err := controller.delete(ctx)

	assert.Nil(t, err)
}

func TestCustomerControllerGetWhenIdIsNotUuid(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/customers/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("123")

	controller := CustomerController{}
	err := controller.get(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrCustomerNotFound, err)
}

func TestCustomerControllerGetWhenErrToGetById(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/customers/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	expectedErr := errors.New(common.RandString(10))
	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(customer.State{}, expectedErr)

	controller := CustomerController{getByIdHandler: application.ProvideGetCustomerByIdHandler(repo)}
	err := controller.get(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCustomerControllerGet(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/customers/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(customer.State{}, nil)

	controller := CustomerController{getByIdHandler: application.ProvideGetCustomerByIdHandler(repo)}
	err := controller.get(ctx)

	assert.Nil(t, err)
}

func TestCustomerControllerGetAllWhenErrOnHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/")

	expectedErr := errors.New(common.RandString(10))
	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("GetAll", mock.Anything, mock.Anything).
		Return(infrastructure.Page[customer.State]{}, expectedErr)

	controller := CustomerController{getAllHandler: application.ProvideGetAllCustomersHandler(repo)}
	err := controller.getAll(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestCustomerControllerGetAll(t *testing.T) {
	q := make(url.Values)
	q.Set("text", common.RandString(10))
	q.Set("page", strconv.Itoa(10))
	q.Set("size", strconv.Itoa(10))
	q.Set("sort", "id_asc")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)

	repo := &infrastructure.MockCustomerRepository{}
	repo.
		On("GetAll", mock.Anything, mock.Anything).
		Return(infrastructure.Page[customer.State]{}, nil)

	controller := CustomerController{getAllHandler: application.ProvideGetAllCustomersHandler(repo)}
	err := controller.getAll(ctx)

	assert.Nil(t, err)
}
