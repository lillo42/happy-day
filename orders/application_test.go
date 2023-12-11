package orders

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"happyday/infra"
	"testing"
	"time"
)

func TestCreateOrChangeOrderNotFound(t *testing.T) {
	repo := new(MockOrderRepository)
	repo.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Order{}, nil)

	command := &Command{repository: repo}

	_, err := command.CreateOrChange(context.Background(), CreateOrChange{
		ID: uuid.New(),
	})

	assert.NotNil(t, err)
	assert.Equal(t, ErrOrderNotFound, err)
}

func TestCreateOrChangeInvalidAddress(t *testing.T) {
	repo := new(MockOrderRepository)
	repo.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Order{}, nil)

	command := &Command{repository: repo}

	t.Run("empty address", func(t *testing.T) {
		_, err := command.CreateOrChange(context.Background(), CreateOrChange{})
		assert.NotNil(t, err)
		assert.Equal(t, ErrAddressIsEmpty, err)
	})

	t.Run("address is too large", func(t *testing.T) {
		_, err := command.CreateOrChange(context.Background(), CreateOrChange{
			Address: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Eget egestas purus viverra accumsan. Semper auctor neque vitae tempus quam pellentesque nec. Tellus in hac habitasse platea dictumst vestibulum rhoncus est. Placerat duis ultricies lacus sed turpis tincidunt id aliquet risus. Augue neque gravida in fermentum et sollicitudin. Ornare aenean euismod elementum nisi quis eleifend quam adipiscing. Orci porta non pulvinar neque. Malesuada fames ac turpis egestas integer. Eget mauris pharetra et ultrices neque. Ut tortor pretium viverra suspendisse potenti. Eget mauris pharetra et ultrices neque. Ipsum suspendisse ultrices gravida dictum fusce ut placerat. Ultrices neque ornare aenean euismod elementum nisi quis. At risus viverra adipiscing at in tellus integer. Eu non diam phasellus vestibulum lorem sed. Diam volutpat commodo sed egestas egestas fringilla. Vitae aliquet nec ullamcorper sit amet risus nullam eget felis. Si",
		})

		assert.NotNil(t, err)
		assert.Equal(t, ErrAddressIsTooLarge, err)
	})
}

func TestCreateOrChangeInvalidDeliveryAt(t *testing.T) {
	repo := new(MockOrderRepository)
	repo.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Order{}, nil)

	command := &Command{repository: repo}

	_, err := command.CreateOrChange(context.Background(), CreateOrChange{
		Address:    "test",
		DeliveryAt: time.Now(),
	})

	assert.NotNil(t, err)
	assert.Equal(t, ErrDeliveryAtIsInvalid, err)

}

func TestCreateOrChangeInvalidPrice(t *testing.T) {
	repo := new(MockOrderRepository)
	repo.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Order{}, nil)

	command := &Command{repository: repo}

	t.Run("total price is invalid", func(t *testing.T) {
		_, err := command.CreateOrChange(context.Background(), CreateOrChange{
			Address:    "test",
			DeliveryAt: time.Now(),
			PickUp:     time.Now().Add(time.Hour * 24),
			TotalPrice: -1,
		})
		assert.NotNil(t, err)
		assert.Equal(t, ErrTotalPriceIsInvalid, err)
	})

	t.Run("discount is invalid", func(t *testing.T) {
		_, err := command.CreateOrChange(context.Background(), CreateOrChange{
			Address:    "test",
			DeliveryAt: time.Now(),
			PickUp:     time.Now().Add(time.Hour * 24),
			TotalPrice: 10,
			Discount:   -1,
		})
		assert.NotNil(t, err)
		assert.Equal(t, ErrDiscountIsInvalid, err)
	})

	t.Run("final price is invalid", func(t *testing.T) {
		_, err := command.CreateOrChange(context.Background(), CreateOrChange{
			Address:    "test",
			DeliveryAt: time.Now(),
			PickUp:     time.Now().Add(time.Hour * 24),
			TotalPrice: 10,
			Discount:   5,
			FinalPrice: -1,
		})
		assert.NotNil(t, err)
		assert.Equal(t, ErrFinalPriceIsInvalid, err)
	})
}

func TestCreateOrChangeInvalidCustomer(t *testing.T) {
	repo := new(MockOrderRepository)
	repo.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Order{}, nil)

	customerService := new(MockCustomerService)
	customerService.
		On("Exists", mock.Anything, mock.Anything).
		Return(false, nil)

	command := &Command{
		repository:      repo,
		customerService: customerService,
	}

	_, err := command.CreateOrChange(context.Background(), CreateOrChange{
		Address:    "test",
		DeliveryAt: time.Now(),
		PickUp:     time.Now().Add(time.Hour * 24),
		TotalPrice: 10,
		Discount:   5,
		FinalPrice: 5,
		CustomerID: uuid.New(),
	})

	assert.NotNil(t, err)
	assert.Equal(t, ErrCustomerNotFound, err)
}

func TestCreateOrChangeInvalidProduct(t *testing.T) {
	repo := new(MockOrderRepository)
	repo.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Order{}, nil)

	customerService := new(MockCustomerService)
	customerService.
		On("Exists", mock.Anything, mock.Anything).
		Return(true, nil)

	productService := new(MockProductService)
	productService.
		On("Get", mock.Anything, mock.Anything).
		Return(ProductProjection{}, nil)

	command := &Command{
		repository:      repo,
		customerService: customerService,
		productService:  productService,
	}

	_, err := command.CreateOrChange(context.Background(), CreateOrChange{
		Address:    "test",
		DeliveryAt: time.Now(),
		PickUp:     time.Now().Add(time.Hour * 24),
		TotalPrice: 10,
		Discount:   5,
		FinalPrice: 5,
		CustomerID: uuid.New(),
		Products: []Product{
			{ID: uuid.New(), Quantity: 10},
		},
	})

	assert.NotNil(t, err)
	assert.Equal(t, ErrProductNotFound, err)
}

func TestCreateOrChangeInvalidPayment(t *testing.T) {
	repo := new(MockOrderRepository)
	repo.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Order{}, nil)

	customerService := new(MockCustomerService)
	customerService.
		On("Exists", mock.Anything, mock.Anything).
		Return(true, nil)

	productID := uuid.New()
	productService := new(MockProductService)
	productService.
		On("Get", mock.Anything, productID).
		Return(Product{ID: productID, Price: 10}, nil)

	command := &Command{
		repository:      repo,
		customerService: customerService,
		productService:  productService,
	}

	t.Run("invalid payment value", func(t *testing.T) {
		_, err := command.CreateOrChange(context.Background(), CreateOrChange{
			Address:    "test",
			DeliveryAt: time.Now(),
			PickUp:     time.Now().Add(time.Hour * 24),
			TotalPrice: 10,
			Discount:   5,
			FinalPrice: 5,
			CustomerID: uuid.New(),
			Products: []Product{
				{ID: productID, Quantity: 10},
			},
			Payments: []Payment{
				{Method: Pix, At: time.Now(), Amount: 0},
			},
		})

		assert.NotNil(t, err)
		assert.Equal(t, ErrPaymentValueIsInvalid, err)
	})
}

func TestCreateOrChange(t *testing.T) {
	repo := new(MockOrderRepository)
	repo.
		On("GetOrCreate", mock.Anything, mock.Anything).
		Return(Order{}, nil)

	repo.
		On("Save", mock.Anything, mock.Anything).
		Return(Order{}, nil)

	customerService := new(MockCustomerService)
	customerService.
		On("Exists", mock.Anything, mock.Anything).
		Return(true, nil)

	productID := uuid.New()
	productService := new(MockProductService)
	productService.
		On("Get", mock.Anything, productID).
		Return(Product{ID: productID, Price: 10}, nil)

	command := &Command{
		repository:      repo,
		customerService: customerService,
		productService:  productService,
	}

	_, err := command.CreateOrChange(context.Background(), CreateOrChange{
		Address:    "test",
		DeliveryAt: time.Now(),
		PickUp:     time.Now().Add(time.Hour * 24),
		TotalPrice: 10,
		Discount:   5,
		FinalPrice: 5,
		CustomerID: uuid.New(),
		Products: []Product{
			{ID: productID, Quantity: 10},
		},
		Payments: []Payment{
			{Method: Pix, At: time.Now(), Amount: 10},
		},
	})

	assert.Nil(t, err)
}

func TestQuote(t *testing.T) {
	chairID := uuid.New()
	tableID := uuid.New()

	productService := new(MockProductService)
	productService.
		On("Get", mock.Anything, chairID).
		Return(ProductProjection{
			ID:    chairID,
			Name:  "chair",
			Price: 2,
		}, nil)

	productService.
		On("Get", mock.Anything, tableID).
		Return(ProductProjection{
			ID:    tableID,
			Name:  "table",
			Price: 4,
		}, nil)

	discountService := new(MockDiscountService)
	discountService.
		On("GetAll", mock.Anything, mock.Anything).
		Return([]DiscountProjection{
			{
				Price: 10,
				Products: []DiscountProducts{
					{ID: chairID, Quantity: 4},
					{ID: tableID, Quantity: 1},
				},
			},
		}, nil)

	command := &Command{
		productService:  productService,
		discountService: discountService,
	}

	t.Run("No discount match", func(t *testing.T) {
		res, _ := command.Quote(context.Background(), QuoteRequest{
			Products: []Product{
				{ID: chairID, Quantity: 2},
				{ID: tableID, Quantity: 1},
			},
		})

		assert.Equal(t, float64(8), res.TotalPrice)
	})

	t.Run("discount match", func(t *testing.T) {
		res, _ := command.Quote(context.Background(), QuoteRequest{
			Products: []Product{
				{ID: chairID, Quantity: 4},
				{ID: tableID, Quantity: 1},
			},
		})

		assert.Equal(t, float64(10), res.TotalPrice)
	})

	t.Run("discount match", func(t *testing.T) {
		res, _ := command.Quote(context.Background(), QuoteRequest{
			Products: []Product{
				{ID: chairID, Quantity: 8},
				{ID: tableID, Quantity: 1},
			},
		})

		assert.Equal(t, float64(18), res.TotalPrice)
	})

	t.Run("discount match", func(t *testing.T) {
		res, _ := command.Quote(context.Background(), QuoteRequest{
			Products: []Product{
				{ID: chairID, Quantity: 8},
				{ID: tableID, Quantity: 2},
			},
		})

		assert.Equal(t, float64(20), res.TotalPrice)
	})
}

var (
	_ OrderRepository = (*MockOrderRepository)(nil)
	_ CustomerService = (*MockCustomerService)(nil)
	_ ProductService  = (*MockProductService)(nil)
	_ DiscountService = (*MockDiscountService)(nil)
)

type (
	MockOrderRepository struct {
		mock.Mock
	}

	MockCustomerService struct {
		mock.Mock
	}

	MockProductService struct {
		mock.Mock
	}

	MockDiscountService struct {
		mock.Mock
	}
)

func (m *MockDiscountService) GetAll(ctx context.Context, productsID []uuid.UUID) ([]DiscountProjection, error) {
	args := m.Called(ctx, productsID)
	return args.Get(0).([]DiscountProjection), args.Error(1)
}

func (m *MockOrderRepository) GetAll(ctx context.Context, filter OrderFilter) (infra.Page[Order], error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(infra.Page[Order]), args.Error(1)
}

func (m *MockOrderRepository) GetOrCreate(ctx context.Context, id uuid.UUID) (Order, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Order), args.Error(1)
}

func (m *MockOrderRepository) Save(ctx context.Context, order Order) (Order, error) {
	args := m.Called(ctx, order)
	return args.Get(0).(Order), args.Error(1)
}

func (m *MockOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCustomerService) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockProductService) Get(ctx context.Context, id uuid.UUID) (ProductProjection, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(ProductProjection), args.Error(1)
}
