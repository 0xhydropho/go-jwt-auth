package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/0xirvan/go-jwt-auth/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func GenerateAccessToken(user *models.User) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")

	if secretKey == "" {
		return "", fmt.Errorf("jwt secret key is not set")
	}

	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)), // Token expires in 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func GenerateRefreshToken() (string, time.Time, error) {
	refreshSecret := os.Getenv("REFRESH_SECRET")

	if refreshSecret == "" {
		return "", time.Time{}, fmt.Errorf("refresh token secret key is not set")
	}

	expiresAt := time.Now().Add(time.Hour * 24 * 7)

	token := jwt.New(jwt.SigningMethodES256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expiresAt.Unix()
	claims["iat"] = time.Now().Unix()

	tokenString, err := token.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

func ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	secretKey := os.Getenv("JWT_SECRET")

	if secretKey == "" {
		return nil, fmt.Errorf("JWT secret key is not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func ValidateRefreshToken(tokenString string) error {
	refreshSecret := os.Getenv("REFRESH_SECRET")
	if refreshSecret == "" {
		return fmt.Errorf("refresh token secret key is not set")
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(refreshSecret), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func GenerateTokenPair(user *models.User) (*TokenPair, *models.RefreshToken, error) {
	accessTokenString, err := GenerateAccessToken(user)
	if err != nil {
		return nil, nil, err
	}

	refreshTokenString, expiresAt, err := GenerateRefreshToken()
	if err != nil {
		return nil, nil, err
	}

	refreshToken := &models.RefreshToken{
		Token:     refreshTokenString,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, refreshToken, nil
}
