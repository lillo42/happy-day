package infrastructure

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Customer struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:255;not null;index"`
	Comment     string
	Phones      []CustomerPhone
	Reservation []Reservation
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null;index"`
}

type CustomerPhone struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	CustomerID uint      `gorm:"not null;index"`
	Number     string    `gorm:"size:20"`
	Customer   *Customer `gorm:"foreignKey:CustomerID"`
}

type Product struct {
	gorm.Model
	Name     string    `gorm:"size:255;not null;index"`
	Price    float64   `gorm:"precision:19;scale:2;not null"`
	Products []Product `gorm:"many2many:products_collection"`
}

type Reservation struct {
	gorm.Model
	Price              float64 `gorm:"precision:19;scale:2;not null"`
	Discount           float64 `gorm:"precision:19;scale:2;not null"`
	FinalPrice         float64 `gorm:"precision:19;scale:2;not null"`
	Comment            *string
	CustomerID         uint                `gorm:"not null"`
	Customer           *Customer           `gorm:"foreignKey:CustomerID"`
	AddressID          uint                `gorm:"not null"`
	Address            *ReservationAddress `gorm:"foreignKey:AddressID"`
	Product            []ReservationProduct
	PaymentInstallment []ReservationPaymentInstallment
	DeliveryInfo       []ReservationDeliveryInformation
}

type ReservationProduct struct {
	ReservationID uint         `gorm:"not null;primaryKey;autoIncrement:false"`
	ProductID     uint         `gorm:"not null;primaryKey;autoIncrement:false"`
	Price         float64      `gorm:"precision:19;scale:2;not null"`
	Quantity      uint         `gorm:"not null"`
	Reservation   *Reservation `gorm:"foreignKey:ReservationID"`
	Product       *Product     `gorm:"foreignKey:ProductID"`
}

type ReservationAddress struct {
	gorm.Model
	Street       string `gorm:"size:255;not null"`
	Number       string `gorm:"size:255;not null"`
	Neighborhood string `gorm:"size:255;not null"`
	Complement   string `gorm:"size:255;not null"`
	PostalCodee  string `gorm:"size:255;not null"`
	City         string `gorm:"size:255;not null"`
}

type ReservationPaymentInstallment struct {
	ID            uint         `gorm:"primaryKey"`
	ReservationID uint         `gorm:"not null;index"`
	Amount        float64      `gorm:"precision:19;scale:2;not null"`
	Method        int          `gorm:"not null"`
	Reservation   *Reservation `gorm:"foreignKey:ReservationID"`
	CreatedAt     time.Time    `gorm:"not null"`
}

type ReservationDeliveryInformation struct {
	ID            uint         `gorm:"primaryKey"`
	ReservationID uint         `gorm:"not null;index"`
	Type          int          `gorm:"not null"`
	At            time.Time    `gorm:"not null"`
	By            string       `gorm:"not null"`
	Reservation   *Reservation `gorm:"foreignKey:ReservationID"`
}
