package domain

import "context"

type SellerRepository interface {
	Create(ctx context.Context, seller *Seller) error
	GetByID(ctx context.Context, id uint) (*Seller, error)
	GetByUserID(ctx context.Context, userID uint) (*Seller, error)
	GetBySlug(ctx context.Context, slug string) (*Seller, error)
	Update(ctx context.Context, seller *Seller) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, filter SellerFilter) ([]*Seller, int64, error)
	GetSellerDetail(ctx context.Context, id uint) (*Seller, error)
	GetSellersByProduct(ctx context.Context, productID uint, pagination Pagination) ([]*SellerProductDetail, int64, error)
	UpdateStats(ctx context.Context, sellerID uint, statsUpdate SellerStatsUpdate) error
}

type SellerProductRepository interface {
	Create(ctx context.Context, sellerProduct *SellerProduct) error
	GetByID(ctx context.Context, id uint) (*SellerProduct, error)
	GetBySellerIDAndProductID(ctx context.Context, sellerID, productID uint) (*SellerProduct, error)
	Update(ctx context.Context, sellerProduct *SellerProduct) error
	Delete(ctx context.Context, id uint) error
	ListBySellerID(ctx context.Context, sellerID uint, filter ProductFilter) ([]*SellerProductDetail, int64, error)
	ListByProductID(ctx context.Context, productID uint, pagination Pagination) ([]*SellerProductDetail, int64, error)
	UpdateStock(ctx context.Context, id uint, quantity int) error
	IncrementSalesCount(ctx context.Context, id uint, quantity int) error
}
