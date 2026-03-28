package services

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/api"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/auth"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/database"
	"github.com/k3forx/CoffeeBeansStock/backend/internal/repository"
)

var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	userRepo   repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{userRepo: userRepo, jwtManager: jwtManager}
}

type SignupInput struct {
	Email    string
	Password string
	Name     string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	User         *UserResponse `json:"user"`
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
}

type UserResponse struct {
	ID                  uuid.UUID `json:"id"`
	Email               string    `json:"email"`
	Name                string    `json:"name"`
	LowStockThreshold   int32    `json:"low_stock_threshold,omitempty"`
	NotificationEnabled bool      `json:"notification_enabled,omitempty"`
	CreatedAt           string    `json:"created_at"`
	UpdatedAt           string    `json:"updated_at,omitempty"`
}

func (s *AuthService) ValidateSignupInput(input *SignupInput) []api.FieldError {
	var errs []api.FieldError

	input.Email = strings.TrimSpace(input.Email)
	input.Name = strings.TrimSpace(input.Name)

	if input.Email == "" {
		errs = append(errs, api.FieldError{Field: "email", Message: "メールアドレスは必須です"})
	} else if _, err := mail.ParseAddress(input.Email); err != nil {
		errs = append(errs, api.FieldError{Field: "email", Message: "有効なメールアドレスを入力してください"})
	}

	if input.Password == "" {
		errs = append(errs, api.FieldError{Field: "password", Message: "パスワードは必須です"})
	} else if len(input.Password) < 8 {
		errs = append(errs, api.FieldError{Field: "password", Message: "パスワードは8文字以上で入力してください"})
	}

	if input.Name == "" {
		errs = append(errs, api.FieldError{Field: "name", Message: "名前は必須です"})
	} else if len(input.Name) > 100 {
		errs = append(errs, api.FieldError{Field: "name", Message: "名前は100文字以内で入力してください"})
	}

	return errs
}

func (s *AuthService) ValidateLoginInput(input *LoginInput) []api.FieldError {
	var errs []api.FieldError

	input.Email = strings.TrimSpace(input.Email)

	if input.Email == "" {
		errs = append(errs, api.FieldError{Field: "email", Message: "メールアドレスは必須です"})
	}
	if input.Password == "" {
		errs = append(errs, api.FieldError{Field: "password", Message: "パスワードは必須です"})
	}

	return errs
}

func (s *AuthService) Signup(ctx context.Context, input *SignupInput) (*AuthResult, error) {
	_, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.Create(ctx, repository.CreateUserParams{
		Email:               input.Email,
		PasswordHash:        hash,
		Name:                input.Name,
		LowStockThreshold:   100,
		NotificationEnabled: true,
	})
	if err != nil {
		return nil, err
	}

	userID, err := uuidFromPgtype(user.ID)
	if err != nil {
		return nil, err
	}

	tokens, err := s.jwtManager.GenerateTokenPair(userID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		User:         toUserResponse(user),
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input *LoginInput) (*AuthResult, error) {
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := auth.CheckPassword(user.PasswordHash, input.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	userID, err := uuidFromPgtype(user.ID)
	if err != nil {
		return nil, err
	}

	tokens, err := s.jwtManager.GenerateTokenPair(userID, user.Email)
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

	tokens, err := s.jwtManager.GenerateTokenPair(uid, user.Email)
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
		Email:               user.Email,
		Name:                user.Name,
		LowStockThreshold:   user.LowStockThreshold,
		NotificationEnabled: user.NotificationEnabled,
	}
	if user.ID.Valid {
		resp.ID = uuid.UUID(user.ID.Bytes)
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
