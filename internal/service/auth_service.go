package service

import (
	"fmt"
	"time"

	"github.com/0xirvan/go-jwt-auth/internal/models"
	"github.com/0xirvan/go-jwt-auth/internal/repository"
	"github.com/0xirvan/go-jwt-auth/pkg/jwt"
	"github.com/0xirvan/go-jwt-auth/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

const AccessExpiresIn = 15 * 60 // 15 minutes

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int         `json:"expires_in"`
	User         models.User `json:"user"`
}

type AuthService interface {
	Register(req RegisterRequest) (*AuthResponse, error)
	Login(req LoginRequest) (*AuthResponse, error)
	Refresh(refreshTokenString string) (*AuthResponse, error)
	Logout(userID uint) error
}

type AuthServiceImpl struct {
	userRepo         repository.UserRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewAuthService(userRepo repository.UserRepository, refreshTokenRepo repository.RefreshTokenRepository) AuthService {
	return &AuthServiceImpl{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *AuthServiceImpl) Register(req RegisterRequest) (*AuthResponse, error) {
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.Log.Errorf("failed to create user: %v", err)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	tokenPair, refreshTokenModel, err := jwt.GenerateTokenPair(user)
	if err != nil {
		logger.Log.Errorf("failed to generate token pair: %v", err)
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	// store refresh token in the database
	if err := s.refreshTokenRepo.Create(refreshTokenModel); err != nil {
		logger.Log.Errorf("failed to store refresh token: %v", err)
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    AccessExpiresIn,
		User:         *user,
	}, nil
}

func (s *AuthServiceImpl) Login(req LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	tokenPair, refreshTokenModel, err := jwt.GenerateTokenPair(user)
	if err != nil {
		logger.Log.Errorf("failed to generate token pair: %v", err)
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	// store refresh token in the database
	if err := s.refreshTokenRepo.Create(refreshTokenModel); err != nil {
		logger.Log.Errorf("failed to store refresh token: %v", err)
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    AccessExpiresIn,
		User:         *user,
	}, nil
}

func (s *AuthServiceImpl) Refresh(refreshTokenString string) (*AuthResponse, error) {
	if err := jwt.ValidateRefreshToken(refreshTokenString); err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	refreshToken, err := s.refreshTokenRepo.FindByToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// check if refresh token is expired, if so delete it
	if time.Now().After(refreshToken.ExpiresAt) {
		_ = s.refreshTokenRepo.DeleteByUserID(refreshToken.UserID)
		return nil, fmt.Errorf("refresh token expired")
	}

	tokenPair, newRefreshToken, err := jwt.GenerateTokenPair(&refreshToken.User)
	if err != nil {
		logger.Log.Errorf("failed to generate token pair: %v", err)
		return nil, fmt.Errorf("failed to generate token pair: %w", err)
	}

	// delete old refresh token
	if err := s.refreshTokenRepo.DeleteByUserID(refreshToken.UserID); err != nil {
		logger.Log.Errorf("failed to delete old refresh token: %v", err)
		return nil, fmt.Errorf("failed to delete old refresh token: %w", err)
	}

	// save new refresh token
	if err := s.refreshTokenRepo.Create(newRefreshToken); err != nil {
		logger.Log.Errorf("failed to store new refresh token: %v", err)
		return nil, fmt.Errorf("failed to store new refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: newRefreshToken.Token,
		ExpiresIn:    AccessExpiresIn,
		User:         refreshToken.User,
	}, nil

}

func (s *AuthServiceImpl) Logout(userID uint) error {
	return s.refreshTokenRepo.DeleteByUserID(userID)
}
