package domain

import "context"

type UserService interface {
	GetUserByID(ctx context.Context, id uint) (*User, error)
	UpdateUser(ctx context.Context, id uint, input UpdateUserInput) (*User, error)
	DeleteUser(ctx context.Context, id uint) error
	GetCustomerProfile(ctx context.Context, userID uint) (*CustomerProfile, error)
	GetSellerProfile(ctx context.Context, userID uint) (*SellerProfile, error)
}