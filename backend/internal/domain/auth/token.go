package auth

import "github.com/google/uuid"

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenClaims struct {
	UserID string
}

type TokenManager interface {
	GenerateTokenPair(userID uuid.UUID) (*TokenPair, error)
	ValidateToken(tokenStr string) (*TokenClaims, error)
}
