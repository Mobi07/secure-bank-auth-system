package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mobi07/secure_bank_auth/internal/auth"
	"github.com/mobi07/secure_bank_auth/internal/config"
	"github.com/mobi07/secure_bank_auth/internal/database"
	"github.com/mobi07/secure_bank_auth/internal/domain"
	"github.com/mobi07/secure_bank_auth/internal/handler"
	"github.com/mobi07/secure_bank_auth/internal/logger"
	"github.com/mobi07/secure_bank_auth/internal/middleware"
	"github.com/mobi07/secure_bank_auth/internal/repository/postgres"
	"github.com/mobi07/secure_bank_auth/internal/service"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.New()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	cfg := config.Load()

	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("database initialization failed: ", zap.Error(err))
	}
	defer db.Close()

	jwtService := auth.NewJWTService(
		cfg.JWTSecret,
		15*time.Minute,
	)

	userRepo := postgres.NewUserRepository(db)

	authService := service.NewAuthService(
		userRepo,
		jwtService,
		log,
	)

	authHandler := handler.NewAuthHandler(authService)
	dashboardHandler := handler.NewDashboardHandler()
	adminHandler := handler.NewAdminHandler()

	router := gin.Default()

	api := router.Group("/api/v1")

	// Public
	api.POST("/auth/register", authHandler.Register)
	api.POST("/auth/login", authHandler.Login)

	// Authenticated
	protected := api.Group("/")
	protected.Use()

	protected.GET("/dashboard", dashboardHandler.GetDashboard)

	// Admin
	admin := api.Group("/admin")
	admin.Use(
		middleware.RequireRole(domain.RoleAdmin),
	)

	admin.GET("/users", adminHandler.GetUsers)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	log.Info("starting SecureBank API", zap.String("address", ":8080"))

	if err := router.Run(":8080"); err != nil {
		log.Fatal("server failed to start ", zap.Error(err))
	}
}
