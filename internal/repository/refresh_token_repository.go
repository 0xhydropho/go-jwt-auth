package repository

import (
	"time"

	"github.com/0xirvan/go-jwt-auth/internal/models"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	FindByToken(token string) (*models.RefreshToken, error)
	DeleteByUserID(userID string) error
	DeleteExpired() error
}

type RefreshTokenRepositoryImpl struct {
	DB *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &RefreshTokenRepositoryImpl{
		DB: db,
	}
}

func (r *RefreshTokenRepositoryImpl) Create(token *models.RefreshToken) error {
	return r.DB.Create(token).Error
}

func (r *RefreshTokenRepositoryImpl) FindByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.DB.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *RefreshTokenRepositoryImpl) DeleteByUserID(userID string) error {
	return r.DB.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

func (r *RefreshTokenRepositoryImpl) DeleteExpired() error {
	return r.DB.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}
