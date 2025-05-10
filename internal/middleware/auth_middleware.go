package middleware

import (
	"net/http"
	"strings"

	"github.com/0xirvan/go-jwt-auth/internal/repository"
	"github.com/0xirvan/go-jwt-auth/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userRepo repository.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			ctx.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header must start with Bearer"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// validate token
		claims, err := jwt.ValidateAccessToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid access token"})
			ctx.Abort()
			return
		}

		user, err := userRepo.FindByID(claims.UserID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Set("userID", user.ID)

		ctx.Next()
	}
}
