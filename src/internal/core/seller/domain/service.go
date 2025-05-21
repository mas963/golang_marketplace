package domain

import "context"

type SellerService interface {
	RegisterSeller(ctx context.Context, userID uint, input RegisterSellerInput) (*Seller, error)
	GetSellerByID(ctx context.Context, id uint) (*SellerDetail, error)
	UpdateSeller(ctx context.Context, id uint, input UpdateSellerInput) (*Seller, error)
	AddProductToSeller(ctx context.Context, sellerID uint, input AddSellerProductInput) (*SellerProduct, error)
	UpdateSellerProduct(ctx context.Context, id uint, input UpdateSellerProductInput) (*SellerProduct, error)
	ListSellerProducts(ctx context.Context, sellerID uint, pagination Pagination) (*SellerProductPagination, error)
	GetSellersByProduct(ctx context.Context, productID uint, pagination Pagination) (*SellersPagination, error)
}