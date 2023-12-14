package discounts

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"mec/infra"
	"mec/products"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

var engine *gin.Engine

func init() {
	if infra.GormFactory == nil {
		infra.GormFactory = func(ctx context.Context) *gorm.DB {
			db, err := gorm.Open(sqlite.Open("../integration-test.db"), &gorm.Config{
				Logger: logger.New(
					log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
					logger.Config{
						LogLevel:                  logger.Silent,
						IgnoreRecordNotFoundError: true,
						Colorful:                  false,
					}),
			})
			if err != nil {
				panic(err)
			}

			return db
		}
	}

	if ProductServiceFactory == nil {
		ProductServiceFactory = func(ctx context.Context) ProductService {
			return products.CreateCommand(ctx)
		}
	}

	err := infra.GormFactory(context.Background()).
		AutoMigrate(&infra.Product{},
			&infra.Discount{},
			&infra.DiscountProducts{})

	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.TestMode)
	engine = gin.Default()
	Map(engine.Group("/api"))
}

func TestHttpDeleteShouldResponse204WhenIdIsNotUuid(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/discounts/1", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNoContent, res.Code)
}

func TestHttpDeleteShouldResponse204WhenNotExists(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/discounts/"+uuid.NewString(), nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNoContent, res.Code)
}

func TestDeleteShouldResponse204WheExists(t *testing.T) {
	discountID := uuid.New()
	db := infra.GormFactory(context.Background())
	result := db.Save(&infra.Discount{
		ExternalID: discountID,
		Name:       "test",
		Price:      10,
		Version:    1,
		CreateAt:   time.Now(),
		UpdateAt:   time.Now(),
	})

	assert.Nil(t, result.Error)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/discounts/"+discountID.String(), nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNoContent, res.Code)

	result = db.First(&infra.Discount{}, "external_id = ?", discountID)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestHttpPostShouldResponse400WhenBodyIsEmpty(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/discounts", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestHttpPostShouldResponse422WhenBodyIsInvalid(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/discounts", strings.NewReader("{}"))
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
}

func TestHttpPostShouldResponse201(t *testing.T) {
	testID := uuid.New()
	productID := uuid.New()

	db := infra.GormFactory(context.Background())
	db.Save(&infra.Product{
		ExternalID: productID,
		Name:       "Lorem Prod" + testID.String(),
		Price:      10,
		Version:    1,
	})

	data, _ := json.Marshal(&CreateOrChange{
		Name:  "Lorem dis" + testID.String(),
		Price: 5,
		Products: []Product{
			{ID: productID, Quantity: 2},
		},
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/discounts", bytes.NewReader(data))
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)
}

func TestHttpPutShouldResponse404WhenIdIsNotUuid(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/discounts/1", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestHttpPutShouldResponse400WhenBodyIsEmpty(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/discounts/"+uuid.NewString(), nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestHttpPutShouldResponse422WhenBodyIsInvalid(t *testing.T) {
	testID := uuid.New()
	id := uuid.New()

	db := infra.GormFactory(context.Background())
	db.Save(&infra.Discount{
		ExternalID: id,
		Name:       "Lorem prod" + testID.String(),
		Price:      5,
		Version:    1,
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/discounts/"+id.String(), strings.NewReader("{}"))
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
}

func TestHttpPut(t *testing.T) {
	testID := uuid.New()
	id := uuid.New()
	productID := uuid.New()

	db := infra.GormFactory(context.Background())
	db.Save(&infra.Product{
		ExternalID: productID,
		Name:       "Lorem prod" + testID.String(),
		Price:      5,
		Version:    1,
	})

	db.Save(&infra.Discount{
		ExternalID: id,
		Name:       "Lorem discount" + testID.String(),
		Price:      5,
		Version:    1,
	})

	data, _ := json.Marshal(&CreateOrChange{
		Name:  "Lorem dis" + testID.String(),
		Price: 5,
		Products: []Product{
			{ID: productID, Quantity: 2},
		},
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/discounts/"+id.String(), bytes.NewReader(data))
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestHttpGetShouldResponse404WhenIdIsNotUuid(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/discounts/1", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestHttpGetShouldResponse404WhenNotExists(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/discounts/"+uuid.NewString(), nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestHttpGetShouldResponse200(t *testing.T) {
	testID := uuid.New()
	id := uuid.New()
	product := &infra.Product{
		ExternalID: uuid.New(),
		Name:       "Lorem prod" + testID.String(),
		Price:      5,
		CreateAt:   time.Now(),
		UpdateAt:   time.Now(),
		Version:    1,
	}

	db := infra.GormFactory(context.Background()).Session(&gorm.Session{})
	db.Save(product)

	db.Save(&infra.Discount{
		ExternalID: id,
		Name:       "Lorem discount" + testID.String(),
		Price:      10,
		Version:    1,
		CreateAt:   time.Now(),
		UpdateAt:   time.Now(),
		Products: []infra.DiscountProducts{
			{ProductID: product.ID, Quantity: 3},
		},
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/discounts/"+id.String(), nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestHttpGetAllShouldResponse200(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/discounts?name=lor&page=1&size=100", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}
