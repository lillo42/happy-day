package reservation

import (
	"time"

	"github.com/google/uuid"
)

const (
	Pix          PaymentMethod = "pix"
	BankTransfer PaymentMethod = "bankTransfer"
	Money        PaymentMethod = "money"
)

type (
	State struct {
		Id                  uuid.UUID            `bson:"id" json:"id"`
		Price               float64              `bson:"price" json:"price"`
		Discount            float64              `bson:"discount" json:"discount"`
		Products            []Product            `bson:"products" json:"products"`
		Delivery            DeliveryOrPickUp     `bson:"delivery" json:"delivery"`
		PickUp              DeliveryOrPickUp     `bson:"pickUp" json:"pickUp"`
		PaymentInstallments []PaymentInstallment `bson:"paymentInstallments" json:"paymentInstallments"`
		Comment             string               `json:"comment"`
		CreateAt            time.Time            `bson:"createAt" json:"createAt"`
		ModifyAt            time.Time            `bson:"modifyAt" json:"modifyAt"`
	}

	Product struct {
		Id     uuid.UUID `bson:"id" json:"id"`
		Price  float64   `bson:"price" json:"price"`
		Amount int64     `bson:"amount" json:"amount"`
	}

	DeliveryOrPickUp struct {
		At time.Time `bson:"at" json:"at"`
		By []string  `bson:"by" json:"by"`
	}

	PaymentMethod      string
	PaymentInstallment struct {
		Amount float64       `bson:"amount" json:"amount"`
		Method PaymentMethod `bson:"method" json:"method"`
		At     time.Time     `bson:"at" json:"at"`
	}
)
