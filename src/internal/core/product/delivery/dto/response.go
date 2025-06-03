package dto

import "golang_marketplace/src/internal/core/product/domain"

type ErrorResponse struct {
	Error string `json:"error"`
}

type ProductListResponse struct {
	Products []*domain.ProductWithVariants `json:"products"`
	Total    int64                         `json:"total"`
	Page     int                           `json:"page"`
	Limit    int                           `json:"limit"`
}
