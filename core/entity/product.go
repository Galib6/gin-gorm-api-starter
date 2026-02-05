package entity

import (
	"myapp/support/base"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string          `json:"name" gorm:"not null"`
	Description string          `json:"description"`
	SKU         string          `json:"sku" gorm:"unique;not null"`
	Price       decimal.Decimal `json:"price" gorm:"type:decimal(15,2);not null"`
	Stock       int             `json:"stock" gorm:"not null;default:0"`
	CategoryID  *uuid.UUID      `json:"category_id" gorm:"type:uuid"`
	IsActive    bool            `json:"is_active" gorm:"not null;default:true"`
	Image       *string         `json:"image"`
	base.Model

	// Relations
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `json:"name" gorm:"unique;not null"`
	Description string    `json:"description"`
	base.Model

	// Relations
	Products []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}
