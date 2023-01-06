package auth

import (
	"testing"
	"time"

	"github.com/cristalhq/jwt/v4"
	"github.com/google/uuid"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestTokenProvider(t *testing.T) {
	var (
		key    = []byte("abcd")
		ttl    = time.Millisecond * 5
		issuer = "app.com"
	)

	p, err := NewProvider(ttl, key, issuer)
	require.NoError(t, err)

	t.Run("should validate valid token and return true", func(t *testing.T) {
		userID := uuid.NewString()
		signer, err := jwt.NewSignerHS(jwt.HS256, key)
		require.NoError(t, err)
		builder := jwt.NewBuilder(signer)
		tokenPayload := Claims{
			UserAuth: UserAuth{
				Role:   domain.RoleCustomer,
				UserID: userID,
			},
			Issuer:    issuer,
			ExpiresAt: time.Now().Add(ttl),
		}
		token, err := builder.Build(tokenPayload)
		require.NoError(t, err)

		// REQUIRED ROLE
		requiredRole := domain.RoleCustomer
		ok, err := p.Validate(token.String(), requiredRole)
		require.NoError(t, err)
		require.True(t, ok)
	})
	t.Run("should validate invalid token and return false and ErrNotEnoughPermissions because required role is higher", func(t *testing.T) {
		userID := uuid.New().String()
		signer, err := jwt.NewSignerHS(jwt.HS256, key)
		require.NoError(t, err)
		builder := jwt.NewBuilder(signer)
		tokenPayload := Claims{
			// Sign with RoleCustomer
			UserAuth: UserAuth{
				Role:   domain.RoleCustomer,
				UserID: userID,
			},
			Issuer:    issuer,
			ExpiresAt: time.Now().Add(ttl),
		}
		token, err := builder.Build(tokenPayload)
		require.NoError(t, err)

		// REQUIRED ROLE (ADMIN)
		requiredRole := domain.RoleAdmin
		ok, err := p.Validate(token.String(), requiredRole)
		require.False(t, ok)
		require.Error(t, err)
		require.Equal(t, ErrNotEnoughPermissions, err)
	})

	t.Run("should generate new pair of refresh and access tokens", func(t *testing.T) {
		data := UserAuth{
			Role:   domain.RoleCustomer,
			UserID: uuid.NewString(),
		}
		tokens, err := p.GenerateNewPair(data)
		require.NoError(t, err)
		require.NotZero(t, tokens.AccessToken)
		require.NotZero(t, tokens.RefreshToken)
	})

	t.Run("should generate new pair and validate access token", func(t *testing.T) {
		data := UserAuth{
			Role:   domain.RoleCustomer,
			UserID: uuid.NewString(),
		}
		tokens, err := p.GenerateNewPair(data)
		require.NoError(t, err)

		requiredRole := domain.RoleCustomer
		ok, err := p.Validate(tokens.AccessToken, requiredRole)
		require.NoError(t, err)
		require.True(t, ok)
	})

	t.Run("should generate pair and validate access token and return false because role is too low", func(t *testing.T) {
		data := UserAuth{
			Role:   domain.RoleCustomer,
			UserID: uuid.NewString(),
		}
		tokens, err := p.GenerateNewPair(data)
		require.NoError(t, err)

		requiredRole := domain.RoleAdmin
		ok, err := p.Validate(tokens.AccessToken, requiredRole)
		require.Error(t, err)
		require.False(t, ok)
		require.Equal(t, ErrNotEnoughPermissions, err)
	})

	t.Run("should generate new pair and wait for expiration then validate and return false with ErrExpired", func(t *testing.T) {
		data := UserAuth{
			Role:   domain.RoleCustomer,
			UserID: uuid.NewString(),
		}

		// see ttl (5ms)
		tokens, err := p.GenerateNewPair(data)
		require.NoError(t, err)

		time.Sleep(time.Millisecond * 6)

		requiredRole := domain.RoleCustomer
		ok, err := p.Validate(tokens.AccessToken, requiredRole)
		require.Equal(t, ErrTokenExpired, err)
		require.False(t, ok)
		require.Error(t, err)
	})

	t.Run("should generate new pair from 1st issuer and validate on 2nd issuer and retreive ErrInvalidIssuer", func(t *testing.T) {
		p1 := p
		p2, err := NewProvider(ttl, key, "some-random-issuer")
		require.NoError(t, err)

		tokensP1, err := p1.GenerateNewPair(UserAuth{
			Role:   domain.RoleCustomer,
			UserID: uuid.NewString(),
		})
		require.NoError(t, err)

		// Required role is OK but p1 and p2 has got different issuers, so expect error
		ok, err := p2.Validate(tokensP1.AccessToken, domain.RoleCustomer)
		require.Error(t, err)
		require.False(t, ok)
		require.Equal(t, ErrInvalidIssuer, err)
	})
}
