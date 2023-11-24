package orders

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"happyday/infra"
	"time"
)

type (
	OrderFilter struct {
		Address         string
		Comment         string
		CustomerName    string
		CustomerPhone   string
		DeliveryBetween []time.Time
		Page            int
		Size            int
	}

	OrderRepository interface {
		GetAll(ctx context.Context, filter OrderFilter) (infra.Page[Order], error)
		GetOrCreate(ctx context.Context, id uuid.UUID) (Order, error)

		Save(ctx context.Context, order Order) (Order, error)
		Delete(ctx context.Context, id uuid.UUID) error
	}

	GormOrderRepository struct {
		db *gorm.DB
	}
)

func (g *GormOrderRepository) GetAll(ctx context.Context, filter OrderFilter) (infra.Page[Order], error) {
	query := g.db.
		WithContext(ctx).
		Model(&infra.Order{}).
		Joins("JOIN customers ON customers.id = orders.customer_id")

	if len(filter.Address) > 0 {
		query.Where("orders.address LIKE ?", "%"+filter.Address+"%")
	}

	if len(filter.Comment) > 0 {
		query.Where("orders.comment LIKE ?", "%"+filter.Comment+"%")
	}

	if len(filter.DeliveryBetween) == 2 {
		query.Where("orders.delivery_at BETWEEN ? AND ?", filter.DeliveryBetween[0], filter.DeliveryBetween[1])
	}

	if len(filter.CustomerName) > 0 {
		query.Where("customers.name LIKE ?", "%"+filter.CustomerName+"%")
	}

	if len(filter.CustomerPhone) > 0 {
		query.Where("customers.phones LIKE ?", "%"+filter.CustomerPhone+"%")
	}

	var counter int64
	result := query.Count(&counter)
	if result.Error != nil {
		return infra.Page[Order]{}, result.Error
	}

	var ordersDB []infra.Order
	result = query.
		Preload("Customer").
		Preload("Payments").
		Preload("Products").
		Preload("Products.Product").
		Order("orders.id").
		Limit(filter.Size).
		Offset(filter.Page * filter.Size).
		Find(&ordersDB)

	if result.Error != nil {
		return infra.Page[Order]{}, result.Error
	}

	var totalPage int64
	if counter > 0 {
		totalPage = counter / int64(filter.Size)
	}

	page := infra.Page[Order]{
		Items:      make([]Order, len(ordersDB)),
		TotalItems: counter,
		TotalPage:  totalPage,
	}

	for i, orderDB := range ordersDB {
		page.Items[i] = mapToOrder(orderDB)
	}

	return page, nil
}

func (g *GormOrderRepository) GetOrCreate(ctx context.Context, id uuid.UUID) (Order, error) {
	var orderDB infra.Order
	result := g.db.
		WithContext(ctx).
		Joins("Customer").
		Preload("Payments").
		Preload("Products").
		Preload("Products.Product").
		First(&orderDB, "orders.external_id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Order{}, nil
		}

		return Order{}, result.Error
	}

	return mapToOrder(orderDB), nil
}

func (g *GormOrderRepository) Save(ctx context.Context, order Order) (Order, error) {
	var orderDB infra.Order
	if order.Version == 0 {
		orderDB.ExternalID = uuid.New()
		orderDB.CreateAt = time.Now()
	} else {
		result := g.db.
			WithContext(ctx).
			First(&orderDB, "external_id = ?", order.ID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return Order{}, ErrOrderNotFound
			}
			return Order{}, result.Error
		}
	}

	orderDB.Comment = order.Comment
	orderDB.Address = order.Address
	orderDB.DeliveryAt = order.DeliveryAt
	orderDB.PickUp = order.PickUpAt
	orderDB.TotalPrice = order.TotalPrice
	orderDB.Discount = order.Discount
	orderDB.FinalPrice = order.FinalPrice
	orderDB.Payments = make([]infra.OrderPayment, len(order.Payments))
	orderDB.Products = make([]infra.OrderProduct, len(order.Products))
	orderDB.Version = order.Version + 1
	orderDB.UpdateAt = time.Now()

	for i, payment := range order.Payments {
		orderDB.Payments[i] = infra.OrderPayment{
			ID:     uuid.New(),
			Method: string(payment.Method),
			At:     payment.At,
			Value:  payment.Amount,
			Info:   payment.Info,
		}
	}

	for i, prod := range order.Products {
		var prodDB infra.Product
		result := g.db.
			WithContext(ctx).
			First(&prodDB, "external_id = ?", prod.ID)

		if result.Error != nil {
			return Order{}, result.Error
		}

		orderDB.Products[i] = infra.OrderProduct{
			ProductID: prodDB.ID,
			Quantity:  prod.Quantity,
			Price:     prod.Price,
		}
	}

	var customerDB infra.Customer
	result := g.db.
		WithContext(ctx).
		First(&customerDB, "external_id = ?", order.Customer.ID)

	if result.Error != nil {
		return Order{}, result.Error
	}

	orderDB.CustomerID = customerDB.ID
	orderDB.Customer = customerDB

	err := g.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			if orderDB.ID > 0 {
				result := tx.Delete(&infra.OrderPayment{}, "order_id = ?", orderDB.ID)
				if result.Error != nil {
					return result.Error
				}

				result = tx.Delete(&infra.OrderProduct{}, "order_id = ?", orderDB.ID)
				if result.Error != nil {
					return result.Error
				}
			}

			result := tx.
				Where("version = ?", order.Version).
				Save(&orderDB)
			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return ErrConcurrencyUpdate
			}

			return nil
		})

	if err != nil {
		return Order{}, err
	}

	return mapToOrder(orderDB), nil
}

func (g *GormOrderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var orderDB infra.Order
	result := g.db.
		WithContext(ctx).
		First(&orderDB, "external_id = ?", id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		return result.Error
	}

	return g.db.
		WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			result := tx.Delete(&infra.OrderPayment{}, "order_id = ?", orderDB.ID)
			if result.Error != nil {
				return result.Error
			}

			result = tx.Delete(&infra.OrderProduct{}, "order_id = ?", orderDB.ID)
			if result.Error != nil {
				return result.Error
			}

			result = tx.Delete(&infra.Order{}, orderDB.ID)
			return result.Error
		})
}

func mapToOrder(orderDB infra.Order) Order {
	var phones []string
	if len(orderDB.Customer.Phones) > 0 {
		_ = json.Unmarshal([]byte(orderDB.Customer.Phones), &phones)
	}

	order := Order{
		ID:         orderDB.ExternalID,
		Address:    orderDB.Address,
		Comment:    orderDB.Comment,
		DeliveryAt: orderDB.DeliveryAt,
		PickUpAt:   orderDB.PickUp,
		TotalPrice: orderDB.TotalPrice,
		Discount:   orderDB.Discount,
		FinalPrice: orderDB.FinalPrice,
		CreateAt:   orderDB.CreateAt,
		UpdateAt:   orderDB.UpdateAt,
		Version:    orderDB.Version,
		Payments:   make([]Payment, len(orderDB.Payments)),
		Products:   make([]Product, len(orderDB.Products)),
		Customer: Customer{
			ID:      orderDB.Customer.ExternalID,
			Name:    orderDB.Customer.Name,
			Comment: orderDB.Customer.Comment,
			Phones:  phones,
		},
	}

	for i, prod := range orderDB.Products {
		order.Products[i] = Product{
			ID:       prod.Product.ExternalID,
			Name:     prod.Product.Name,
			Price:    prod.Price,
			Quantity: prod.Quantity,
		}
	}

	for i, payment := range orderDB.Payments {
		order.Payments[i] = Payment{
			Amount: payment.Value,
			At:     payment.At,
			Info:   payment.Info,
			Method: PaymentMethod(payment.Method),
		}
	}

	return order
}
