package service

import (
	"context"
	"errors"

	"github.com/dob-calculator/internal/logger"
	"github.com/dob-calculator/internal/models"
	"github.com/dob-calculator/internal/repository"
	"go.uber.org/zap"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input")
)

// UserService defines the interface for user business logic
type UserService interface {
	CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error)
	GetUserByID(ctx context.Context, id int32) (*models.UserWithAgeResponse, error)
	ListUsers(ctx context.Context, page, pageSize int) (*models.PaginatedUsersResponse, error)
	UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id int32) error
}

// userService implements UserService
type userService struct {
	repo repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(ctx context.Context, req models.CreateUserRequest) (*models.UserResponse, error) {
	dob, err := models.ParseDate(req.DOB)
	if err != nil {
		logger.Error("Failed to parse DOB", zap.Error(err), zap.String("dob", req.DOB))
		return nil, ErrInvalidInput
	}

	user, err := s.repo.Create(ctx, req.Name, dob)
	if err != nil {
		logger.Error("Failed to create user", zap.Error(err), zap.String("name", req.Name))
		return nil, err
	}

	logger.Info("User created successfully", zap.Int32("user_id", user.ID), zap.String("name", user.Name))

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  models.FormatDate(user.Dob),
	}, nil
}

// GetUserByID retrieves a user by ID with calculated age
func (s *userService) GetUserByID(ctx context.Context, id int32) (*models.UserWithAgeResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("Failed to get user by ID", zap.Error(err), zap.Int32("user_id", id))
		return nil, ErrUserNotFound
	}

	age := models.CalculateAge(user.Dob)

	logger.Info("User retrieved successfully", zap.Int32("user_id", user.ID))

	return &models.UserWithAgeResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  models.FormatDate(user.Dob),
		Age:  age,
	}, nil
}

// ListUsers retrieves all users with pagination and calculated ages
func (s *userService) ListUsers(ctx context.Context, page, pageSize int) (*models.PaginatedUsersResponse, error) {
	// Set defaults
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	users, err := s.repo.List(ctx, int32(pageSize), int32(offset))
	if err != nil {
		logger.Error("Failed to list users", zap.Error(err))
		return nil, err
	}

	total, err := s.repo.Count(ctx)
	if err != nil {
		logger.Error("Failed to count users", zap.Error(err))
		return nil, err
	}

	// Convert to response with age
	var usersWithAge []models.UserWithAgeResponse
	for _, user := range users {
		usersWithAge = append(usersWithAge, models.UserWithAgeResponse{
			ID:   user.ID,
			Name: user.Name,
			DOB:  models.FormatDate(user.Dob),
			Age:  models.CalculateAge(user.Dob),
		})
	}

	// Calculate total pages
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	logger.Info("Users listed successfully", zap.Int("count", len(users)), zap.Int64("total", total))

	return &models.PaginatedUsersResponse{
		Users:      usersWithAge,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateUser updates a user
func (s *userService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (*models.UserResponse, error) {
	dob, err := models.ParseDate(req.DOB)
	if err != nil {
		logger.Error("Failed to parse DOB", zap.Error(err), zap.String("dob", req.DOB))
		return nil, ErrInvalidInput
	}

	user, err := s.repo.Update(ctx, id, req.Name, dob)
	if err != nil {
		logger.Error("Failed to update user", zap.Error(err), zap.Int32("user_id", id))
		return nil, ErrUserNotFound
	}

	logger.Info("User updated successfully", zap.Int32("user_id", user.ID))

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  models.FormatDate(user.Dob),
	}, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(ctx context.Context, id int32) error {
	// Check if user exists first
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		logger.Error("User not found for deletion", zap.Error(err), zap.Int32("user_id", id))
		return ErrUserNotFound
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		logger.Error("Failed to delete user", zap.Error(err), zap.Int32("user_id", id))
		return err
	}

	logger.Info("User deleted successfully", zap.Int32("user_id", id))

	return nil
}
