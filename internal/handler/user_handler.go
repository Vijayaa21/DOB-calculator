package handler

import (
	"errors"
	"strconv"

	"github.com/dob-calculator/internal/logger"
	"github.com/dob-calculator/internal/models"
	"github.com/dob-calculator/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	userService service.UserService
	validate    *validator.Validate
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		validate:    validator.New(),
	}
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: formatValidationErrors(err),
		})
	}

	user, err := h.userService.CreateUser(c.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "invalid_input",
				Message: "Invalid date format. Use YYYY-MM-DD",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create user",
		})
	}

	logger.Info("User created via API", zap.Int32("user_id", user.ID))
	return c.Status(fiber.StatusCreated).JSON(user)
}

// GetUserByID handles GET /users/:id
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		logger.Error("Invalid user ID", zap.String("id", idParam))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid user ID",
		})
	}

	user, err := h.userService.GetUserByID(c.Context(), int32(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error:   "not_found",
				Message: "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get user",
		})
	}

	return c.JSON(user)
}

// ListUsers handles GET /users
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	result, err := h.userService.ListUsers(c.Context(), page, pageSize)
	if err != nil {
		logger.Error("Failed to list users", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to list users",
		})
	}

	// Return just the array of users for backward compatibility
	// or the full paginated response based on query params
	if c.Query("paginated", "false") == "true" {
		return c.JSON(result)
	}

	return c.JSON(result.Users)
}

// UpdateUser handles PUT /users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		logger.Error("Invalid user ID", zap.String("id", idParam))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid user ID",
		})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid request body",
		})
	}

	if err := h.validate.Struct(req); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "validation_error",
			Message: formatValidationErrors(err),
		})
	}

	user, err := h.userService.UpdateUser(c.Context(), int32(id), req)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error:   "not_found",
				Message: "User not found",
			})
		}
		if errors.Is(err, service.ErrInvalidInput) {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Error:   "invalid_input",
				Message: "Invalid date format. Use YYYY-MM-DD",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to update user",
		})
	}

	logger.Info("User updated via API", zap.Int32("user_id", user.ID))
	return c.JSON(user)
}

// DeleteUser handles DELETE /users/:id
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		logger.Error("Invalid user ID", zap.String("id", idParam))
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Error:   "bad_request",
			Message: "Invalid user ID",
		})
	}

	err = h.userService.DeleteUser(c.Context(), int32(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse{
				Error:   "not_found",
				Message: "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete user",
		})
	}

	logger.Info("User deleted via API", zap.Int64("user_id", id))
	return c.SendStatus(fiber.StatusNoContent)
}

// formatValidationErrors formats validation errors into a readable message
func formatValidationErrors(err error) string {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		var messages []string
		for _, e := range errs {
			switch e.Tag() {
			case "required":
				messages = append(messages, e.Field()+" is required")
			case "min":
				messages = append(messages, e.Field()+" must be at least "+e.Param()+" characters")
			case "max":
				messages = append(messages, e.Field()+" must be at most "+e.Param()+" characters")
			case "datetime":
				messages = append(messages, e.Field()+" must be in format YYYY-MM-DD")
			default:
				messages = append(messages, e.Field()+" is invalid")
			}
		}
		if len(messages) > 0 {
			return messages[0]
		}
	}
	return "Validation failed"
}
