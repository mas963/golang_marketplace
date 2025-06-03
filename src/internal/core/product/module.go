package product

import (
	"github.com/gin-gonic/gin"
	"golang_marketplace/src/internal/core/product/delivery/http"
	"golang_marketplace/src/internal/core/product/domain"
	"golang_marketplace/src/internal/core/product/repository"
	"golang_marketplace/src/internal/core/product/service"
	"golang_marketplace/src/internal/platform/cache"
	"gorm.io/gorm"
)

type Module struct {
	Service domain.ProductService
}

func NewModule(db *gorm.DB, cache *cache.Cache) *Module {
	productRepo := repository.NewProductRepository(db)
	variantRepo := repository.NewProductVariantRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	productService := service.NewProductService(productRepo, variantRepo, categoryRepo, cache)

	return &Module{
		Service: productService,
	}
}

func (m *Module) RegisterRoutes(router *gin.RouterGroup) {
	http.RegisterRoutes(router, m.Service)
}
