package domain

import "context"

type SellerService interface {
	RegisterSeller(ctx context.Context, userID uint, input RegisterSellerInput) (*Seller, error)
	GetSellerByID(ctx context.Context, id uint) (*Seller, error)
	GetSellerBySlug(ctx context.Context, slug string) (*Seller, error)
	UpdateSeller(ctx context.Context, id uint, input UpdateSellerInput) (*Seller, error)
	AdminUpdateSeller(ctx context.Context, id uint, input AdminUpdateSellerInput) (*Seller, error)
	DeleteSeller(ctx context.Context, id uint) error
	AddProductToSeller(ctx context.Context, sellerID uint, input AddSellerProductInput) (*SellerProduct, error)
	UpdateSellerProduct(ctx context.Context, id uint, input UpdateSellerProductInput) (*SellerProduct, error)
	ListSellerProducts(ctx context.Context, sellerID uint, filter ProductFilter) (*SellerProductPagination, error)
	GetSellersByProduct(ctx context.Context, productID uint, pagination Pagination) (*SellersPagination, error)
	UpdateSellerStats(ctx context.Context, sellerID uint, statsUpdate SellerStatsUpdate) error
}
