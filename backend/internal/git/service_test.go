package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitService_CreateRepo(t *testing.T) {
	// 创建临时目录作为仓库存储路径
	tempDir, err := os.MkdirTemp("", "laima-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 初始化Git服务
	service := NewService(tempDir)

	// 测试创建仓库
	repoPath, err := service.CreateRepo("test-user", "test-repo")
	assert.NoError(t, err)
	assert.NotEmpty(t, repoPath)

	// 验证仓库目录是否存在
	expectedPath := filepath.Join(tempDir, "test-user", "test-repo")
	assert.DirExists(t, expectedPath)
}

func TestGitService_ListBranches(t *testing.T) {
	// 创建临时目录作为仓库存储路径
	tempDir, err := os.MkdirTemp("", "laima-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 初始化Git服务
	service := NewService(tempDir)

	// 创建仓库
	_, err = service.CreateRepo("test-user", "test-repo")
	assert.NoError(t, err)

	// 测试列出分支
	branches, err := service.ListBranches("test-user", "test-repo")
	assert.NoError(t, err)
	assert.NotEmpty(t, branches)
	assert.Contains(t, branches, "main")
}
