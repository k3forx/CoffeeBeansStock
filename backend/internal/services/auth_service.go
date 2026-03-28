package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/auth"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/repository"
)

type AuthService struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{userRepo: userRepo, jwtManager: jwtManager}
}

type AuthResult struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}

type UserResponse struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name,omitempty"`
	LowStockThreshold   int32    `json:"low_stock_threshold,omitempty"`
	NotificationEnabled bool      `json:"notification_enabled,omitempty"`
	CreatedAt           string    `json:"created_at"`
	UpdatedAt           string    `json:"updated_at,omitempty"`
}

func (s *AuthService) RegisterAnonymous(ctx context.Context) (*AuthResult, error) {
	user, err := s.userRepo.CreateAnonymous(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := uuidFromPgtype(user.ID)
	if err != nil {
		return nil, err
	}

	tokens, err := s.jwtManager.GenerateTokenPair(userID)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         toUserResponse(user),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

type RefreshResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*RefreshResult, error) {
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, auth.ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, auth.ErrInvalidToken
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, auth.ErrInvalidToken
	}

	uid, err := uuidFromPgtype(user.ID)
	if err != nil {
		return nil, err
	}

	tokens, err := s.jwtManager.GenerateTokenPair(uid)
	if err != nil {
		return nil, err
	}

	return &RefreshResult{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *AuthService) GetMe(ctx context.Context, userID uuid.UUID) (*UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return toUserResponse(user), nil
}

func toUserResponse(user database.User) *UserResponse {
	resp := &UserResponse{
		LowStockThreshold:   user.LowStockThreshold,
		NotificationEnabled: user.NotificationEnabled,
	}
	if user.ID.Valid {
		resp.ID = uuid.UUID(user.ID.Bytes)
	}
	if user.Name.Valid {
		resp.Name = user.Name.String
	}
	if user.CreatedAt.Valid {
		resp.CreatedAt = user.CreatedAt.Time.Format("2006-01-02T15:04:05Z")
	}
	if user.UpdatedAt.Valid {
		resp.UpdatedAt = user.UpdatedAt.Time.Format("2006-01-02T15:04:05Z")
	}
	return resp
}

func uuidFromPgtype(id pgtype.UUID) (uuid.UUID, error) {
	if !id.Valid {
		return uuid.Nil, errors.New("invalid UUID")
	}
	return uuid.UUID(id.Bytes), nil
}
