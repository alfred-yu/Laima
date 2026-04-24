package app

import (
	"laima/internal/user/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserService_Register(t *testing.T) {
	// 创建内存数据库
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 自动迁移 - 移除 jsonb 支持，使用 TEXT 替代
	err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		avatar_url TEXT,
		bio TEXT,
		settings TEXT DEFAULT '{}',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`).Error
	assert.NoError(t, err)

	// 初始化服务
	service := NewUserService(db)

	// 测试注册
	req := &domain.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := service.Register(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "testuser", resp.User.Username)
	assert.Equal(t, "test@example.com", resp.User.Email)
}

func TestUserService_Login(t *testing.T) {
	// 创建内存数据库
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 自动迁移 - 移除 jsonb 支持，使用 TEXT 替代
	err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		avatar_url TEXT,
		bio TEXT,
		settings TEXT DEFAULT '{}',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`).Error
	assert.NoError(t, err)

	// 初始化服务
	service := NewUserService(db)

	// 先注册用户
	req := &domain.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	_, err = service.Register(req)
	assert.NoError(t, err)

	// 测试登录
	resp, err := service.Login("testuser", "password123")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "testuser", resp.User.Username)
}

func TestUserService_CreateOrganization(t *testing.T) {
	// 创建内存数据库
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 自动迁移 - 移除 jsonb 支持，使用 TEXT 替代
	err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		avatar_url TEXT,
		bio TEXT,
		settings TEXT DEFAULT '{}',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`).Error
	assert.NoError(t, err)

	err = db.Exec(`CREATE TABLE IF NOT EXISTS organizations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		display_name TEXT,
		description TEXT,
		owner_id INTEGER NOT NULL,
		settings TEXT DEFAULT '{}',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`).Error
	assert.NoError(t, err)

	err = db.Exec(`CREATE TABLE IF NOT EXISTS organization_members (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		organization_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		role TEXT NOT NULL DEFAULT 'member',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(organization_id, user_id)
	)`).Error
	assert.NoError(t, err)

	// 初始化服务
	service := NewUserService(db)

	// 先注册用户
	req := &domain.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	resp, err := service.Register(req)
	assert.NoError(t, err)

	// 测试创建组织
	org, err := service.CreateOrganization(resp.User.ID, "testorg", "Test Organization", "Test org description")
	assert.NoError(t, err)
	assert.NotNil(t, org)
	assert.Equal(t, "testorg", org.Name)
	assert.Equal(t, "Test Organization", org.DisplayName)
	assert.Equal(t, resp.User.ID, org.OwnerID)
}
