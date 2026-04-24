package domain

import (
	"time"
)

// Pages 状态常量
const (
	PagesStatusDraft     = "draft"
	PagesStatusPublished = "published"
	PagesStatusArchived  = "archived"
)

// Pages 模型
type Pages struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	RepositoryID int       `json:"repository_id" gorm:"not null;index"`
	Slug         string    `json:"slug" gorm:"not null;uniqueIndex"`
	Title        string    `json:"title" gorm:"not null;size:512"`
	Content      string    `json:"content" gorm:"type:text"`
	Status       string    `json:"status" gorm:"not null;default:'draft'"`
	AuthorID     int       `json:"author_id" gorm:"index"`
	LastEditorID int       `json:"last_editor_id" gorm:"index"`
	PublishAt    time.Time `json:"publish_at" gorm:"index"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// PagesConfig 模型
type PagesConfig struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	RepositoryID int       `json:"repository_id" gorm:"not null;uniqueIndex"`
	CustomDomain string    `json:"custom_domain" gorm:"size:512"`
	Theme        string    `json:"theme" gorm:"size:128;default:'default'"`
	BasePath     string    `json:"base_path" gorm:"size:256;default:'/'"`
	EnableHTTPS  bool      `json:"enable_https" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// CreatePagesRequest 创建Pages请求
type CreatePagesRequest struct {
	RepositoryID int    `json:"repository_id" binding:"required"`
	Slug         string `json:"slug" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Content      string `json:"content"`
	Status       string `json:"status" binding:"omitempty,oneof=draft published archived"`
}

// UpdatePagesRequest 更新Pages请求
type UpdatePagesRequest struct {
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status" binding:"omitempty,oneof=draft published archived"`
}

// PagesFilter 过滤条件
type PagesFilter struct {
	RepositoryID int    `json:"repository_id"`
	Status       string `json:"status"`
	AuthorID     int    `json:"author_id"`
	Search       string `json:"search"`
	Page         int    `json:"page" binding:"min=1"`
	PerPage      int    `json:"per_page" binding:"min=1,max=100"`
}

// UpdatePagesConfigRequest 更新Pages配置请求
type UpdatePagesConfigRequest struct {
	CustomDomain string `json:"custom_domain"`
	Theme        string `json:"theme"`
	BasePath     string `json:"base_path"`
	EnableHTTPS  *bool  `json:"enable_https"`
}
