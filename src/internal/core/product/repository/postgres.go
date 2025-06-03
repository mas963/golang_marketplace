package repository

import (
	"context"
	"github.com/google/uuid"
	"golang_marketplace/src/internal/core/product/domain"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		First(&product, "sku = ?", sku).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.Product{}, id).Error
}

func (r *productRepository) List(ctx context.Context, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	var products []*domain.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.Product{}).Preload("Category")

	query = r.applyFilters(query, filter)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	if filter.SortBy != "" {
		orderBy := filter.SortBy
		if filter.SortOrder == "desc" {
			orderBy += " DESC"
		} else {
			orderBy += " ASC"
		}
		query = query.Order(orderBy)
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Search(ctx context.Context, query string, filter domain.ProductFilter) ([]*domain.Product, int64, error) {
	var products []*domain.Product
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&domain.Product{}).
		Preload("Category").
		Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")

	dbQuery = r.applyFilters(dbQuery, filter)

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		dbQuery = dbQuery.Offset(offset).Limit(filter.Limit)
	}

	if err := dbQuery.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) applyFilters(query *gorm.DB, filter domain.ProductFilter) *gorm.DB {
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", filter.CategoryID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	return query
}

type productVariantRepository struct {
	db *gorm.DB
}

func NewProductVariantRepository(db *gorm.DB) domain.ProductVariantRepository {
	return &productVariantRepository{db: db}
}

func (r *productVariantRepository) Create(ctx context.Context, variant *domain.ProductVariant) error {
	return r.db.WithContext(ctx).Create(variant).Error
}

func (r *productVariantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.ProductVariant, error) {
	var variant domain.ProductVariant
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Seller").
		First(&variant, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *productVariantRepository) GetByProductID(ctx context.Context, productID uuid.UUID) ([]*domain.ProductVariant, error) {
	var variants []*domain.ProductVariant
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("product_id = ? AND is_active = ?", productID, true).
		Find(&variants).Error
	return variants, err
}

func (r *productVariantRepository) GetBySellerID(ctx context.Context, sellerID uuid.UUID, filter domain.ProductFilter) ([]*domain.ProductVariant, int64, error) {
	var variants []*domain.ProductVariant
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.ProductVariant{}).
		Preload("Product").
		Where("seller_id = ?", sellerID)

	if filter.MinPrice != nil {
		query = query.Where("price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		query = query.Where("price <= ?", *filter.MaxPrice)
	}
	if filter.InStock != nil && *filter.InStock {
		query = query.Where("quantity > 0")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Page > 0 && filter.Limit > 0 {
		offset := (filter.Page - 1) * filter.Limit
		query = query.Offset(offset).Limit(filter.Limit)
	}

	if err := query.Find(&variants).Error; err != nil {
		return nil, 0, err
	}

	return variants, total, nil
}

func (r *productVariantRepository) Update(ctx context.Context, variant *domain.ProductVariant) error {
	return r.db.WithContext(ctx).Save(variant).Error
}

func (r *productVariantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.ProductVariant{}, id).Error
}

func (r *productVariantRepository) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	return r.db.WithContext(ctx).Model(&domain.ProductVariant{}).
		Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", quantity)).Error
}
