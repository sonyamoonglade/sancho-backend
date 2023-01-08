package auth

import (
	"errors"
	"time"

	"github.com/cristalhq/jwt/v4"
	"github.com/google/uuid"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
	"github.com/sonyamoonglade/sancho-backend/logger"
	"go.uber.org/zap"
)

var (
	ErrTokenExpired  = errors.New("token has expired")
	ErrInvalidIssuer = errors.New("invalid issuer")
	ErrInvalidToken  = errors.New("invalid token")
)

type Pair struct {
	RefreshToken string
	AccessToken  string
}

type UserAuth struct {
	// User role
	Role domain.Role `json:"rl"`
	// UserID of issuer
	UserID string `json:"usrId"`
}

type Claims struct {
	UserAuth
	Issuer    string    `json:"iss"`
	ExpiresAt time.Time `json:"exp"`
}

type TokenProvider interface {
	GenerateNewPair(data UserAuth) (Pair, error)
	// GenerateNewPairWithTTL will return access token with custom ttl
	GenerateNewPairWithTTL(data UserAuth, ttl time.Duration) (Pair, error)
	ParseAndValidate(token string) (UserAuth, error)
}

type provider struct {
	// Duration added to time.Now() in ExpiresAt field
	ttl        time.Duration
	signingKey []byte

	builder  *jwt.Builder
	verifier *jwt.HSAlg

	issuer string
}

func NewProvider(ttl time.Duration, signingKey []byte, issuer string) (TokenProvider, error) {
	signer, err := jwt.NewSignerHS(jwt.HS256, signingKey)
	if err != nil {
		return nil, err
	}
	verifier, err := jwt.NewVerifierHS(jwt.HS256, signingKey)
	if err != nil {
		return nil, err
	}
	return &provider{
		ttl:        ttl,
		signingKey: signingKey,
		builder:    jwt.NewBuilder(signer),
		verifier:   verifier,
		issuer:     issuer,
	}, nil
}

func (p provider) GenerateNewPair(data UserAuth) (Pair, error) {
	payload := Claims{
		UserAuth:  data,
		Issuer:    p.issuer,
		ExpiresAt: time.Now().Add(p.ttl),
	}
	accessToken, err := p.builder.Build(payload)
	if err != nil {
		return Pair{}, err
	}

	return Pair{
		RefreshToken: uuid.NewString(),
		AccessToken:  accessToken.String(),
	}, nil
}

func (p provider) GenerateNewPairWithTTL(data UserAuth, ttl time.Duration) (Pair, error) {
	payload := Claims{
		UserAuth:  data,
		Issuer:    p.issuer,
		ExpiresAt: time.Now().Add(ttl),
	}
	accessToken, err := p.builder.Build(payload)
	if err != nil {
		return Pair{}, err
	}

	return Pair{
		RefreshToken: uuid.NewString(),
		AccessToken:  accessToken.String(),
	}, nil
}

func (p provider) ParseAndValidate(token string) (UserAuth, error) {
	var (
		tokenPayload Claims
		now          = time.Now()
		tokenBytes   = []byte(token)
	)
	if err := jwt.ParseClaims(tokenBytes, p.verifier, &tokenPayload); err != nil {
		if errors.Is(err, jwt.ErrInvalidFormat) || errors.Is(err, jwt.ErrInvalidSignature) {
			return UserAuth{}, ErrInvalidToken
		}
		logger.Get().Error("jwt.ParseClaims: ", zap.Error(err))
		return UserAuth{}, err
	}
	if tokenPayload.ExpiresAt.Before(now) {
		return UserAuth{}, ErrTokenExpired
	}
	if tokenPayload.Issuer != p.issuer {
		return UserAuth{}, ErrInvalidIssuer
	}
	// Validation successful
	return tokenPayload.UserAuth, nil
}
