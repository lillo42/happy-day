package infra

import (
	"github.com/google/uuid"
	"time"
)

type (
	Customer struct {
		ID         uint      `gorm:"primaryKey"`
		ExternalID uuid.UUID `gorm:"size:36;uniqueIndex;not null"`
		Name       string    `gorm:"size:255;not null"`
		Comment    string
		Phones     string
		Pix        string    `gorm:"size:50"`
		Orders     []Order   `gorm:"foreignKey:CustomerID;references:ID"`
		CreateAt   time.Time `gorm:"not null"`
		UpdateAt   time.Time `gorm:"not null"`
		Version    uint      `gorm:"not null"`
	}

	Product struct {
		ID         uint               `gorm:"primaryKey"`
		ExternalID uuid.UUID          `gorm:"size:36;uniqueIndex;not null"`
		Name       string             `gorm:"size:255;not null"`
		Price      float64            `gorm:"precision:10;scale:2;not null"`
		Discounts  []DiscountProducts `gorm:"foreignKey:ProductID;references:ID"`
		Order      []OrderProduct     `gorm:"foreignKey:ProductID;references:ID"`
		CreateAt   time.Time          `gorm:"not null"`
		UpdateAt   time.Time          `gorm:"not null"`
		Version    uint               `gorm:"not null"`
	}

	Discount struct {
		ID         uint               `gorm:"primaryKey"`
		ExternalID uuid.UUID          `gorm:"size:36;uniqueIndex;not null"`
		Name       string             `gorm:"size:255;not null"`
		Price      float64            `gorm:"precision:10;scale:2;not null"`
		Products   []DiscountProducts `gorm:"foreignKey:DiscountID;references:ID"`
		CreateAt   time.Time          `gorm:"not null"`
		UpdateAt   time.Time          `gorm:"not null"`
		Version    uint               `gorm:"not null"`
	}

	DiscountProducts struct {
		DiscountID uint     `gorm:"primaryKey;autoIncrement:false;not null"`
		Discount   Discount `gorm:"foreignKey:DiscountID;references:ID"`
		ProductID  uint     `gorm:"primaryKey;autoIncrement:false;not null"`
		Product    Product  `gorm:"foreignKey:ProductID;references:ID"`
		Quantity   uint     `gorm:"not null"`
	}

	Order struct {
		ID         uint           `gorm:"primaryKey"`
		ExternalID uuid.UUID      `gorm:"size:36;uniqueIndex;not null"`
		Address    string         `gorm:"size:1000;not null"`
		DeliveryAt time.Time      `gorm:"not null"`
		PickUp     time.Time      `gorm:"not null"`
		TotalPrice float64        `gorm:"precision:10;scale:2;not null"`
		Discount   float64        `gorm:"precision:10;scale:2;not null"`
		FinalPrice float64        `gorm:"precision:10;scale:2;not null"`
		CustomerID uint           `gorm:"not null;index"`
		Customer   Customer       `gorm:"foreignKey:CustomerID;references:ID"`
		Payments   []OrderPayment `gorm:"foreignKey:OrderID;references:ID"`
		Products   []OrderProduct `gorm:"foreignKey:OrderID;references:ID"`
		CreateAt   time.Time      `gorm:"not null"`
		UpdateAt   time.Time      `gorm:"not null"`
		Version    uint           `gorm:"not null"`
		Comment    string
	}

	OrderPayment struct {
		ID      uuid.UUID `gorm:"size:36;primaryKey;not null"`
		Method  string    `gorm:"size:50;not null"`
		Info    string    `gorm:"size:100;not null"`
		Value   float64   `gorm:"precision:10;scale:2;not null"`
		At      time.Time `gorm:"not null"`
		OrderID uint      `gorm:"not null;index"`
		Order   Order     `gorm:"foreignKey:OrderID;references:ID"`
	}

	OrderProduct struct {
		OrderID   uint    `gorm:"primaryKey;autoIncrement:false;not null"`
		Order     Order   `gorm:"foreignKey:OrderID;references:ID"`
		ProductID uint    `gorm:"primaryKey;autoIncrement:false;not null"`
		Product   Product `gorm:"foreignKey:ProductID;references:ID"`
		Quantity  uint    `gorm:"not null"`
		Price     float64 `gorm:"precision:10;scale:2;not null"`
	}
)
