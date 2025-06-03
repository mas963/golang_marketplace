package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"golang_marketplace/src/internal/core/product/domain"
	"golang_marketplace/src/internal/platform/cache"
	"golang_marketplace/src/pkg/validator"
	"gorm.io/gorm"
	"time"
)

type productService struct {
	productRepo  domain.ProductRepository
	variantRepo  domain.ProductVariantRepository
	categoryRepo domain.CategoryRepository
	cache        *cache.Cache
}

func NewProductService(
	productRepo domain.ProductRepository,
	variantRepo domain.ProductVariantRepository,
	categoryRepo domain.CategoryRepository,
	cache *cache.Cache,
) domain.ProductService {
	return &productService{
		productRepo:  productRepo,
		variantRepo:  variantRepo,
		categoryRepo: categoryRepo,
		cache:        cache,
	}
}

func (s *productService) CreateProduct(ctx context.Context, req domain.CreateProductRequest) (*domain.Product, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	// sku must be unique
	if _, err := s.productRepo.GetBySKU(ctx, req.SKU); err == nil {
		return nil, fmt.Errorf("sku %s already exists", req.SKU)
	} else if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("failed to check if sku %s exists: %w", req.SKU, err)
	}

	// category must exist
	if _, err := s.categoryRepo.GetByID(ctx, req.CategoryID); err != nil {
		return nil, fmt.Errorf("category %s does not exist", req.CategoryID)
	}

	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Brand:       req.Brand,
		SKU:         req.SKU,
		Status:      "active",
		Images:      req.Images,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (s *productService) GetProduct(ctx context.Context, id uuid.UUID) (*domain.ProductWithVariants, error) {
	cacheKey := fmt.Sprintf("product:%s", id.String())
	var cached domain.ProductWithVariants
	if err := s.cache.Get(cacheKey, &cached); err == nil {
		return &cached, nil
	}

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	variants, err := s.variantRepo.GetByProductID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product variants: %w", err)
	}

	result := &domain.ProductWithVariants{
		Product:  product,
		Variants: variants,
	}

	s.cache.Set(cacheKey, result, 5*time.Minute)

	return result, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id uuid.UUID, req domain.UpdateProductRequest) (*domain.Product, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.CategoryID != nil {
		if _, err := s.categoryRepo.GetByID(ctx, *req.CategoryID); err != nil {
			return nil, fmt.Errorf("category not found: %w", *req.CategoryID)
		}
		product.CategoryID = *req.CategoryID
	}
	if req.Brand != nil {
		product.Brand = *req.Brand
	}
	if req.Images != nil {
		product.Images = req.Images
	}
	if req.Status != nil {
		product.Status = *req.Status
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	cacheKey := fmt.Sprintf("product:%s", id.String())
	_ = s.cache.Delete(cacheKey)

	return product, nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if err := s.productRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	cacheKey := fmt.Sprintf("product:%s", id.String())
	_ = s.cache.Delete(cacheKey)

	return nil
}

func (s *productService) ListProducts(ctx context.Context, filter domain.ProductFilter) ([]*domain.ProductWithVariants, int64, error) {
	products, total, err := s.productRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	var result []*domain.ProductWithVariants
	for _, product := range products {
		variants, err := s.variantRepo.GetByProductID(ctx, product.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get variants for product %s: %w", product.ID, err)
		}
		result = append(result, &domain.ProductWithVariants{
			Product:  product,
			Variants: variants,
		})
	}

	return result, total, nil
}

func (s *productService) SearchProducts(ctx context.Context, query string, filter domain.ProductFilter) ([]*domain.ProductWithVariants, int64, error) {
	products, total, err := s.productRepo.Search(ctx, query, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search products: %w", err)
	}

	var result []*domain.ProductWithVariants
	for _, product := range products {
		variants, err := s.variantRepo.GetByProductID(ctx, product.ID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get variants for product %s: %w", product.ID, err)
		}
		result = append(result, &domain.ProductWithVariants{
			Product:  product,
			Variants: variants,
		})
	}

	return result, total, nil
}

func (s *productService) AddProductVariant(ctx context.Context, req domain.CreateVariantRequest) (*domain.ProductVariant, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	if _, err := s.productRepo.GetByID(ctx, req.ProductID); err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	if req.DiscountPrice != nil && *req.DiscountPrice >= req.Price {
		return nil, fmt.Errorf("discount price must be less than price")
	}

	variant := &domain.ProductVariant{
		ProductID:     req.ProductID,
		SellerID:      req.SellerID,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		Stock:         req.Stock,
		Attributes:    req.Attributes,
		IsActive:      true,
	}

	if err := s.variantRepo.Create(ctx, variant); err != nil {
		return nil, fmt.Errorf("failed to create variant: %w", err)
	}

	cacheKey := fmt.Sprintf("product:%s", req.ProductID.String())
	s.cache.Delete(cacheKey)

	return variant, nil
}

func (s *productService) UpdateProductVariant(ctx context.Context, id uuid.UUID, req domain.UpdateVariantRequest) (*domain.ProductVariant, error) {
	if err := validator.ValidateStruct(req); err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	variant, err := s.variantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("variant not found: %w", err)
	}

	if req.Price != nil {
		variant.Price = *req.Price
	}
	if req.DiscountPrice != nil {
		variant.DiscountPrice = req.DiscountPrice
	}
	if req.Stock != nil {
		variant.Stock = *req.Stock
	}
	if req.Attributes != nil {
		variant.Attributes = req.Attributes
	}
	if req.IsActive != nil {
		variant.IsActive = *req.IsActive
	}

	if variant.DiscountPrice != nil && *variant.DiscountPrice >= variant.Price {
		return nil, fmt.Errorf("discount price must be less than regular price")
	}

	if err := s.variantRepo.Update(ctx, variant); err != nil {
		return nil, fmt.Errorf("failed to update variant: %w", err)
	}

	cacheKey := fmt.Sprintf("product:%s", variant.ProductID.String())
	_ = s.cache.Delete(cacheKey)

	return variant, nil
}

func (s *productService) DeleteProductVariant(ctx context.Context, id uuid.UUID) error {
	variant, err := s.variantRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("variant not found: %w", err)
	}

	if err := s.variantRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete variant: %w", err)
	}

	cacheKey := fmt.Sprintf("product:%s", variant.ProductID.String())
	_ = s.cache.Delete(cacheKey)

	return nil
}

func (s *productService) GetVariantsByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.ProductVariant, error) {
	variants, err := s.variantRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get variants: %w", err)
	}
	return variants, nil
}

func (s *productService) GetVariantsBySeller(ctx context.Context, sellerID uuid.UUID, filter domain.ProductFilter) ([]*domain.ProductVariant, int64, error) {
	variants, total, err := s.variantRepo.GetBySellerID(ctx, sellerID, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get variants: %w", err)
	}
	return variants, total, nil
}

func (s *productService) UpdateStock(ctx context.Context, variantID uuid.UUID, quantity int) error {
	if err := s.variantRepo.UpdateStock(ctx, variantID, quantity); err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	variant, err := s.variantRepo.GetByID(ctx, variantID)
	if err == nil {
		cacheKey := fmt.Sprintf("product:%s", variant.ProductID.String())
		_ = s.cache.Delete(cacheKey)
	}

	return nil
}

func (s *productService) CheckStock(ctx context.Context, variantID uuid.UUID, quantity int) (bool, error) {
	variant, err := s.variantRepo.GetByID(ctx, variantID)
	if err != nil {
		return false, fmt.Errorf("variant not found: %w", err)
	}

	return variant.Stock >= quantity, nil
}
