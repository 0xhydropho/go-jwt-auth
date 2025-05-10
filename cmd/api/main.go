package main

import (
	"fmt"
	"os"
	"time"

	"github.com/0xirvan/go-jwt-auth/internal/database"
	"github.com/0xirvan/go-jwt-auth/internal/handler"
	"github.com/0xirvan/go-jwt-auth/internal/middleware"
	"github.com/0xirvan/go-jwt-auth/internal/repository"
	"github.com/0xirvan/go-jwt-auth/internal/service"
	"github.com/0xirvan/go-jwt-auth/pkg/logger"
	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Fatalln("Error loading .env file")
	}

	db := database.Connect()

	// periodic cleanup of expired refresh tokens in the background
	go func() {
		ticker := time.NewTicker(24 * time.Hour) // Run every 24 hours
		refreshTokenRepo := repository.NewRefreshTokenRepository(db)
		for range ticker.C {
			if err := refreshTokenRepo.DeleteExpired(); err != nil {
				logger.Log.Error("Error deleting expired refresh tokens: ", err)
			}
		}
	}()

	// in mem store for rate limiting
	store := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Second * 30,
		Limit: 5,
	})

	// Rate limiting middleware
	rateLimit := ratelimit.RateLimiter(store, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	// Init repositories
	userRepo := repository.NewUserRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository(db)

	// Init services
	authService := service.NewAuthService(userRepo, refreshTokenRepo)

	// Init handlers
	authHandler := handler.NewAuthHandler(authService)

	// Routes
	router := gin.Default()

	v1 := router.Group("/api/v1")

	authRoutes := v1.Group("auth")
	{
		authRoutes.POST("/register", rateLimit, authHandler.Register)
		authRoutes.POST("/login", rateLimit, authHandler.Login)
		authRoutes.POST("/refresh", authHandler.RefreshToken)
	}

	protectedAuthRoutes := v1.Group("auth")
	protectedAuthRoutes.Use(middleware.AuthMiddleware(userRepo))
	{
		protectedAuthRoutes.POST("/logout", authHandler.Logout)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Log.Infof("Starting server on port %s", port)
	if err := router.Run(fmt.Sprintf(":%s", port)); err != nil {
		logger.Log.Fatalln("Error starting server: ", err)
	}

}

// utility function to set up rate limiting
func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}
