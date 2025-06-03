package domain

import (
	"context"
	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*ProductWithVariants, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error
	ListProducts(ctx context.Context, filter ProductFilter) ([]*ProductWithVariants, int64, error)
	SearchProducts(ctx context.Context, query string, filter ProductFilter) ([]*ProductWithVariants, int64, error)

	AddProductVariant(ctx context.Context, req CreateVariantRequest) (*ProductVariant, error)
	UpdateProductVariant(ctx context.Context, id uuid.UUID, req UpdateVariantRequest) (*ProductVariant, error)
	DeleteProductVariant(ctx context.Context, id uuid.UUID) error
	GetVariantsByProduct(ctx context.Context, productID uuid.UUID) ([]*ProductVariant, error)
	GetVariantsBySeller(ctx context.Context, sellerID uuid.UUID, filter ProductFilter) ([]*ProductVariant, int64, error)

	UpdateStock(ctx context.Context, variantID uuid.UUID, quantity int) error
	CheckStock(ctx context.Context, variantID uuid.UUID, quantity int) (bool, error)
}

type CreateProductRequest struct {
	Name        string    `json:"name" validate:"required,min=2,max=255"`
	Description string    `json:"description" validate:"max=2000"`
	CategoryID  uuid.UUID `json:"category_id" validate:"required"`
	Brand       string    `json:"brand" validate:"max=100"`
	SKU         string    `json:"sku" validate:"required"`
	Images      []string  `json:"images"`
}

type UpdateProductRequest struct {
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Description *string    `json:"description,omitempty" validate:"omitempty,max=2000"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	Brand       *string    `json:"brand,omitempty" validate:"omitempty,max=100"`
	Images      []string   `json:"images,omitempty"`
	Status      *string    `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
}

type CreateVariantRequest struct {
	ProductID     uuid.UUID              `json:"product_id" validate:"required"`
	SellerID      uuid.UUID              `json:"seller_id" validate:"required"`
	Price         float64                `json:"price" validate:"required,gt=0"`
	DiscountPrice *float64               `json:"discount_price,omitempty" validate:"omitempty,gt=0"`
	Stock         int                    `json:"stock" validate:"gte=0"`
	Attributes    map[string]interface{} `json:"attributes"`
}

type UpdateVariantRequest struct {
	Price         *float64               `json:"price,omitempty" validate:"omitempty,gt=0"`
	DiscountPrice *float64               `json:"discount_price,omitempty" validate:"omitempty,gt=0"`
	Stock         *int                   `json:"stock,omitempty" validate:"omitempty,gte=0"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
	IsActive      *bool                  `json:"is_active,omitempty"`
}

type ProductWithVariants struct {
	*Product
	Variants []*ProductVariant `json:"variants"`
}
