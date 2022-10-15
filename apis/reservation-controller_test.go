package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"happy_day/application"
	"happy_day/common"
	"happy_day/domain/customer"
	"happy_day/domain/product"
	"happy_day/domain/reservation"
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

func TestReservationCreateWhenBodyIsNotJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(common.RandString(10)))

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservations")

	controller := ReservationController{}
	err := controller.create(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidBody, err)
}

func TestReservationCreateWhenErrToHandler(t *testing.T) {
	b, _ := json.Marshal(application.CreateReservationRequest{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservations")

	controller := ReservationController{}
	err := controller.create(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, application.ErrProductListIsEmpty, err)
}

func TestReservationCreate(t *testing.T) {
	b, _ := json.Marshal(application.CreateReservationRequest{
		Products: []application.CreateReservationProductRequest{
			{
				Id:     uuid.New(),
				Amount: 10,
			},
		},
		Price:    100,
		Discount: 0,
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservations")

	productRepository := &infrastructure.MockProductRepository{}
	productRepository.
		On("GetByProducts", mock.Anything, mock.Anything).
		Return([]product.State{}, nil)

	reservationRepository := &infrastructure.MockReservationRepository{}
	reservationRepository.
		On("Save", mock.Anything, mock.Anything).
		Return(reservation.State{}, nil)

	controller := ReservationController{
		createHandler: application.ProvideCreateReservationHandler(productRepository, reservationRepository),
	}
	err := controller.create(ctx)

	assert.Nil(t, err)
}

func TestReservationQuoteWhenBodyIsNotJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(common.RandString(10)))

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservations/quote")

	controller := ReservationController{}
	err := controller.quote(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidBody, err)
}

func TestReservationQuoteWhenHandlerReturnErr(t *testing.T) {
	quote := application.QuoteReservationRequest{
		Products: []application.QuoteReservationProductRequest{
			{
				Id:     uuid.New(),
				Amount: 10,
			},
		},
	}

	b, _ := json.Marshal(quote)

	repository := &infrastructure.MockProductRepository{}
	repository.
		On("GetComposed", mock.Anything, mock.Anything).
		Return(([]product.State)(nil), infrastructure.ErrOneProductNotFound)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservations/quote", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservations/quote")

	controller := ReservationController{quoteHandler: application.ProvideQuoteReservationHandler(repository)}
	err := controller.quote(ctx)
	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrOneProductNotFound, err)
}

func TestReservationQuote(t *testing.T) {
	quote := application.QuoteReservationRequest{
		Products: []application.QuoteReservationProductRequest{
			{
				Id:     uuid.New(),
				Amount: 10,
			},
		},
	}

	j, _ := json.Marshal(quote)

	repository := &infrastructure.MockProductRepository{}
	repository.
		On("GetComposed", mock.Anything, mock.Anything).
		Return([]product.State{}, nil)

	repository.
		On("GetByProducts", mock.Anything, mock.Anything).
		Return([]product.State{
			{
				Id:    quote.Products[0].Id,
				Price: 2.5,
			},
		}, nil)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservations/quote", bytes.NewBuffer(j))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)

	controller := ReservationController{quoteHandler: application.ProvideQuoteReservationHandler(repository)}
	err := controller.quote(ctx)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var res application.QuoteReservationProductResponse
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.Nil(t, err)

	assert.Equal(t, application.QuoteReservationProductResponse{Price: 25}, res)

}

func TestReservationControllerUpdateWhenIsNotJson(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(common.RandString(10)))

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservation/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	controller := ReservationController{}
	err := controller.update(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, ErrInvalidBody, err)
}

func TestReservationControllerUpdateWhenIdIsNotUuid(t *testing.T) {
	b, _ := json.Marshal(application.CreateOrChangeProductRequest{})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservations/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("123")

	controller := ReservationController{changeHandler: application.ChangeReservationHandler{}}
	err := controller.update(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrReservationNotFound, err)
}

func TestReservationControllerUpdateWhenErrOnHandler(t *testing.T) {
	b, _ := json.Marshal(application.ChangeReservationRequest{
		PaymentInstallments: []reservation.PaymentInstallment{
			{
				Amount: -1,
			},
		},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservations/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	controller := ReservationController{changeHandler: application.ChangeReservationHandler{}}
	err := controller.update(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, application.ErrInvalidPaymentInstallmentAmount, err)
}

func TestReservationControllerUpdate(t *testing.T) {
	b, _ := json.Marshal(application.ChangeReservationRequest{
		Customer: reservation.Customer{
			State: customer.State{
				Name: "John Doe",
				Phones: []customer.Phone{
					{
						Number: "123456789",
					},
				},
			},
		},
		Address: reservation.Address{
			Street:       common.RandString(10),
			Number:       common.RandString(10),
			Neighborhood: common.RandString(10),
			PostalCode:   common.RandString(10),
			City:         common.RandString(10),
		},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewBuffer(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservation/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	repo := &infrastructure.MockReservationRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(reservation.State{
			Id:    uuid.New(),
			Price: 10,
		}, nil)

	repo.
		On("Save", mock.Anything, mock.Anything).
		Return(reservation.State{
			Id:    uuid.New(),
			Price: 3,
		}, nil)

	controller := ReservationController{changeHandler: application.ProvideChangeReservationHandler(repo)}
	err := controller.update(ctx)

	assert.Nil(t, err)
}

func TestReservationControllerDeleteWhenIdIsNotUuid(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservations/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("123")

	controller := ReservationController{}
	err := controller.delete(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrReservationNotFound, err)
}

func TestReservationControllerDeleteWhenErrOnHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservations/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	expectedErr := errors.New(common.RandString(10))
	repo := &infrastructure.MockReservationRepository{}
	repo.
		On("Delete", mock.Anything, mock.Anything).
		Return(expectedErr)

	controller := ReservationController{deleteHandler: application.ProvideDeleteReservationHandler(repo)}
	err := controller.delete(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestReservationControllerDelete(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservations/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	repo := &infrastructure.MockReservationRepository{}
	repo.
		On("Delete", mock.Anything, mock.Anything).
		Return(nil)

	controller := ReservationController{deleteHandler: application.ProvideDeleteReservationHandler(repo)}
	err := controller.delete(ctx)

	assert.Nil(t, err)
}

func TestReservationControllerGetWhenIdIsNotUuid(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservations/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("123")

	controller := ReservationController{}
	err := controller.get(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, infrastructure.ErrReservationNotFound, err)
}

func TestReservationControllerGetWhenErrToGetById(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservations/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	expectedErr := errors.New(common.RandString(10))
	repo := &infrastructure.MockReservationRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(reservation.State{}, expectedErr)

	controller := ReservationController{getByIdHandler: application.ProvideGetReservationByIdHandler(repo)}
	err := controller.get(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestReservationControllerGet(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/api/v1/reservations/:id")
	ctx.SetParamNames("id")
	ctx.SetParamValues("cf9fe14e-393d-4d5d-9800-3f4448a04e9c")

	repo := &infrastructure.MockReservationRepository{}
	repo.
		On("GetById", mock.Anything, mock.Anything).
		Return(reservation.State{}, nil)

	controller := ReservationController{getByIdHandler: application.ProvideGetReservationByIdHandler(repo)}
	err := controller.get(ctx)

	assert.Nil(t, err)
}

func TestReservationControllerGetAllWhenErrOnHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)
	ctx.SetPath("/")

	expectedErr := errors.New(common.RandString(10))
	repo := &infrastructure.MockReservationRepository{}
	repo.
		On("GetAll", mock.Anything, mock.Anything).
		Return(infrastructure.Page[reservation.State]{}, expectedErr)

	controller := ReservationController{getAllHandler: application.ProvideGetAllReservationHandler(repo)}
	err := controller.getAll(ctx)

	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err)
}

func TestReservationControllerGetAll(t *testing.T) {
	q := make(url.Values)
	q.Set("text", common.RandString(10))
	q.Set("page", strconv.Itoa(10))
	q.Set("size", strconv.Itoa(10))
	q.Set("sort", "id_asc")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)

	ec := echo.New()
	ctx := ec.NewContext(req, rec)

	repo := &infrastructure.MockReservationRepository{}
	repo.
		On("GetAll", mock.Anything, mock.Anything).
		Return(infrastructure.Page[reservation.State]{}, nil)

	controller := ReservationController{getAllHandler: application.ProvideGetAllReservationHandler(repo)}
	err := controller.getAll(ctx)

	assert.Nil(t, err)
}
