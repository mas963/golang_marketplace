package domain

import "time"

type Seller struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"uniqueIndex"`
	StoreName    string    `json:"store_name"`
	Slug         string    `json:"slug" gorm:"uniqueIndex"`
	Status       string    `json:"status"` // pending, active, suspended, banned
	Rating       float64   `json:"rating"`
	RatingCount  int       `json:"rating_count"`
	ContactEmail string    `json:"contact_email"`
	ContactPhone string    `json:"contact_phone"`
	Address      string    `json:"address"`
	City         string    `json:"city"`
	District     string    `json:"district"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	User     *User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Products []SellerProduct `json:"products,omitempty" gorm:"foreignKey:SellerID"`
}

type RegisterSellerInput struct {
	StoreName    string `json:"store_name"`
	ContactEmail string `json:"contact_email"`
	ContactPhone string `json:"contact_phone"`
	Address      string `json:"address"`
	City         string `json:"city"`
	District     string `json:"district"`
}

type UpdateSellerInput struct {
	StoreName    *string `json:"store_name"`
	ContactEmail *string `json:"contact_email"`
	ContactPhone *string `json:"contact_phone"`
	Address      *string `json:"address"`
	City         *string `json:"city"`
	District     *string `json:"district"`
}

type AdminUpdateSellerInput struct {
	UpdateSellerInput
	ComissionRate *float64 `json:"comission_rate"`
	IsVerified    *bool    `json:"is_verified"`
	Status        *string  `json:"status"`
}

type SellerProduct struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	SellerID       uint      `json:"seller_id" gorm:"index"`
	ProductID      uint      `json:"product_id" gorm:"index"`
	Price          float64   `json:"price"`
	DiscountPrice  float64   `json:"discount_price"`
	DiscountType   string    `json:"discount_type"` // percentage, fixed
	DiscountRate   float64   `json:"discount_rate"`
	StockQuantity  int       `json:"stock_quantity"`
	StockCode      string    `json:"stock_code"`
	IsActive       bool      `json:"is_active"`
	IsFeatured     bool      `json:"is_featured"`
	ShippingTime   string    `json:"shipping_time"`
	ShippingCost   float64   `json:"shipping_cost"`
	ShippingOption string    `json:"shipping_option"` // free, flat_rate, variable
	TaxRate        float64   `json:"tax_rate"`
	SalesCount     int       `json:"sales_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Seller  *Seller  `json:"seller,omitempty" gorm:"foreignKey:SellerID"`
	Product *Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

type AddSellerProductInput struct {
	ProductID      uint    `json:"product_id"`
	Price          float64 `json:"price"`
	DiscountPrice  float64 `json:"discount_price"`
	DiscountType   string  `json:"discount_type"`
	DiscountRate   float64 `json:"discount_rate"`
	StockQuantity  int     `json:"stock_quantity"`
	StockCode      string  `json:"stock_code"`
	IsActive       bool    `json:"is_active"`
	IsFeatured     bool    `json:"is_featured"`
	ShippingTime   string  `json:"shipping_time"`
	ShippingCost   float64 `json:"shipping_cost"`
	ShippingOption string  `json:"shipping_option"`
	TaxRate        float64 `json:"tax_rate"`
}

type UpdateSellerProductInput struct {
	Price          *float64 `json:"price"`
	DiscountPrice  *float64 `json:"discount_price"`
	DiscountType   *string  `json:"discount_type"`
	DiscountRate   *float64 `json:"discount_rate"`
	StockQuantity  *int     `json:"stock_quantity"`
	StockCode      *string  `json:"stock_code"`
	IsActive       *bool    `json:"is_active"`
	IsFeatured     *bool    `json:"is_featured"`
	ShippingTime   *string  `json:"shipping_time"`
	ShippingCost   *float64 `json:"shipping_cost"`
	ShippingOption *string  `json:"shipping_option"`
	TaxRate        *float64 `json:"tax_rate"`
}

type SellerProductDetail struct {
	SellerProduct
	SellerName    string  `json:"seller_name"`
	SellerRating  float64 `json:"seller_rating"`
	ProductName   string  `json:"product_name"`
	ProductImage  string  `json:"product_image"`
	CategoryName  string  `json:"category_name"`
	Brand         string  `json:"brand"`
	TotalReviews  int     `json:"total_reviews"`
	AverageRating float64 `json:"average_rating"`
}

type SellerProductPagination struct {
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	Limit      int                   `json:"limit"`
	TotalPages int                   `json:"total_pages"`
	HasNext    bool                  `json:"has_next"`
	HasPrev    bool                  `json:"has_prev"`
	Items      []SellerProductDetail `json:"items"`
}
