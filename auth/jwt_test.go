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

	t.Run("should validate token, return true and payload", func(t *testing.T) {
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

		userAuth, err := p.ParseAndValidate(token.String())
		require.NoError(t, err)
		require.EqualValues(t, tokenPayload.UserAuth, userAuth)
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

		userAuth, err := p.ParseAndValidate(tokens.AccessToken)
		require.NoError(t, err)
		require.EqualValues(t, data, userAuth)
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

		userAuth, err := p.ParseAndValidate(tokens.AccessToken)
		require.Error(t, err)
		require.Equal(t, ErrTokenExpired, err)
		require.Zero(t, userAuth)
	})

	t.Run("should generate new pair from 1st issuer and validate on 2nd issuer and retrieve ErrInvalidIssuer", func(t *testing.T) {
		p1 := p
		p2, err := NewProvider(ttl, key, "some-random-issuer")
		require.NoError(t, err)

		tokensP1, err := p1.GenerateNewPair(UserAuth{
			Role:   domain.RoleCustomer,
			UserID: uuid.NewString(),
		})
		require.NoError(t, err)

		// Required role is OK but p1 and p2 has got different issuers, so expect error
		userAuth, err := p2.ParseAndValidate(tokensP1.AccessToken)
		require.Error(t, err)
		require.Equal(t, ErrInvalidIssuer, err)
		require.Zero(t, userAuth)
	})

	t.Run("should generate new pair with custom ttl and return ErrExpired because new ttl is lower than default", func(t *testing.T) {
		var (
			// default ttl
			ttl    = time.Millisecond * 5
			key    = []byte("abcd")
			issuer = "app.com"
		)
		p, err := NewProvider(ttl, key, issuer)
		require.NoError(t, err)

		customTTl := time.Millisecond * 2
		tokens, err := p.GenerateNewPairWithTTL(UserAuth{
			Role:   domain.RoleCustomer,
			UserID: uuid.NewString(),
		}, customTTl)
		require.NoError(t, err)

		// access token has ttl about 2 milliseconds so if we sleep for 2ms
		// token will be expired by that time, because default ttl is overridden.
		time.Sleep(time.Millisecond * 2)

		userAuth, err := p.ParseAndValidate(tokens.AccessToken)
		require.Error(t, err)
		require.Equal(t, ErrTokenExpired, err)
		require.Zero(t, userAuth)
	})

	t.Run("should generate new pair with custom ttl and return true because new ttl is much higher than default", func(t *testing.T) {
		var (
			// default ttl
			ttl    = time.Millisecond * 2
			key    = []byte("abcd")
			issuer = "app.com"
		)
		p, err := NewProvider(ttl, key, issuer)
		require.NoError(t, err)

		customTTl := time.Millisecond * 200
		tokens, err := p.GenerateNewPairWithTTL(UserAuth{
			Role:   domain.RoleCustomer,
			UserID: uuid.NewString(),
		}, customTTl)
		require.NoError(t, err)

		// access token has ttl about 200ms and thus if we sleep for 4ms (higher than default ttl of 2ms)
		// we'll get absolutely fine result and token will be valid.
		time.Sleep(time.Millisecond * 4)

		userAuth, err := p.ParseAndValidate(tokens.AccessToken)
		require.NoError(t, err)
		require.NotZero(t, userAuth)
	})

	t.Run("should return ErrInvalidToken because token is in invalid format", func(t *testing.T) {
		token := "some-random-token"
		userAuth, err := p.ParseAndValidate(token)
		require.Error(t, err)
		require.Zero(t, userAuth)
		require.Equal(t, ErrInvalidToken, err)
	})

	t.Run("should return ErrInvalidToken because token has corrupted header", func(t *testing.T) {
		tokens, _ := p.GenerateNewPair(UserAuth{
			Role:   domain.RoleCustomer,
			UserID: uuid.NewString(),
		})
		var (
			byteToken      = []byte(tokens.AccessToken)
			corruptedToken string
		)
		// Corrupt. Don't change 0th, 1st, 2nd chars, see jwt/parse.go:39
		byteToken[4] = byte('o')
		byteToken[5] = byte('j')
		byteToken[6] = byte('m')
		corruptedToken = string(byteToken)

		userAuth, err := p.ParseAndValidate(corruptedToken)
		require.Error(t, err)
		require.Zero(t, userAuth)
		require.Equal(t, ErrInvalidToken, err)
	})

	t.Run("should return ErrInvalidToken because token has corrupted signature", func(t *testing.T) {
		tokens, _ := p.GenerateNewPair(UserAuth{
			Role:   domain.RoleCustomer,
			UserID: uuid.NewString(),
		})
		var (
			byteToken      = []byte(tokens.AccessToken)
			corruptedToken string
		)
		tl := len(byteToken)
		byteToken[tl-2] = byte('o')
		byteToken[tl-3] = byte('j')
		byteToken[tl-4] = byte('m')
		corruptedToken = string(byteToken)

		userAuth, err := p.ParseAndValidate(corruptedToken)
		require.Error(t, err)
		require.Zero(t, userAuth)
		require.Equal(t, ErrInvalidToken, err)
	})
}
