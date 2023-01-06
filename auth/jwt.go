package auth

import (
	"errors"
	"time"

	"github.com/cristalhq/jwt/v4"
	"github.com/google/uuid"
	"github.com/sonyamoonglade/sancho-backend/internal/domain"
)

var (
	ErrTokenExpired         = errors.New("token has expired")
	ErrNotEnoughPermissions = errors.New("not enough permissions")
	ErrInvalidIssuer        = errors.New("invalid issuer")
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
	Validate(token string, role domain.Role) (bool, error)
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

func (p provider) Validate(token string, requiredRole domain.Role) (bool, error) {
	var (
		now        = time.Now()
		tokenBytes = []byte(token)
	)
	jwtToken, err := jwt.Parse(tokenBytes, p.verifier)
	if err != nil {
		return false, err
	}

	if err := p.verifier.Verify(jwtToken); err != nil {
		return false, err
	}

	var tokenPayload Claims
	if err := jwt.ParseClaims(tokenBytes, p.verifier, &tokenPayload); err != nil {
		return false, err
	}

	if tokenPayload.ExpiresAt.Before(now) {
		return false, ErrTokenExpired
	}

	if !tokenPayload.Role.CheckPermissions(requiredRole) {
		return false, ErrNotEnoughPermissions
	}

	if tokenPayload.Issuer != p.issuer {
		return false, ErrInvalidIssuer
	}

	// Validation successful
	return true, nil
}
