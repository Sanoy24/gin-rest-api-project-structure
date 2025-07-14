package handlers

import (
	"net/http"
	"time"
	"user-management-api/internal/models"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) HealthCheck(ctx *gin.Context) {
	ctx.JSON((http.StatusOK), models.APIResponse{
		Success: true,
		Message: "Service is running",
		Data: gin.H{
			"status":    "OK",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})
}
