package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	domainauth "github.com/k3forx/CoffeeBeansStock/backend/internal/domain/auth"
)

const (
	AccessTokenDuration  = 1 * time.Hour
	RefreshTokenDuration = 7 * 24 * time.Hour
)

// jwtClaims is the internal JWT claims structure.
type jwtClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token generation and validation.
// It implements domainauth.TokenManager.
type JWTManager struct {
	secret []byte
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{secret: []byte(secret)}
}

func (m *JWTManager) GenerateTokenPair(userID uuid.UUID) (*domainauth.TokenPair, error) {
	accessToken, err := m.generateToken(userID.String(), AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.generateToken(userID.String(), RefreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &domainauth.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *JWTManager) ValidateToken(tokenStr string) (*domainauth.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{UserID: "", RegisteredClaims: jwt.RegisteredClaims{}}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrExpiredToken
		}
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	return &domainauth.TokenClaims{
		UserID: claims.UserID,
	}, nil
}

func (m *JWTManager) generateToken(userID string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := &jwtClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}
