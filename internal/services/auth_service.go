package service

import (
	"context"
	"time"

	"github.com/sonyamoonglade/sancho-backend/auth"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
)

type TTLStrategy struct {
	AccessTokenTTL struct {
		Admin, Worker, Customer time.Duration
	}
	RefreshTokenTTL struct {
		Admin, Worker, Customer time.Duration
	}
}

type authService struct {
	tokenProvider auth.TokenProvider
	ttlStrategy   TTLStrategy
}

func NewAuthService(tokenProvider auth.TokenProvider, ttlStrategy TTLStrategy) Auth {
	return &authService{
		tokenProvider: tokenProvider,
		ttlStrategy:   ttlStrategy,
	}
}

func (a authService) GetAdminByRefreshToken(ctx context.Context, token string) (domain.Admin, error) {
	return domain.Admin{}, nil
}
