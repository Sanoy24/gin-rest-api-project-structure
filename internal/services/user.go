package services

import (
	"context"
	"math"
	"user-management-api/internal/models"
	"user-management-api/internal/repository/interfaces"
	"user-management-api/pkg/errors"
	"user-management-api/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	userRepo interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetByID(ctx context.Context, id primitive.ObjectID) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrInternalServer
	}
	return user.ToResponse(), nil
}

func (s *UserService) Create(ctx context.Context, req *models.CreateUserRequest) (*models.UserResponse, error) {
	if _, err := s.userRepo.GetByEmail(ctx, req.Email); err == nil {
		return nil, errors.ErrUserExists
	}
	if _, err := s.userRepo.GetByUsername(ctx, req.Username); err == nil {
		return nil, errors.ErrUserExists
	}
	// hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// create user
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      req.Role,
		IsActive:  true,
	}

	if err = s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.ErrInternalServer
	}

	return user.ToResponse(), nil
}

func (s *UserService) Update(ctx context.Context, id primitive.ObjectID, req *models.UpdateUserRequest) (*models.UserResponse, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ErrUserNotFound
		}
		return nil, errors.ErrInternalServer
	}

	// Update fields if provided
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, errors.ErrInternalServer
	}

	return user.ToResponse(), nil
}

func (s *UserService) Delete(ctx context.Context, id primitive.ObjectID) error {
	// Check if user exists
	if _, err := s.userRepo.GetByID(ctx, id); err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.ErrUserNotFound
		}
		return errors.ErrInternalServer
	}

	return s.userRepo.Delete(ctx, id)
}

func (s *UserService) List(ctx context.Context, page, limit int) (*models.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	users, total, err := s.userRepo.List(ctx, page, limit)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Convert to response format
	userResponses := make([]*models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &models.PaginatedResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    userResponses,
		Pagination: models.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      int(total),
			TotalPages: totalPages,
		},
	}, nil
}
