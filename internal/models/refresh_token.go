package models

import "time"

/*
RefreshToken represents a refresh token for a user.
It contains the token string, the user ID it belongs to,
Use for get new access token without re-login.
*/
type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"uniqueIndex;not null"`
	UserID    uint      `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
}
