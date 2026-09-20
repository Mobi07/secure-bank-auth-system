package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

	if err := godotenv.Load("../.env"); err != nil {
		log.Error("Root .env not found; using system environment variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("configuration initialization failed", zap.Error(err))
	}

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

	auth := api.Group("/auth")

	// Public
	auth.POST("/auth/register", authHandler.Register)
	auth.POST("/auth/login", authHandler.Login)

	// Authenticated
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtService))

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
