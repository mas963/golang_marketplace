package domain

import "context"

type ProductService interface {
	CreateProduct(ctx context.Context, input CreateProductInput) (*Product, error)
	GetProductByID(ctx context.Context, id uint) (*ProductDetail, error)
	UpdateProduct(ctx context.Context, id uint, input UpdateProductInput) (*Product, error)
	DeleteProduct(ctx context.Context, id uint) error
	ListProducts(ctx context.Context, filter ProductFilter) (*ProductPagination, error)
	GetProductByCategory(ctx context.Context, categoryID uint, pagination Pagination) (*ProductPagination, error)
	SearchProducts(ctx context.Context, query string, filter ProductFilter) (*ProductPagination, error)
}