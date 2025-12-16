package repository

import (
	"context"
	"time"

	"github.com/dob-calculator/db/sqlc"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, name string, dob time.Time) (sqlc.User, error)
	GetByID(ctx context.Context, id int32) (sqlc.User, error)
	List(ctx context.Context, limit, offset int32) ([]sqlc.User, error)
	Update(ctx context.Context, id int32, name string, dob time.Time) (sqlc.User, error)
	Delete(ctx context.Context, id int32) error
	Count(ctx context.Context) (int64, error)
}

// userRepository implements UserRepository
type userRepository struct {
	queries *sqlc.Queries
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(queries *sqlc.Queries) UserRepository {
	return &userRepository{
		queries: queries,
	}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, name string, dob time.Time) (sqlc.User, error) {
	return r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Name: name,
		Dob:  dob,
	})
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id int32) (sqlc.User, error) {
	return r.queries.GetUserByID(ctx, id)
}

// List retrieves all users with pagination
func (r *userRepository) List(ctx context.Context, limit, offset int32) ([]sqlc.User, error) {
	return r.queries.ListUsers(ctx, sqlc.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

// Update updates a user
func (r *userRepository) Update(ctx context.Context, id int32, name string, dob time.Time) (sqlc.User, error) {
	return r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:   id,
		Name: name,
		Dob:  dob,
	})
}

// Delete deletes a user
func (r *userRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}

// Count returns the total number of users
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}
