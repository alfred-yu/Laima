package git

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"laima/internal/user/domain"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// 测试用的用户服务实现
type testUserService struct {
	db *gorm.DB
}

func (s *testUserService) GetSSHKeyByFingerprint(fingerprint string) (*domain.SSHKey, error) {
	var sshKey domain.SSHKey
	result := s.db.Where("fingerprint = ?", fingerprint).First(&sshKey)
	return &sshKey, result.Error
}

func (s *testUserService) AddSSHKey(userID int, key, title string) (*domain.SSHKey, error) {
	return nil, nil
}

func (s *testUserService) RemoveSSHKey(id, userID int) error {
	return nil
}

func (s *testUserService) GetSSHKeys(userID int) ([]*domain.SSHKey, error) {
	return nil, nil
}

func (s *testUserService) Login(username, password string) (*domain.AuthResponse, error) {
	return nil, nil
}

func (s *testUserService) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	return nil, nil
}

func (s *testUserService) GetUserByID(id int) (*domain.User, error) {
	return nil, nil
}

func (s *testUserService) GetUserByUsername(username string) (*domain.User, error) {
	return nil, nil
}

func (s *testUserService) UpdateUser(id int, updates map[string]interface{}) (*domain.User, error) {
	return nil, nil
}

func (s *testUserService) CreateOrganization(userID int, name, displayName, description string) (*domain.Organization, error) {
	return nil, nil
}

func (s *testUserService) GetOrganizationByID(id int) (*domain.Organization, error) {
	return nil, nil
}

func (s *testUserService) GetOrganizationByName(name string) (*domain.Organization, error) {
	return nil, nil
}

func (s *testUserService) UpdateOrganization(id int, updates map[string]interface{}) (*domain.Organization, error) {
	return nil, nil
}

func (s *testUserService) DeleteOrganization(id int) error {
	return nil
}

func (s *testUserService) AddOrganizationMember(orgID, userID int, role string) (*domain.OrganizationMember, error) {
	return nil, nil
}

func (s *testUserService) RemoveOrganizationMember(orgID, userID int) error {
	return nil
}

func (s *testUserService) UpdateOrganizationMemberRole(orgID, userID int, role string) (*domain.OrganizationMember, error) {
	return nil, nil
}

func (s *testUserService) GetOrganizationMembers(orgID int) ([]*domain.OrganizationMember, error) {
	return nil, nil
}

func (s *testUserService) GetOrganizationMember(orgID, userID int) (*domain.OrganizationMember, error) {
	return nil, nil
}

func (s *testUserService) AddRepositoryMember(repoID, userID int, role string) (*domain.RepositoryMember, error) {
	return nil, nil
}

func (s *testUserService) RemoveRepositoryMember(repoID, userID int) error {
	return nil
}

func (s *testUserService) UpdateRepositoryMemberRole(repoID, userID int, role string) (*domain.RepositoryMember, error) {
	return nil, nil
}

func (s *testUserService) GetRepositoryMembers(repoID int) ([]*domain.RepositoryMember, error) {
	return nil, nil
}

func (s *testUserService) GetRepositoryMember(repoID, userID int) (*domain.RepositoryMember, error) {
	return nil, nil
}

func (s *testUserService) CheckOrganizationPermission(orgID, userID int, requiredRole string) (bool, error) {
	return true, nil
}

func (s *testUserService) CheckRepositoryPermission(repoID, userID int, requiredRole string) (bool, error) {
	return true, nil
}

func (s *testUserService) CheckUserPermission(userID, targetUserID int) (bool, error) {
	return true, nil
}

func TestSSHServer_Start(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "ssh-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 初始化 Git 服务
	gitService := NewService(filepath.Join(tempDir, "repos"))

	// 初始化数据库
	db, err := gorm.Open(sqlite.Open(filepath.Join(tempDir, "test.db")), &gorm.Config{})
	assert.NoError(t, err)
	db.AutoMigrate(&domain.SSHKey{})

	// 初始化用户服务
	userService := &testUserService{db: db}

	// 创建 SSH 服务器
	sshServer := NewSSHServer(":2223", filepath.Join(tempDir, "host_key"), gitService.GetRepoBasePath(), gitService, userService)

	// 启动 SSH 服务器
	ctx, cancel := context.WithCancel(context.Background())
	err = sshServer.Start(ctx)
	assert.NoError(t, err)

	// 等待服务器启动
	time.Sleep(500 * time.Millisecond)

	// 停止服务器
	cancel()
	time.Sleep(500 * time.Millisecond)
}

func TestSSHServer_loadHostKey(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "ssh-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 初始化 Git 服务
	gitService := NewService(filepath.Join(tempDir, "repos"))

	// 初始化用户服务
	userService := &testUserService{}

	// 创建 SSH 服务器
	sshServer := NewSSHServer(":2223", filepath.Join(tempDir, "host_key"), gitService.GetRepoBasePath(), gitService, userService)

	// 加载主机密钥（应该生成新的）
	hostKey, err := sshServer.loadHostKey()
	assert.NoError(t, err)
	assert.NotNil(t, hostKey)

	// 再次加载主机密钥（应该使用已生成的）
	hostKey2, err := sshServer.loadHostKey()
	assert.NoError(t, err)
	assert.NotNil(t, hostKey2)
}

func TestSSHServer_generateHostKey(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "ssh-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 初始化 Git 服务
	gitService := NewService(filepath.Join(tempDir, "repos"))

	// 初始化用户服务
	userService := &testUserService{}

	// 创建 SSH 服务器
	sshServer := NewSSHServer(":2223", filepath.Join(tempDir, "host_key"), gitService.GetRepoBasePath(), gitService, userService)

	// 生成主机密钥
	err = sshServer.generateHostKey()
	assert.NoError(t, err)

	// 检查文件是否存在
	exists, err := fileExists(filepath.Join(tempDir, "host_key"))
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestSSHServer_handleGitCommand(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "ssh-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// 初始化 Git 服务
	gitService := NewService(filepath.Join(tempDir, "repos"))

	// 创建测试仓库
	err = gitService.CreateRepo(context.Background(), "test", "repo", true)
	assert.NoError(t, err)

	// 初始化用户服务
	userService := &testUserService{}

	// 创建 SSH 服务器
	sshServer := NewSSHServer(":2223", filepath.Join(tempDir, "host_key"), gitService.GetRepoBasePath(), gitService, userService)

	// 创建测试通道
	ch := &mockChannel{}

	// 测试 git-receive-pack 命令
	sshServer.handleGitCommand(ch, "git-receive-pack 'test/repo.git'")

	// 测试 git-upload-pack 命令
	sshServer.handleGitCommand(ch, "git-upload-pack 'test/repo.git'")

	// 测试无效的 Git 命令
	sshServer.handleGitCommand(ch, "git-invalid 'test/repo.git'")

	// 测试无效的仓库路径
	sshServer.handleGitCommand(ch, "git-receive-pack 'invalid'")
}

// 辅助函数：检查文件是否存在
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 模拟 SSH 通道
 type mockChannel struct {
	data []byte
}

func (m *mockChannel) Read(data []byte) (int, error) {
	return 0, io.EOF
}

func (m *mockChannel) Write(data []byte) (int, error) {
	m.data = append(m.data, data...)
	return len(data), nil
}

func (m *mockChannel) Close() error {
	return nil
}

func (m *mockChannel) CloseWrite() error {
	return nil
}

func (m *mockChannel) ReadWindow() uint32 {
	return 0
}

func (m *mockChannel) WriteWindow() uint32 {
	return 0
}

func (m *mockChannel) SendRequest(name string, wantReply bool, payload []byte) (bool, error) {
	return false, nil
}

func (m *mockChannel) Stderr() io.ReadWriter {
	return m
}
