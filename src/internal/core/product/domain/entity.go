package domain

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CategoryID  uuid.UUID `json:"category_id" gorm:"type:uuid;not null" validate:"required"`
	Brand       string    `json:"brand"`
	SKU         string    `json:"sku"`
	Status      string    `json:"status" gorm:"default:'active'" validate:"oneof=active inactive"`
	Images      []string  `json:"images" gorm:"type:text[]"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductVariant struct {
	ID            uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID     uuid.UUID              `json:"product_id" gorm:"type:uuid;not null"`
	SellerID      uuid.UUID              `json:"seller_id" gorm:"type:uuid;not null"`
	Price         float64                `json:"price" validate:"required,gt=0"`
	DiscountPrice *float64               `json:"discount_price,omitempty" validate:"omitempty,gt=0"`
	Stock         int                    `json:"stock" validate:"gte=0"`
	Attributes    map[string]interface{} `json:"attributes" gorm:"type:jsonb"`
	IsActive      bool                   `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type Category struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name      string     `json:"name" gorm:"not null" validate:"required"`
	ParentID  *uuid.UUID `json:"parent_id,omitempty" gorm:"type:uuid"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Seller struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CompanyName string    `json:"company_name" gorm:"not null"`
	Email       string    `json:"email" gorm:"unique;not null"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	TaxNumber   string    `json:"tax_number" gorm:"unique;"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
