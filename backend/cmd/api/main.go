package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mobi07/secure_bank_auth/internal/config"
	"github.com/mobi07/secure_bank_auth/internal/database"
	"github.com/mobi07/secure_bank_auth/internal/logger"
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

	router := gin.Default()

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
