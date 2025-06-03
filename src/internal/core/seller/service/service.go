package service

import (
	"golang_marketplace/src/internal/core/seller/domain"
	"golang_marketplace/src/internal/platform/cache"
)

type sellerService struct {
	repo        domain.SellerRepository
	productRepo domain.SellerProductRepository
	cache       cache.RedisCache
}
