package domain

import (
	"time"
)

// User 用户模型
type User struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"uniqueIndex;not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	AvatarURL    string    `json:"avatar_url"`
	Bio          string    `json:"bio"`
	Settings     string    `json:"settings" gorm:"type:jsonb;default:'{}'"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// Organization 组织模型
type Organization struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	OwnerID     int       `json:"owner_id" gorm:"not null"`
	Settings    string    `json:"settings" gorm:"type:jsonb;default:'{}'"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// OrganizationMember 组织成员模型
type OrganizationMember struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	OrganizationID int       `json:"organization_id" gorm:"not null;index"`
	UserID        int       `json:"user_id" gorm:"not null;index"`
	Role          string    `json:"role" gorm:"not null;default:'member'"`
	CreatedAt     time.Time `json:"created_at" gorm:"not null;default:now()"`
}

// RepositoryMember 仓库成员模型
type RepositoryMember struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	RepositoryID int       `json:"repository_id" gorm:"not null;index"`
	UserID       int       `json:"user_id" gorm:"not null;index"`
	Role         string    `json:"role" gorm:"not null;default:'developer'"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
}

// SSHKey SSH密钥模型
type SSHKey struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	UserID      int       `json:"user_id" gorm:"not null;index"`
	Key         string    `json:"key" gorm:"not null"`
	Fingerprint string    `json:"fingerprint" gorm:"not null;uniqueIndex"`
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;default:now()"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}
