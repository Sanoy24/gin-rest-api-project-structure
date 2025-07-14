package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-management-api/internal/config"
	"user-management-api/internal/handlers"
	"user-management-api/internal/middleware"
	"user-management-api/internal/repository/mongo"
	"user-management-api/internal/services"
	"user-management-api/pkg/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", err)
	}

	mongoDb, err := database.NewMongoDB(cfg.Database.URI, cfg.Database.Name, cfg.Database.Timeout)
	if err != nil {
		log.Fatal("failed to connect to mongodb")
	}
	defer mongoDb.Close(context.Background())
	// initialize repositories
	userRepo := mongo.NewUserRepository(mongoDb.Database)

	// initialize services
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiresIn.String())
	userService := services.NewUserService(userRepo)

	// initialize handler

	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// setup router
	router := setupRouter(cfg, healthHandler, authHandler, userHandler)

	// start server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.Server.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown ", err)
	}
	log.Println("server exited")

}

func setupRouter(cfg *config.Config, healthHandler *handlers.HealthHandler, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler) *gin.Engine {
	if cfg.Server.Env == "Production" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()

	// middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(gin.Recovery())

	router.GET("/health", healthHandler.HealthCheck)

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		users := v1.Group("/users")
		{
			users.GET("/profile", middleware.AuthMidddleware(cfg), userHandler.GetProfile)
			users.GET("", middleware.AuthMidddleware(cfg), middleware.RequireRole("admin"), userHandler.ListUsers)
			users.POST("", middleware.AuthMidddleware(cfg), middleware.RequireRole("admin"), userHandler.CreateUser)
			users.GET("/:id", middleware.AuthMidddleware(cfg), middleware.RequireRole("admin"), userHandler.GetUser)
			users.PUT("/:id", middleware.AuthMidddleware(cfg), middleware.RequireRole("admin"), userHandler.UpdateUser)
			users.DELETE("/:id", middleware.AuthMidddleware(cfg), middleware.RequireRole("admin"), userHandler.DeleteUser)
		}
	}
	return router

}
