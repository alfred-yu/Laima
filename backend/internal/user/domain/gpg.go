package domain

import "time"

type GPGKey struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	UserID      int       `json:"user_id" gorm:"not null;index"`
	Name        string    `json:"name" gorm:"not null;size:255"`
	Key         string    `json:"-" gorm:"not null;type:text"`
	Fingerprint string    `json:"fingerprint" gorm:"not null;size:128;uniqueIndex"`
	Verified    bool      `json:"verified" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;default:now()"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type GPGKeyRequest struct {
	Name        string    `json:"name" binding:"required"`
	Key         string    `json:"key" binding:"required"`
	Fingerprint string    `json:"fingerprint" binding:"required"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type GPGKeyResponse struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Fingerprint string    `json:"fingerprint"`
	Verified    bool      `json:"verified"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}
