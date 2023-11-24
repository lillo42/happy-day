package orders

import (
	"context"
	"github.com/google/uuid"
	"math"
	"time"
)

type (
	Command struct {
		repository      OrderRepository
		productService  ProductService
		customerService CustomerService
		discountService DiscountService
	}

	CreateOrChange struct {
		ID         uuid.UUID `json:"id"`
		Address    string    `json:"address"`
		Comment    string    `json:"comment,omitempty"`
		DeliveryAt time.Time `json:"deliveryAt"`
		PickUp     time.Time `json:"pickUpAt"`
		TotalPrice float64   `json:"totalPrice"`
		Discount   float64   `json:"discount"`
		FinalPrice float64   `json:"finalPrice"`
		CustomerID uuid.UUID `json:"customerId"`
		Products   []Product `json:"products,omitempty"`
		Payments   []Payment `json:"payments"`
	}

	QuoteRequest struct {
		Products []Product `json:"products"`
	}

	QuoteResponse struct {
		TotalPrice float64 `json:"totalPrice"`
	}
)

func (c *Command) CreateOrChange(context context.Context, req CreateOrChange) (Order, error) {
	order, err := c.repository.GetOrCreate(context, req.ID)
	if err != nil {
		return Order{}, err
	}

	if req.ID != uuid.Nil && order.Version == 0 {
		return Order{}, ErrOrderNotFound
	}

	if len(req.Address) == 0 {
		return Order{}, ErrAddressIsEmpty
	}

	if len(req.Address) > 1000 {
		return Order{}, ErrAddressIsTooLarge
	}

	if !req.DeliveryAt.Before(req.PickUp) && !req.DeliveryAt.Equal(req.PickUp) {
		return Order{}, ErrDeliveryAtIsInvalid
	}

	if req.TotalPrice < 0 {
		return Order{}, ErrTotalPriceIsInvalid
	}

	if req.Discount < 0 {
		return Order{}, ErrDiscountIsInvalid
	}

	if req.FinalPrice < 0 {
		return Order{}, ErrFinalPriceIsInvalid
	}

	exists, err := c.customerService.Exists(context, req.CustomerID)
	if err != nil {
		return Order{}, err
	}

	if !exists {
		return Order{}, ErrCustomerNotFound
	}

	order.Products = make([]Product, len(req.Products))
	for i, prod := range req.Products {
		prodDB, err := c.productService.Get(context, prod.ID)
		if err != nil {
			return Order{}, err
		}

		if prodDB.ID != prod.ID {
			return Order{}, ErrProductNotFound
		}

		order.Products[i] = Product{
			ID:       prod.ID,
			Quantity: prod.Quantity,
			Price:    prodDB.Price,
		}
	}

	order.Payments = make([]Payment, len(req.Payments))
	for _, payment := range req.Payments {
		if payment.Value <= 0 {
			return Order{}, ErrPaymentValueIsInvalid
		}

		if len(payment.Info) == 0 {
			return Order{}, ErrPaymentInfoIsEmpty
		}

		if len(payment.Info) > 100 {
			return Order{}, ErrPaymentInfoIsTooLarge
		}
	}

	order.Address = req.Address
	order.Comment = req.Comment
	order.DeliveryAt = req.DeliveryAt
	order.PickUpAt = req.PickUp
	order.TotalPrice = req.TotalPrice
	order.Discount = req.Discount
	order.FinalPrice = req.FinalPrice
	order.Customer = Customer{ID: req.CustomerID}

	return c.repository.Save(context, order)
}

func (c *Command) Delete(ctx context.Context, id uuid.UUID) error {
	return c.repository.Delete(ctx, id)
}

func (c *Command) Quote(ctx context.Context, quote QuoteRequest) (QuoteResponse, error) {
	products := map[uuid.UUID]*Product{}
	productsID := make([]uuid.UUID, len(quote.Products))

	for i, p := range quote.Products {
		productsID[i] = p.ID
		prod, err := c.productService.Get(ctx, p.ID)
		if err != nil {
			return QuoteResponse{}, err
		}

		products[p.ID] = &Product{
			ID:       p.ID,
			Quantity: p.Quantity,
			Price:    prod.Price,
		}
	}

	discounts, err := c.discountService.GetAll(ctx, productsID)
	if err != nil {
		return QuoteResponse{}, err
	}

	var totalPrice float64
	for _, discount := range discounts {
		totalDiscount := math.MaxFloat64
		for _, prod := range discount.Products {
			amount, exists := products[prod.ID]
			if !exists {
				totalDiscount = 0
				break
			}

			totalDiscount = math.Min(totalDiscount, math.Floor(float64(amount.Quantity/prod.Quantity)))
		}

		if totalDiscount == 0 || totalDiscount == math.MaxFloat64 {
			continue
		}

		totalPrice += totalDiscount * discount.Price
		for _, prod := range discount.Products {
			products[prod.ID].Quantity -= prod.Quantity * uint(totalDiscount)
		}
	}

	for _, prod := range products {
		totalPrice += prod.Price * float64(prod.Quantity)
	}

	return QuoteResponse{TotalPrice: totalPrice}, nil
}
