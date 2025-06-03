package http

import (
	"github.com/gin-gonic/gin"
	"golang_marketplace/src/internal/core/product/domain"
)

func RegisterRoutes(router *gin.RouterGroup, service domain.ProductService) {
	handler := NewProductHandler(service)

	products := router.Group("/products")
	{
		products.POST("", handler.CreateProduct)
		products.GET("", handler.ListProducts)
		products.GET("/search", handler.SearchProducts)
		products.GET("/:id", handler.GetProduct)
		products.PUT("/:id", handler.UpdateProduct)
		products.DELETE("/:id", handler.DeleteProduct)
		products.GET("/:product_id/variants", handler.GetVariantsByProduct)
	}

	variants := router.Group("/products/variants")
	{
		variants.POST("", handler.AddProductVariant)
		variants.PUT("/:id", handler.UpdateProductVariant)
		variants.DELETE("/:id", handler.DeleteProductVariant)
	}
}
