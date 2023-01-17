package service

import (
	"context"
	"errors"
	"time"

	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/internal/services/dto"
	"github.com/sonyamoonglade/sancho-backend/pkg/auth"
)

var (
	ErrInvalidPassword = errors.New("password is invalid")
)

type TTLStrategy struct {
	AccessTokenTTLs map[domain.Role]time.Duration
	RefreshTokenTTL map[domain.Role]time.Duration
}

type authService struct {
	tokenProvider  auth.TokenProvider
	passwordHasher Hasher
	userService    User
	ttlStrategy    TTLStrategy
}

func NewAuthService(userService User, tokenProvider auth.TokenProvider, hasher Hasher, ttlStrategy TTLStrategy) Auth {
	return &authService{
		userService:    userService,
		tokenProvider:  tokenProvider,
		passwordHasher: hasher,
		ttlStrategy:    ttlStrategy,
	}
}

func (a authService) RegisterAdmin(ctx context.Context, dto dto.RegisterAdminDTO) (string, error) {
	// Create an admin without session. It's acquired via login.
	admin := domain.Admin{
		Login:    dto.Login,
		Role:     domain.RoleAdmin,
		Password: a.passwordHasher.Hash(dto.Password),
	}

	adminID, err := a.userService.SaveAdmin(ctx, admin)
	if err != nil {
		return "", err
	}

	return adminID, nil
}

func (a authService) LoginAdmin(ctx context.Context, loginDto dto.LoginAdminDTO) (auth.Pair, error) {
	var (
		accessTokenTTL  = a.ttlStrategy.AccessTokenTTLs[domain.RoleAdmin]
		refreshTokenTTL = a.ttlStrategy.RefreshTokenTTL[domain.RoleAdmin]
	)
	admin, err := a.userService.GetAdminByLogin(ctx, loginDto.Login)
	if err != nil {
		return auth.Pair{}, err
	}

	hashedDTOPassword := a.passwordHasher.Hash(loginDto.Password)
	if hashedDTOPassword != admin.Password {
		return auth.Pair{}, ErrInvalidPassword
	}

	tokens, err := a.tokenProvider.GenerateNewPairWithTTL(auth.UserAuth{
		Role:   domain.RoleAdmin,
		UserID: admin.UserID.Hex(),
	}, accessTokenTTL)
	if err != nil {
		return auth.Pair{}, err
	}

	session := domain.Session{
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    domain.NewExpiresAt(refreshTokenTTL),
	}
	err = a.userService.SaveSession(ctx, dto.SaveSessionDTO{
		UserID:  admin.UserID.Hex(),
		Role:    admin.Role,
		Session: session,
	})
	if err != nil {
		return auth.Pair{}, err
	}

	return tokens, nil
}

func (a authService) RefreshAdminToken(ctx context.Context, adminID, token string) (auth.Pair, error) {
	var (
		accessTokenTTL  = a.ttlStrategy.AccessTokenTTLs[domain.RoleAdmin]
		refreshTokenTTL = a.ttlStrategy.RefreshTokenTTL[domain.RoleAdmin]
	)

	admin, err := a.userService.GetAdminByRefreshToken(ctx, adminID, token)
	if err != nil {
		return auth.Pair{}, err
	}

	tokens, err := a.tokenProvider.GenerateNewPairWithTTL(auth.UserAuth{
		Role:   domain.RoleAdmin,
		UserID: admin.UserID.Hex(),
	}, accessTokenTTL)
	if err != nil {
		return auth.Pair{}, err
	}

	session := domain.Session{
		RefreshToken: tokens.RefreshToken,
		ExpiresAt:    domain.NewExpiresAt(refreshTokenTTL),
	}

	err = a.userService.SaveSession(ctx, dto.SaveSessionDTO{
		UserID:  admin.UserID.Hex(),
		Role:    admin.Role,
		Session: session,
	})
	if err != nil {
		return auth.Pair{}, err
	}

	return tokens, nil
}
