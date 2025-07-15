package handlers

import (
	"net/http"
	"strconv"
	"user-management-api/internal/middleware"
	"user-management-api/internal/models"
	"user-management-api/internal/services"
	"user-management-api/pkg/errors"
	"user-management-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Get the profile of the currently authenticated user
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.APIResponse{data=models.UserResponse} "Profile retrieved successfully"
// @Failure      401  {object}  models.APIResponse "Unauthorized"
// @Failure      500  {object}  models.APIResponse "Internal server error"
// @Router       /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := middleware.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, models.APIResponse{
				Success: false,
				Message: appErr.Message,
				Error:   appErr.Type,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    user,
	})
}

// GetUser godoc
// @Summary      Get a user by ID
// @Description  Get a single user by their ID (Admin only)
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Security     BearerAuth
// @Success      200  {object}  models.APIResponse{data=models.UserResponse} "User retrieved successfully"
// @Failure      400  {object}  models.APIResponse "Invalid user ID"
// @Failure      404  {object}  models.APIResponse "User not found"
// @Failure      500  {object}  models.APIResponse "Internal server error"
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, models.APIResponse{
				Success: false,
				Message: appErr.Message,
				Error:   appErr.Type,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user (Admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.CreateUserRequest  true  "New User Info"
// @Security     BearerAuth
// @Success      201  {object}  models.APIResponse{data=models.UserResponse} "User created successfully"
// @Failure      400  {object}  models.APIResponse{error=map[string]string} "Validation failed or invalid request"
// @Failure      409  {object}  models.APIResponse "User already exists"
// @Failure      500  {object}  models.APIResponse "Internal server error"
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Validation failed",
			Error:   utils.FormatValidationError(err, models.CreateUserRequest{}),
		})
		return
	}

	user, err := h.userService.Create(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, models.APIResponse{
				Success: false,
				Message: appErr.Message,
				Error:   appErr.Type,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Update an existing user's details by ID (Admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string                   true  "User ID"
// @Param        user  body      models.UpdateUserRequest  true  "User Update Info"
// @Security     BearerAuth
// @Success      200  {object}  models.APIResponse{data=models.UserResponse} "User updated successfully"
// @Failure      400  {object}  models.APIResponse{error=map[string]string} "Validation failed or invalid request"
// @Failure      404  {object}  models.APIResponse "User not found"
// @Failure      500  {object}  models.APIResponse "Internal server error"
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Validation failed",
			Error:   utils.FormatValidationError(err, models.UpdateUserRequest{}),
		})
		return
	}

	user, err := h.userService.Update(c.Request.Context(), userID, &req)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, models.APIResponse{
				Success: false,
				Message: appErr.Message,
				Error:   appErr.Type,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	})
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user by their ID (Admin only)
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Security     BearerAuth
// @Success      200  {object}  models.APIResponse "User deleted successfully"
// @Failure      400  {object}  models.APIResponse "Invalid user ID"
// @Failure      404  {object}  models.APIResponse "User not found"
// @Failure      500  {object}  models.APIResponse "Internal server error"
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	err = h.userService.Delete(c.Request.Context(), userID)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, models.APIResponse{
				Success: false,
				Message: appErr.Message,
				Error:   appErr.Type,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}

// ListUsers godoc
// @Summary      List users
// @Description  Get a paginated list of all users (Admin only)
// @Tags         users
// @Produce      json
// @Param        page   query     int  false  "Page number"  default(1)
// @Param        limit  query     int  false  "Items per page" default(10)
// @Security     BearerAuth
// @Success      200  {object}  models.PaginatedUserResponse "Users retrieved successfully"
// @Failure      401  {object}  models.APIResponse "Unauthorized"
// @Failure      500  {object}  models.APIResponse "Internal server error"
// @Router       /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	result, err := h.userService.List(c.Request.Context(), page, limit)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, models.APIResponse{
				Success: false,
				Message: appErr.Message,
				Error:   appErr.Type,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}
