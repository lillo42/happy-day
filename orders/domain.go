package orders

import (
	"github.com/google/uuid"
	"time"
)

type (
	PaymentMethod string

	Order struct {
		ID         uuid.UUID `json:"id"`
		Address    string    `json:"address"`
		Comment    string    `json:"comment,omitempty"`
		DeliveryAt time.Time `json:"deliveryAt"`
		PickUp     time.Time `json:"pickUp"`
		TotalPrice float64   `json:"totalPrice"`
		Discount   float64   `json:"discount"`
		FinalPrice float64   `json:"finalPrice"`
		Customer   Customer  `json:"customer"`
		Products   []Product `json:"products,omitempty"`
		Payments   []Payment `json:"payments"`
		CreateAt   time.Time `json:"createAt"`
		UpdateAt   time.Time `json:"updateAt"`
		Version    uint      `json:"-"`
	}

	Customer struct {
		ID     uuid.UUID `json:"id"`
		Name   string    `json:"name"`
		Phones []string  `json:"phones"`
	}

	Product struct {
		ID       uuid.UUID `json:"id"`
		Name     string    `json:"name"`
		Quantity uint      `json:"quantity"`
		Price    float64   `json:"price"`
	}

	Payment struct {
		Method PaymentMethod `json:"method"`
		Info   string        `json:"paymentInfo"`
		Value  float64       `json:"value"`
		At     time.Time     `json:"at"`
	}
)

const (
	Pix          PaymentMethod = "pix"
	BankTransfer PaymentMethod = "bank-transfer"
	Cash         PaymentMethod = "cash"
)