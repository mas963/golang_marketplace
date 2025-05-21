package domain

import "context"

type OrderService interface {
	CreateOrder(ctx context.Context, customerID uint, input CreateOrderInput) (*Order, error)
	GetOrderByID(ctx context.Context, id uint) (*OrderDetail, error)
	UpdateOrderStatus(ctx context.Context, id uint, status string) (*Order, error)
	CancelOrder(ctx context.Context, id uint) error
	GetCustomerOrders(ctx context.Context, customerID uint, pagination Pagination) (*OrderPagination, error)
	GetSellerOrders(ctx context.Context, sellerID uint, pagination Pagination) (*OrderPagination, error)
}