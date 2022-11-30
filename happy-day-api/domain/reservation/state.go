package reservation

import (
	"happy_day/domain/customer"
	"time"

	"github.com/google/uuid"
)

const (
	Pix          PaymentMethod = "pix"
	BankTransfer PaymentMethod = "bankTransfer"
	Cash         PaymentMethod = "cash"
)

type (
	State struct {
		Id                  uuid.UUID            `bson:"id" json:"id"`
		Price               float64              `bson:"price" json:"price"`
		Discount            float64              `bson:"discount" json:"discount"`
		FinalPrice          float64              `bson:"finalPrice" json:"finalPrice"`
		Products            []Product            `bson:"products" json:"products"`
		Delivery            DeliveryOrPickUp     `bson:"delivery" json:"delivery"`
		PickUp              DeliveryOrPickUp     `bson:"pickUp" json:"pickUp"`
		PaymentInstallments []PaymentInstallment `bson:"paymentInstallments" json:"paymentInstallments"`
		Comment             string               `bson:"comment" json:"comment,omitempty"`
		Customer            Customer             `bson:"customer" json:"customer"`
		Address             Address              `bson:"address" json:"address"`
		CreatedAt           time.Time            `bson:"createdAt" json:"createdAt"`
		ModifiedAt          time.Time            `bson:"modifiedAt" json:"modifiedAt"`
	}

	Product struct {
		Id       uuid.UUID `bson:"id" json:"id"`
		Price    float64   `bson:"price" json:"price"`
		Quantity int64     `bson:"quantity" json:"quantity"`
	}

	DeliveryOrPickUp struct {
		At time.Time `bson:"at" json:"at"`
		By []string  `bson:"by" json:"by"`
	}

	Customer struct {
		customer.State
	}

	Address struct {
		Street       string `bson:"street" json:"street"`
		Number       string `bson:"number" json:"number"`
		Neighborhood string `bson:"neighborhood" json:"neighborhood"`
		Complement   string `bson:"complement" json:"complement,omitempty"`
		PostalCode   string `bson:"postalCode" json:"postalCode"`
		City         string `bson:"city" json:"city"`
	}

	PaymentMethod      string
	PaymentInstallment struct {
		Amount float64       `bson:"amount" json:"amount"`
		Method PaymentMethod `bson:"method" json:"method"`
		At     time.Time     `bson:"at" json:"at"`
	}
)
