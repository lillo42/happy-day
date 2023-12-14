package products

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"mec/infra"
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

	_ = infra.GormFactory(context.Background()).
		AutoMigrate(&infra.Product{})

	gin.SetMode(gin.TestMode)
	engine = gin.Default()
	Map(engine.Group("/api"))
}

func TestHttpDeleteShouldResponse204WhenCustomerIdIsNotUuid(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/products/1", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNoContent, res.Code)
}

func TestHttpDelete(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/products/"+uuid.NewString(), nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNoContent, res.Code)
}

func TestHttpPostShouldResponse400WhenBodyIsEmpty(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/products", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestHttpPostShouldResponse422WhenBodyIsInvalid(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/products", strings.NewReader("{}"))
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
}

func TestHttpPostShouldResponse201(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/products", strings.NewReader("{\"name\": \"test\", \"price\": 10}"))
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)
}

func TestHttpPutShouldResponse404WhenIdIsNotUuid(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/products/1", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestHttpPutShouldResponse400WhenBodyIsEmpty(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/products/"+uuid.NewString(), nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestHttpPutShouldResponse422WhenBodyIsInvalid(t *testing.T) {
	id := uuid.New()

	db := infra.GormFactory(context.Background()).Session(&gorm.Session{})
	db.Save(&infra.Product{
		ExternalID: id,
		Name:       "Lorem Ipsum",
		Price:      5,
		Version:    1,
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/products/"+id.String(), strings.NewReader("{}"))
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
}

func TestHttpPut(t *testing.T) {
	id := uuid.New()

	db := infra.GormFactory(context.Background()).Session(&gorm.Session{})
	db.Save(&infra.Product{
		ExternalID: id,
		Name:       "Lorem Ipsum",
		Price:      5,
		Version:    1,
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/products/"+id.String(), strings.NewReader("{\"name\": \"Lorem\", \"price\": 100}"))
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestHttpGetShouldResponse404WhenIdIsNotUuid(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/products/1", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestHttpGetShouldResponse404WhenProductNotExists(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/products/"+uuid.NewString(), nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestHttpGetShouldResponse200(t *testing.T) {
	id := uuid.New()

	db := infra.GormFactory(context.Background()).Session(&gorm.Session{})
	db.Save(&infra.Product{
		ExternalID: id,
		Name:       "Lorem Ipsum",
		Price:      5,
		CreateAt:   time.Now(),
		UpdateAt:   time.Now(),
		Version:    1,
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/products/"+id.String(), nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestHttpGetAllShouldResponse200(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/products?name=lor&page=1&size=100", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}
