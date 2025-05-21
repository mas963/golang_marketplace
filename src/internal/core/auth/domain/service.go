package domain

import "context"

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*UserWithToken, error)
	Login(ctx context.Context, email, password string) (*UserWithToken, error)
	RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
	ValidateToken(ctx context.Context, token string) (*Claims, error)
}