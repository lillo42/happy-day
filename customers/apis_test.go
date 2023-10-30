package customers

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
	"happyday/infra"
	"log"
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

	_ = infra.GormFactory(context.Background()).AutoMigrate(&infra.Customer{})

	gin.SetMode(gin.TestMode)
	engine = gin.Default()
	Map(engine.Group("/api"))
}

func TestDeleteShouldResponse204WhenCustomerIdIsNotUuid(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/customers/1", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNoContent, res.Code)
}

func TestDeleteShouldResponse204WhenCustomerNotExits(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/customers/"+uuid.NewString(), nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNoContent, res.Code)
}

func TestDeleteShouldResponse204WhenCustomerExists(t *testing.T) {
	customerID := uuid.New()
	db := infra.GormFactory(context.Background()).Session(&gorm.Session{})
	result := db.Save(&infra.Customer{
		ExternalID: customerID,
		Name:       "test",
		Pix:        "teste@test.com",
		Comment:    "",
		Version:    1,
		CreateAt:   time.Now(),
		UpdateAt:   time.Now(),
	})

	assert.Nil(t, result.Error)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/customers/"+customerID.String(), nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNoContent, res.Code)

	result = db.First(&infra.Customer{}, "external_id = ?", customerID)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestPostShouldReturn400WhenBodyNotExists(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/customers", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestPostShouldReturn422WhenDataIsInvalid(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/customers", strings.NewReader("{}"))
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
}

func TestPostShouldReturn201WhenCreated(t *testing.T) {
	data, _ := json.Marshal(&CreateOrChangeCustomer{
		Name: "Lorem Ipsum",
		Phones: []string{
			"123456789",
		},
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/customers", bytes.NewBuffer(data))
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusCreated, res.Code)
	assert.Greater(t, res.Body.Len(), 0)
}

func TestPutShouldReturn404WhenIdIsNotUuid(t *testing.T) {

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/customers/1", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestPutShouldReturn400WhenBodyIsEmpty(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/customers/"+uuid.NewString(), nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestPutShouldReturn422WhenBodyIsInvalid(t *testing.T) {
	customerID := uuid.New()
	db := infra.GormFactory(context.Background()).Session(&gorm.Session{})
	result := db.Save(&infra.Customer{
		ExternalID: customerID,
		Name:       "Lorem Ipsum",
		Version:    1,
		CreateAt:   time.Now(),
		UpdateAt:   time.Now(),
	})

	assert.Nil(t, result.Error)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/customers/"+customerID.String(), strings.NewReader("{}"))
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
}

func TestPutShouldReturn200WhenBodyIsValid(t *testing.T) {
	customerID := uuid.New()
	db := infra.GormFactory(context.Background()).Session(&gorm.Session{})
	result := db.Save(&infra.Customer{
		ExternalID: customerID,
		Name:       "Lorem Ipsum",
		Version:    1,
		CreateAt:   time.Now(),
		UpdateAt:   time.Now(),
	})

	assert.Nil(t, result.Error)

	data, _ := json.Marshal(&CreateOrChangeCustomer{
		Name:    "Lorem",
		Comment: "Ipsum",
		Pix:     "pix-key",
		Phones: []string{
			"123456789",
		},
	})
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/customers/"+customerID.String(), bytes.NewBuffer(data))
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetShouldReturn404WhenIdIsNotUuid(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/customers/1", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestGetShouldReturn200WhenCustomerExists(t *testing.T) {
	customerID := uuid.New()
	db := infra.GormFactory(context.Background()).Session(&gorm.Session{})
	result := db.Save(&infra.Customer{
		ExternalID: customerID,
		Name:       "Lorem Ipsum",
		Version:    1,
		CreateAt:   time.Now(),
		UpdateAt:   time.Now(),
	})

	assert.Nil(t, result.Error)

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/customers/"+customerID.String(), nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestGetAll(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/customers?page=1&size=0&name=lorem&phones=123", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}
