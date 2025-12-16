package models

import (
	"time"
)

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=255"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02"`
}

// UserResponse represents the response body for a user
type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
}

// UserWithAgeResponse represents the response body for a user with calculated age
type UserWithAgeResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  int    `json:"age"`
}

// PaginatedUsersResponse represents the paginated response for users
type PaginatedUsersResponse struct {
	Users      []UserWithAgeResponse `json:"users"`
	Total      int64                 `json:"total"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"page_size"`
	TotalPages int                   `json:"total_pages"`
}

// PaginationParams represents pagination query parameters
type PaginationParams struct {
	Page     int `query:"page"`
	PageSize int `query:"page_size"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// CalculateAge calculates age from date of birth
func CalculateAge(dob time.Time) int {
	now := time.Now()
	years := now.Year() - dob.Year()

	// Check if birthday hasn't occurred yet this year
	if now.Month() < dob.Month() || (now.Month() == dob.Month() && now.Day() < dob.Day()) {
		years--
	}

	return years
}

// DateFormat is the standard date format used in the application
const DateFormat = "2006-01-02"

// ParseDate parses a date string in the standard format
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse(DateFormat, dateStr)
}

// FormatDate formats a time.Time to the standard date string
func FormatDate(t time.Time) string {
	return t.Format(DateFormat)
}
