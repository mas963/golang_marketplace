package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang_marketplace/src/internal/core/product/delivery/dto"
	"golang_marketplace/src/internal/core/product/domain"
	"log"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	service domain.ProductService
}

func NewProductHandler(service domain.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product
// @Tags products
// @Accept json
// @Produce json
// @Param product body domain.CreateProductRequest true "Product data"
// @Success 201 {object} domain.Product
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req domain.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	product, err := h.service.CreateProduct(c.Request.Context(), req)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Description Get a product by ID with variants
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} domain.ProductWithVariants
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid product ID"})
		return
	}

	product, err := h.service.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body domain.UpdateProductRequest true "Product update data"
// @Success 200 {object} domain.Product
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid product ID"})
		return
	}

	var req domain.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	product, err := h.service.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid product ID"})
		return
	}

	if err := h.service.DeleteProduct(c.Request.Context(), id); err != nil {
		log.Println("Failed to delete product: ", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListProducts godoc
// @Summary List products
// @Description List products with filters and pagination
// @Tags products
// @Accept json
// @Produce json
// @Param category_id query string false "Category ID"
// @Param seller_id query string false "Seller ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param in_stock query boolean false "In stock filter"
// @Param status query string false "Product status"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param sort_by query string false "Sort by field"
// @Param sort_order query string false "Sort order (asc/desc)"
// @Success 200 {object} ProductListResponse
// @Router /products [get]
func (h *ProductHandler) ListProducts(c *gin.Context) {
	filter := h.parseProductFilter(c)

	products, total, err := h.service.ListProducts(c.Request.Context(), filter)
	if err != nil {
		log.Println("Failed to list products: ", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	response := ProductListResponse{
		Products: products,
		Total:    total,
		Page:     filter.Page,
		Limit:    filter.Limit,
	}

	c.JSON(http.StatusOK, response)
}

// SearchProducts godoc
// @Summary Search products
// @Description Search products by name or description
// @Tags products
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param category_id query string false "Category ID"
// @Param min_price query number false "Minimum price"
// @Param max_price query number false "Maximum price"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} ProductListResponse
// @Router /products/search [get]
func (h *ProductHandler) SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Search query is required"})
		return
	}

	filter := h.parseProductFilter(c)

	products, total, err := h.service.SearchProducts(c.Request.Context(), query, filter)
	if err != nil {
		log.Println("Failed to search products: ", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	response := ProductListResponse{
		Products: products,
		Total:    total,
		Page:     filter.Page,
		Limit:    filter.Limit,
	}

	c.JSON(http.StatusOK, response)
}

// AddProductVariant godoc
// @Summary Add a variant to a product
// @Description Add a new variant to an existing product
// @Tags variants
// @Accept json
// @Produce json
// @Param variant body domain.CreateVariantRequest true "Variant data"
// @Success 201 {object} domain.ProductVariant
// @Failure 400 {object} ErrorResponse
// @Router /products/variants [post]
func (h *ProductHandler) AddProductVariant(c *gin.Context) {
	var req domain.CreateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	variant, err := h.service.AddProductVariant(c.Request.Context(), req)
	if err != nil {
		log.Println("Failed to add product variant: ", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, variant)
}

// UpdateProductVariant godoc
// @Summary Update a product variant
// @Description Update a product variant
// @Tags variants
// @Accept json
// @Produce json
// @Param id path string true "Variant ID"
// @Param variant body domain.UpdateVariantRequest true "Variant update data"
// @Success 200 {object} domain.ProductVariant
// @Failure 400 {object} ErrorResponse
// @Router /products/variants/{id} [put]
func (h *ProductHandler) UpdateProductVariant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid variant ID"})
		return
	}

	var req domain.UpdateVariantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	variant, err := h.service.UpdateProductVariant(c.Request.Context(), id, req)
	if err != nil {
		log.Println("Failed to update product variant: ", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, variant)
}

// DeleteProductVariant godoc
// @Summary Delete a product variant
// @Description Delete a product variant
// @Tags variants
// @Accept json
// @Produce json
// @Param id path string true "Variant ID"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Router /products/variants/{id} [delete]
func (h *ProductHandler) DeleteProductVariant(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid variant ID"})
		return
	}

	if err := h.service.DeleteProductVariant(c.Request.Context(), id); err != nil {
		log.Println("Failed to delete product variant: ", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetVariantsByProduct godoc
// @Summary Get variants by product ID
// @Description Get all variants for a specific product
// @Tags variants
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {array} domain.ProductVariant
// @Failure 400 {object} ErrorResponse
// @Router /products/{product_id}/variants [get]
func (h *ProductHandler) GetVariantsByProduct(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid product ID"})
		return
	}

	variants, err := h.service.GetVariantsByProduct(c.Request.Context(), productID)
	if err != nil {
		log.Println("Failed to get variants by product: ", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, variants)
}

func (h *ProductHandler) parseProductFilter(c *gin.Context) domain.ProductFilter {
	filter := domain.ProductFilter{
		Page:  1,
		Limit: 20,
	}

	if categoryIDStr := c.Query("category_id"); categoryIDStr != "" {
		if categoryID, err := uuid.Parse(categoryIDStr); err == nil {
			filter.CategoryID = &categoryID
		}
	}

	if sellerIDStr := c.Query("seller_id"); sellerIDStr != "" {
		if sellerID, err := uuid.Parse(sellerIDStr); err == nil {
			filter.SellerID = &sellerID
		}
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filter.MinPrice = &minPrice
		}
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filter.MaxPrice = &maxPrice
		}
	}

	if inStockStr := c.Query("in_stock"); inStockStr != "" {
		if inStock, err := strconv.ParseBool(inStockStr); err == nil {
			filter.InStock = &inStock
		}
	}

	if status := c.Query("status"); status != "" {
		filter.Status = status
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filter.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			filter.Limit = limit
		}
	}

	if sortBy := c.Query("sort_by"); sortBy != "" {
		filter.SortBy = sortBy
	}

	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		filter.SortOrder = sortOrder
	}

	return filter
}
