package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// ErrInvalidToken is returned when a token is malformed or has an invalid signature.
var ErrInvalidToken = errors.New("invalid token")

// ErrExpiredToken is returned when a token has expired.
var ErrExpiredToken = errors.New("expired token")

// AccessTokenDuration is the lifetime of an access token.
const AccessTokenDuration = 1 * time.Hour

// RefreshTokenDuration is the lifetime of a refresh token.
const RefreshTokenDuration = 7 * 24 * time.Hour

// TokenClaims represents the JWT claims including the user ID.
type TokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// TokenPair holds an access token and a refresh token.
type TokenPair struct {
	AccessToken  string `json:"access_token"`  //nolint:gosec // JWT response field, not a hardcoded secret
	RefreshToken string `json:"refresh_token"` //nolint:gosec // JWT response field, not a hardcoded secret
}

// JWTManager handles JWT token generation and validation.
type JWTManager struct {
	secret []byte
}

// NewJWTManager creates a new JWTManager with the given secret.
func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{secret: []byte(secret)}
}

// GenerateTokenPair creates a new access/refresh token pair for the given user.
func (m *JWTManager) GenerateTokenPair(userID uuid.UUID) (*TokenPair, error) {
	accessToken, err := m.generateToken(userID.String(), AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := m.generateToken(userID.String(), RefreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// ValidateToken parses and validates a JWT token string.
func (m *JWTManager) ValidateToken(tokenStr string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

func (m *JWTManager) generateToken(userID string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := &TokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}
