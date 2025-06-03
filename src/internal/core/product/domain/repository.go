package domain

import (
	"context"
	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	GetBySKU(ctx context.Context, sku string) (*Product, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filter ProductFilter) ([]*Product, int64, error)
	Search(ctx context.Context, query string, filter ProductFilter) ([]*Product, int64, error)
}

type ProductVariantRepository interface {
	Create(ctx context.Context, variant *ProductVariant) error
	GetByID(ctx context.Context, id uuid.UUID) (*ProductVariant, error)
	GetByProductID(ctx context.Context, productID uuid.UUID) ([]*ProductVariant, error)
	GetBySellerID(ctx context.Context, sellerID uuid.UUID, filter ProductFilter) ([]*ProductVariant, int64, error)
	Update(ctx context.Context, variant *ProductVariant) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*Category, error)
	List(ctx context.Context) ([]*Category, error)
	Update(ctx context.Context, category *Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ProductFilter struct {
	CategoryID *uuid.UUID
	SellerID   *uuid.UUID
	MinPrice   *float64
	MaxPrice   *float64
	InStock    *bool
	Status     string
	Page       int
	Limit      int
	SortBy     string
	SortOrder  string
}
