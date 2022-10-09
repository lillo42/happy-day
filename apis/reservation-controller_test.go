package apis

import (
	"bytes"
	"encoding/json"
	"happy_day/application"
	"happy_day/common"
	"happy_day/domain/product"
	"happy_day/infrastructure"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newEcho() *echo.Echo {
	e := echo.New()

	return e
}

func TestReservationQuoteWhenBodyIsNotJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservations/quote", strings.NewReader(common.RandString(10)))
	rec := httptest.NewRecorder()

	ec := echo.New()
	controller := ReservationController{}
	ctx := ec.NewContext(req, rec)
	err := controller.quoteReservation(ctx)
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

	json, _ := json.Marshal(quote)

	repository := &infrastructure.MockProductRepository{}
	repository.
		On("GetComposed", mock.Anything, mock.Anything).
		Return(([]product.State)(nil), infrastructure.ErrOneProductNotFound)

	controller := ReservationController{
		quoteHandler: application.ProvideQuoteReservationHandler(repository),
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservations/quote", bytes.NewBuffer(json))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ec := echo.New()

	ctx := ec.NewContext(req, rec)
	err := controller.quoteReservation(ctx)
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

	controller := ReservationController{
		quoteHandler: application.ProvideQuoteReservationHandler(repository),
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservations/quote", bytes.NewBuffer(j))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ec := echo.New()

	ctx := ec.NewContext(req, rec)
	err := controller.quoteReservation(ctx)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var res application.QuoteReservationProductResponse
	err = json.Unmarshal(rec.Body.Bytes(), &res)
	assert.Nil(t, err)

	assert.Equal(t, application.QuoteReservationProductResponse{Price: 25}, res)

}
