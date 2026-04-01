package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/k3forx/CoffeeBeansStock/backend/internal/auth"
	domain "github.com/k3forx/CoffeeBeansStock/backend/internal/domain"
	domainauth "github.com/k3forx/CoffeeBeansStock/backend/internal/domain/auth"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/domain/user"
)

// AuthService handles authentication business logic.
type AuthService struct {
	userRepo      user.Repository
	tokens        domainauth.TokenManager
	refreshTokens domainauth.RefreshTokenRepository
}

// NewAuthService creates a new AuthService.
func NewAuthService(userRepo user.Repository, tokens domainauth.TokenManager, refreshTokens domainauth.RefreshTokenRepository) *AuthService {
	return &AuthService{userRepo: userRepo, tokens: tokens, refreshTokens: refreshTokens}
}

// AuthResult holds the result of an authentication operation.
type AuthResult struct {
	User         *user.User
	AccessToken  string
	RefreshToken string
}

// RegisterAnonymous creates a new anonymous user and returns tokens.
func (s *AuthService) RegisterAnonymous(ctx context.Context) (*AuthResult, error) {
	u := user.NewAnonymousUser()

	tokens, err := s.tokens.GenerateTokenPair(u.ID())
	if err != nil {
		return nil, err
	}

	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, err
	}

	tokenHash := auth.HashToken(tokens.RefreshToken)
	expiresAt := time.Now().UTC().Add(auth.RefreshTokenDuration)
	if err := s.refreshTokens.Store(ctx, u.ID(), tokenHash, expiresAt); err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         u,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// RefreshInput holds input for refreshing tokens.
type RefreshInput struct {
	RefreshToken string
}

// RefreshResult holds new tokens after a refresh operation.
type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}

// Refresh validates a refresh token and issues a new token pair.
func (s *AuthService) Refresh(ctx context.Context, in RefreshInput) (*RefreshResult, error) {
	claims, err := s.tokens.ValidateToken(in.RefreshToken)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	u, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, domain.ErrInvalidToken
		}
		return nil, err
	}

	oldHash := auth.HashToken(in.RefreshToken)
	exists, err := s.refreshTokens.ExistsByHash(ctx, oldHash)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, domain.ErrInvalidToken
	}
	if err = s.refreshTokens.DeleteByHash(ctx, oldHash); err != nil {
		return nil, err
	}

	tokens, err := s.tokens.GenerateTokenPair(u.ID())
	if err != nil {
		return nil, err
	}

	newHash := auth.HashToken(tokens.RefreshToken)
	expiresAt := time.Now().UTC().Add(auth.RefreshTokenDuration)
	if err := s.refreshTokens.Store(ctx, u.ID(), newHash, expiresAt); err != nil {
		return nil, err
	}

	return &RefreshResult{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

// GetMeInput holds input for getting the current user.
type GetMeInput struct {
	UserID uuid.UUID
}

// GetMeOutput holds the result of getting the current user.
type GetMeOutput struct {
	User *user.User
}

// GetMe returns the user for the given user ID.
func (s *AuthService) GetMe(ctx context.Context, in GetMeInput) (*GetMeOutput, error) {
	u, err := s.userRepo.GetByID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}
	return &GetMeOutput{User: u}, nil
}

// UpdateMeInput holds input for updating the current user's settings.
type UpdateMeInput struct {
	UserID            uuid.UUID
	GramsPerCup       *int32
	LowStockThreshold *int32
}

// UpdateMeOutput holds the result of updating the current user.
type UpdateMeOutput struct {
	User *user.User
}

// UpdateMe updates the authenticated user's settings.
func (s *AuthService) UpdateMe(ctx context.Context, in UpdateMeInput) (*UpdateMeOutput, error) {
	u, err := s.userRepo.GetByID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	if err := u.Update(in.GramsPerCup, in.LowStockThreshold); err != nil {
		return nil, err
	}

	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	return &UpdateMeOutput{User: u}, nil
}
