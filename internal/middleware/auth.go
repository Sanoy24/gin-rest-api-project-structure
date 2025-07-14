package middleware

import (
	"errors"
	"net/http"
	"strings"
	"user-management-api/internal/config"
	"user-management-api/internal/models"
	"user-management-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMidddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Authorization header is required",
			})
			c.Abort()
			return
		}
		// check if a token starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Invalid authorization heade",
			})
			c.Abort()
			return
		}
		// extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(token, cfg.JWT.Secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.APIResponse{
				Success: false,
				Message: "Invalid or expired token",
			})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole, exists := ctx.Get("user_role")
		if !exists {
			ctx.JSON(http.StatusForbidden, models.APIResponse{
				Success: false,
				Message: "Unauthorized access",
			})
			ctx.Abort()
			return
		}
		role := userRole.(string)
		for _, requriedRoles := range roles {
			if role == requriedRoles {
				ctx.Next()
				return
			}
		}
		ctx.JSON(http.StatusForbidden, models.APIResponse{
			Success: false,
			Message: "Insufficient permissions",
		})
		ctx.Abort()
	}
}

func GetUserId(ctx *gin.Context) (primitive.ObjectID, error) {
	userId, exists := ctx.Get("user_id")
	if !exists {
		return primitive.NilObjectID, errors.New("user ID not found in context")
	}
	return userId.(primitive.ObjectID), nil
}
