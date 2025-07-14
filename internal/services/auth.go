package services

import (
	"context"
	"time"
	"user-management-api/internal/models"
	"user-management-api/internal/repository/interfaces"
	"user-management-api/pkg/errors"
	"user-management-api/pkg/utils"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuthService struct {
	userRepo  interfaces.UserRepository
	jwtSecret string
	jwtExpiry string
}

func NewAuthService(userRepo interfaces.UserRepository, jwtSecret, jwtExpiry string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *AuthService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrInvalidCredentials
		}
		return nil, errors.ErrInternalServer
	}
	// Check if user is active
	if !user.IsActive {
		return nil, errors.ErrUnAuthorized
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user.ToResponse(),
	}, nil
}

func (s *AuthService) Register(ctx context.Context, req *models.CreateUserRequest) (*models.AuthResponse, error) {
	// Check if user already exists
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, errors.ErrUserExists
	}

	if _, err := s.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, errors.ErrUserExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Create user
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  true,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.ErrInternalServer
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, s.jwtSecret, 24*time.Hour)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return &models.AuthResponse{
		Token: token,
		User:  *user.ToResponse(),
	}, nil
}
