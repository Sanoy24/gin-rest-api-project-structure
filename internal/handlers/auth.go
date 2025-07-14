package handlers

import (
	"net/http"
	"user-management-api/internal/models"
	"user-management-api/internal/services"
	"user-management-api/pkg/errors"
	"user-management-api/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}
	if err := utils.ValidateStruct(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Validation failed",
			Error:   utils.FormatValidationError(err),
		})
		return
	}
	authResponse, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		if appError, ok := err.(*errors.AppError); ok {
			c.JSON(appError.Code, models.APIResponse{
				Success: false,
				Message: appError.Message,
				Error:   appError.Type,
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
		Data:    authResponse,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
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
			Error:   utils.FormatValidationError(err),
		})
		return
	}

	// Login user
	authResponse, err := h.authService.Login(c.Request.Context(), &req)
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
		Message: "Login successful",
		Data:    authResponse,
	})
}
