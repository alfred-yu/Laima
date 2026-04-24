package app

import (
	"context"
	"laima/internal/git"
	repodomain "laima/internal/repo/domain"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRepoService_CreateRepo(t *testing.T) {
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

	err = db.Exec(`CREATE TABLE IF NOT EXISTS repositories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		full_path TEXT NOT NULL,
		description TEXT,
		owner_type TEXT NOT NULL,
		owner_id INTEGER NOT NULL,
		visibility TEXT NOT NULL DEFAULT 'private',
		default_branch TEXT NOT NULL DEFAULT 'main',
		size INTEGER NOT NULL DEFAULT 0,
		is_fork INTEGER NOT NULL DEFAULT 0,
		fork_parent_id INTEGER,
		is_mirror INTEGER NOT NULL DEFAULT 0,
		mirror_url TEXT,
		settings TEXT DEFAULT '{}',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`).Error
	assert.NoError(t, err)

	// 创建测试用户
	err = db.Exec(`INSERT INTO users (username, email, password_hash, avatar_url, bio, settings) VALUES 
		('test-user', 'test@example.com', 'password_hash', 'avatar_url', 'bio', '{}')`).Error
	assert.NoError(t, err)

	// 创建临时目录作为Git仓库存储路径
	tempDir, err := os.MkdirTemp("", "laima-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 初始化Git服务
	gitService := git.NewService(tempDir)

	// 初始化仓库服务
	service := NewRepoService(db, gitService, nil)

	// 测试创建仓库
	req := &CreateRepoRequest{
		Name:        "test-repo",
		Description: "Test repository",
		OwnerType:   repodomain.OwnerTypeUser,
		OwnerID:     1,
		Visibility:  "public",
		AutoInit:    true,
	}

	repo, err := service.CreateRepo(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, repo)
	assert.Equal(t, "test-repo", repo.Name)
	assert.Equal(t, "Test repository", repo.Description)
	assert.Equal(t, repodomain.OwnerTypeUser, repo.OwnerType)
	assert.Equal(t, int64(1), repo.OwnerID)
	assert.Equal(t, repodomain.VisibilityPublic, repo.Visibility)
}

func TestRepoService_GetRepo(t *testing.T) {
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

	err = db.Exec(`CREATE TABLE IF NOT EXISTS repositories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		full_path TEXT NOT NULL,
		description TEXT,
		owner_type TEXT NOT NULL,
		owner_id INTEGER NOT NULL,
		visibility TEXT NOT NULL DEFAULT 'private',
		default_branch TEXT NOT NULL DEFAULT 'main',
		size INTEGER NOT NULL DEFAULT 0,
		is_fork INTEGER NOT NULL DEFAULT 0,
		fork_parent_id INTEGER,
		is_mirror INTEGER NOT NULL DEFAULT 0,
		mirror_url TEXT,
		settings TEXT DEFAULT '{}',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`).Error
	assert.NoError(t, err)

	// 创建测试用户
	err = db.Exec(`INSERT INTO users (username, email, password_hash, avatar_url, bio, settings) VALUES 
		('test-user', 'test@example.com', 'password_hash', 'avatar_url', 'bio', '{}')`).Error
	assert.NoError(t, err)

	// 创建临时目录作为Git仓库存储路径
	tempDir, err := os.MkdirTemp("", "laima-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 初始化Git服务
	gitService := git.NewService(tempDir)

	// 初始化仓库服务
	service := NewRepoService(db, gitService, nil)

	// 先创建仓库
	req := &CreateRepoRequest{
		Name:        "test-repo",
		Description: "Test repository",
		OwnerType:   repodomain.OwnerTypeUser,
		OwnerID:     1,
		Visibility:  "public",
		AutoInit:    true,
	}

	repo, err := service.CreateRepo(context.Background(), req)
	assert.NoError(t, err)

	// 测试获取仓库
	retrievedRepo, err := service.GetRepo(context.Background(), repo.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrievedRepo)
	assert.Equal(t, repo.ID, retrievedRepo.ID)
	assert.Equal(t, repo.Name, retrievedRepo.Name)
}
