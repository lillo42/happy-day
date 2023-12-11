package orders

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"happyday/customers"
	"happyday/discounts"
	"happyday/infra"
	"happyday/products"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

var engine *gin.Engine

func init() {
	if infra.GormFactory == nil {
		infra.GormFactory = func(ctx context.Context) *gorm.DB {
			db, err := gorm.Open(sqlite.Open("../integration-test.db"), &gorm.Config{
				Logger: &infra.SlogGorm{
					Logger:                    slog.Default(),
					LogLevel:                  gormlogger.Info,
					IgnoreRecordNotFoundError: true,
				},
			})
			if err != nil {
				panic(err)
			}

			return db
		}
	}

	if ProductServiceFactory == nil {
		ProductServiceFactory = func(ctx context.Context) ProductService {
			return &TestProductService{
				repository: products.CreateRepository(ctx),
			}
		}
	}

	if CustomerServiceFactory == nil {
		CustomerServiceFactory = func(ctx context.Context) CustomerService {
			return &TestCustomerService{
				repository: customers.CreateRepository(ctx),
			}
		}
	}

	if DiscountServiceFactory == nil {
		DiscountServiceFactory = func(ctx context.Context) DiscountService {
			return &TestDiscountService{
				repository: discounts.CreateRepository(ctx),
			}
		}
	}

	err := infra.GormFactory(context.Background()).
		AutoMigrate(
			&infra.Customer{},
			&infra.Product{},
			&infra.Discount{}, &infra.DiscountProducts{},
			&infra.Order{}, &infra.OrderProduct{}, &infra.OrderPayment{})

	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.TestMode)
	engine = gin.Default()
	Map(engine.Group("/api"))
}

func TestHttpDelete(t *testing.T) {
	t.Run("should return 204 when id is not uuid", func(t *testing.T) {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/api/orders/1", nil)
		engine.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)
	})

	t.Run("should return 204 when order not exits", func(t *testing.T) {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/api/orders/"+uuid.NewString(), nil)
		engine.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)
	})

	t.Run("should return 204 when exits", func(t *testing.T) {
		testID := uuid.New()
		orderID := uuid.New()

		db := infra.GormFactory(context.Background())
		db.Save(&infra.Order{
			ExternalID: orderID,
			Address:    "Lorem address" + testID.String(),
			Customer: infra.Customer{
				ExternalID: uuid.New(),
				Name:       "Lorem customer" + testID.String(),
			},
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
			Version:  0,
		})

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, "/api/orders/"+orderID.String(), nil)
		engine.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNoContent, res.Code)

		result := db.First(&infra.Order{}, "external_id = ?", orderID)
		assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
	})
}

func TestHttpPost(t *testing.T) {
	t.Run("should return 400 when body is when", func(t *testing.T) {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/orders", nil)
		engine.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Code)
	})

	t.Run("should return 422 when body data is invalid", func(t *testing.T) {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/orders", strings.NewReader("{}"))
		engine.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
	})

	t.Run("should return 201 when body data is valid", func(t *testing.T) {
		testID := uuid.New()
		customerID := uuid.New()
		productID := uuid.New()

		db := infra.GormFactory(context.Background())
		db.Save(&infra.Customer{
			ExternalID: customerID,
			Name:       "Lorem order customer" + testID.String(),
			CreateAt:   time.Now(),
			UpdateAt:   time.Now(),
			Version:    1,
		})

		db.Save(&infra.Product{
			ExternalID: productID,
			Name:       "Lorem product" + testID.String(),
			Price:      2,
			CreateAt:   time.Now(),
			UpdateAt:   time.Now(),
			Version:    1,
		})

		data, _ := json.Marshal(&CreateOrChange{
			Address:    "Lorem address" + testID.String(),
			DeliveryAt: time.Now(),
			PickUp:     time.Now().Add(time.Hour * 2),
			TotalPrice: 10,
			Discount:   2,
			FinalPrice: 8,
			CustomerID: customerID,
			Payments: []Payment{
				{At: time.Now(), Amount: 8, Method: Pix},
			},
			Products: []Product{
				{ID: productID, Quantity: 5},
			},
		})

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/orders", bytes.NewReader(data))
		engine.ServeHTTP(res, req)

		assert.Equal(t, http.StatusCreated, res.Code)
	})
}

func TestHttpPut(t *testing.T) {
	t.Run("should return 404 when body is when", func(t *testing.T) {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/api/orders/1", nil)
		engine.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("should return 400 when body is when", func(t *testing.T) {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/api/orders"+uuid.NewString(), nil)
		engine.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("should return 422 when body data is invalid", func(t *testing.T) {
		testID := uuid.New()
		orderID := uuid.New()
		customerID := uuid.New()

		db := infra.GormFactory(context.Background())
		db.Save(&infra.Order{
			ExternalID: orderID,
			Address:    "Lorem address" + testID.String(),
			DeliveryAt: time.Now(),
			PickUp:     time.Now().Add(time.Hour * 2),
			TotalPrice: 10,
			Discount:   2,
			FinalPrice: 8,
			Customer: infra.Customer{
				ExternalID: customerID,
				Name:       "Lorem order customer" + testID.String(),
				CreateAt:   time.Now(),
				UpdateAt:   time.Now(),
				Version:    1,
			},
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
			Version:  1,
		})

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/api/orders/"+orderID.String(), strings.NewReader("{}"))
		engine.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
	})

	t.Run("should return 200 when body data is valid", func(t *testing.T) {
		testID := uuid.New()
		orderID := uuid.New()
		customerID := uuid.New()
		productID := uuid.New()

		db := infra.GormFactory(context.Background())
		db.Save(&infra.Product{
			ExternalID: productID,
			Name:       "Lorem product" + testID.String(),
			Price:      2,
			CreateAt:   time.Now(),
			UpdateAt:   time.Now(),
			Version:    1,
		})

		db.Save(&infra.Order{
			ExternalID: orderID,
			Address:    "Lorem address" + testID.String(),
			DeliveryAt: time.Now(),
			PickUp:     time.Now().Add(time.Hour * 2),
			TotalPrice: 10,
			Discount:   2,
			FinalPrice: 8,
			Customer: infra.Customer{
				ExternalID: customerID,
				Name:       "Lorem order customer" + testID.String(),
				CreateAt:   time.Now(),
				UpdateAt:   time.Now(),
				Version:    1,
			},
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
			Version:  1,
		})

		data, _ := json.Marshal(&CreateOrChange{
			Address:    "Lorem address" + testID.String(),
			DeliveryAt: time.Now(),
			PickUp:     time.Now().Add(time.Hour * 2),
			TotalPrice: 12,
			Discount:   2,
			FinalPrice: 10,
			CustomerID: customerID,
			Payments: []Payment{
				{At: time.Now(), Amount: 8, Method: Pix},
			},
			Products: []Product{
				{ID: productID, Quantity: 6},
			},
		})

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/api/orders/"+orderID.String(), bytes.NewReader(data))
		engine.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
}

func TestHttpGetByID(t *testing.T) {
	t.Run("should return 404 when id is not uuid", func(t *testing.T) {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/orders/1", nil)
		engine.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("should return 404 when order not exists", func(t *testing.T) {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/orders/"+uuid.NewString(), nil)
		engine.ServeHTTP(res, req)

		assert.Equal(t, http.StatusNotFound, res.Code)
	})

	t.Run("should return 200 when order exits", func(t *testing.T) {
		testID := uuid.New()
		orderID := uuid.New()

		db := infra.GormFactory(context.Background())
		db.Save(&infra.Order{
			ExternalID: orderID,
			Address:    "Lorem address" + testID.String(),
			Customer: infra.Customer{
				ExternalID: uuid.New(),
				Name:       "Lorem customer" + testID.String(),
				CreateAt:   time.Now(),
				UpdateAt:   time.Now(),
				Version:    1,
			},
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
			Version:  1,
		})

		res := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/orders/"+orderID.String(), nil)
		engine.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Code)
	})
}

func TestHttpGet(t *testing.T) {
	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/orders?page=1&size=51&address=l&comment=o&customerName=r&customerPhone=3", nil)
	engine.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

type (
	TestProductService struct {
		repository products.ProductRepository
	}

	TestCustomerService struct {
		repository customers.CustomerRepository
	}

	TestDiscountService struct {
		repository discounts.DiscountRepository
	}
)

func (t *TestDiscountService) GetAll(ctx context.Context, productsID []uuid.UUID) ([]DiscountProjection, error) {
	discounts, err := t.repository.GetAllWithProducts(ctx, productsID)
	if err != nil {
		return nil, err
	}

	proj := make([]DiscountProjection, len(discounts))
	for i, discount := range discounts {
		prods := make([]DiscountProducts, len(discount.Products))
		for j, prod := range discount.Products {
			prods[j] = DiscountProducts{
				ID:       prod.ID,
				Quantity: prod.Quantity,
			}
		}

		proj[i] = DiscountProjection{
			Price:    discount.Price,
			Products: prods,
		}
	}

	return proj, nil
}

func (c *TestCustomerService) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	cus, err := c.repository.GetOrCreate(ctx, id)
	if err != nil {
		return false, err
	}

	return cus.Version > 0, nil
}

func (s *TestProductService) Get(ctx context.Context, id uuid.UUID) (ProductProjection, error) {
	prod, err := s.repository.GetOrCreate(ctx, id)
	if err != nil {
		return ProductProjection{}, err
	}

	return ProductProjection{
		ID:    prod.ID,
		Name:  prod.Name,
		Price: prod.Price,
	}, err
}
