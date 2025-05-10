package repository

import (
	"time"

	"github.com/0xirvan/go-jwt-auth/internal/models"
	"github.com/0xirvan/go-jwt-auth/pkg/logger"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	FindByToken(token string) (*models.RefreshToken, error)
	DeleteByUserID(userID uint) error
	DeleteExpired() error
	ReplaceToken(userID uint, newToken *models.RefreshToken) error
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
	err := r.DB.Where("token = ?", token).Preload("User").First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *RefreshTokenRepositoryImpl) DeleteByUserID(userID uint) error {
	return r.DB.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

func (r *RefreshTokenRepositoryImpl) DeleteExpired() error {
	return r.DB.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}

func (r *RefreshTokenRepositoryImpl) ReplaceToken(userID uint, newToken *models.RefreshToken) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		logger.Log.Infof("Replacing token for user ID: %d", userID)

		if err := tx.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error; err != nil {
			logger.Log.Errorf("Failed to delete old tokens: %v", err)
			return err
		}

		if err := tx.Create(newToken).Error; err != nil {
			logger.Log.Errorf("Failed to create new token: %v", err)
			return err
		}

		return nil
	})
}
