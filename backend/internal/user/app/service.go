package app

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"laima/internal/user/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
)

// JWT 密钥（生产环境应放在环境变量中）
var jwtSecret = []byte("laima-secret-key-change-in-production")

// JWT Claims 结构
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// UserService 用户服务接口
type UserService interface {
	// 认证相关
	Login(username, password string) (*domain.AuthResponse, error)
	Register(req *domain.RegisterRequest) (*domain.AuthResponse, error)
	GetUserByID(id int) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	UpdateUser(id int, updates map[string]interface{}) (*domain.User, error)

	// 组织相关
	CreateOrganization(userID int, name, displayName, description string) (*domain.Organization, error)
	GetOrganizationByID(id int) (*domain.Organization, error)
	GetOrganizationByName(name string) (*domain.Organization, error)
	UpdateOrganization(id int, updates map[string]interface{}) (*domain.Organization, error)
	DeleteOrganization(id int) error

	// 组织成员管理
	AddOrganizationMember(orgID, userID int, role string) (*domain.OrganizationMember, error)
	RemoveOrganizationMember(orgID, userID int) error
	UpdateOrganizationMemberRole(orgID, userID int, role string) (*domain.OrganizationMember, error)
	GetOrganizationMembers(orgID int) ([]*domain.OrganizationMember, error)
	GetOrganizationMember(orgID, userID int) (*domain.OrganizationMember, error)

	// 仓库成员管理
	AddRepositoryMember(repoID, userID int, role string) (*domain.RepositoryMember, error)
	RemoveRepositoryMember(repoID, userID int) error
	UpdateRepositoryMemberRole(repoID, userID int, role string) (*domain.RepositoryMember, error)
	GetRepositoryMembers(repoID int) ([]*domain.RepositoryMember, error)
	GetRepositoryMember(repoID, userID int) (*domain.RepositoryMember, error)

	// SSH密钥管理
	AddSSHKey(userID int, key, title string) (*domain.SSHKey, error)
	RemoveSSHKey(id, userID int) error
	GetSSHKeys(userID int) ([]*domain.SSHKey, error)
	GetSSHKeyByFingerprint(fingerprint string) (*domain.SSHKey, error)

	// 权限检查
	CheckOrganizationPermission(orgID, userID int, requiredRole string) (bool, error)
	CheckRepositoryPermission(repoID, userID int, requiredRole string) (bool, error)
	CheckUserPermission(userID, targetUserID int) (bool, error)
	CheckRepositoryPermissionByPerm(repoID, userID int, permission string) (bool, error)
	CheckOrganizationPermissionByPerm(orgID, userID int, permission string) (bool, error)
}

// userService 用户服务实现
type userService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

// hashPassword 密码哈希
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPassword 验证密码
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateToken 生成 JWT token
func generateToken(userID int, username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour * 7) // 7天过期

	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Login 用户登录
func (s *userService) Login(username, password string) (*domain.AuthResponse, error) {
	var user domain.User
	result := s.db.Where("username = ? OR email = ?", username, username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, result.Error
	}

	if !checkPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	token, err := generateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User: domain.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
			Bio:       user.Bio,
		},
	}, nil
}

// Register 用户注册
func (s *userService) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	var count int64
	s.db.Model(&domain.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("username %s already exists", req.Username)
	}

	s.db.Model(&domain.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return nil, fmt.Errorf("email %s already exists", req.Email)
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		AvatarURL:    fmt.Sprintf("https://ui-avatars.com/api/?name=%s&background=random", req.Username),
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	token, err := generateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User: domain.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			AvatarURL: user.AvatarURL,
			Bio:       user.Bio,
		},
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id int) (*domain.User, error) {
	var user domain.User
	result := s.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *userService) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	result := s.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (s *userService) UpdateUser(id int, updates map[string]interface{}) (*domain.User, error) {
	var user domain.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateOrganization 创建组织
func (s *userService) CreateOrganization(userID int, name, displayName, description string) (*domain.Organization, error) {
	org := &domain.Organization{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		OwnerID:     userID,
	}

	if err := s.db.Create(org).Error; err != nil {
		return nil, err
	}

	// 添加创建者为组织成员
	member := &domain.OrganizationMember{
		OrganizationID: org.ID,
		UserID:        userID,
		Role:          "owner",
	}
	if err := s.db.Create(member).Error; err != nil {
		return nil, err
	}

	return org, nil
}

// GetOrganizationByID 根据ID获取组织
func (s *userService) GetOrganizationByID(id int) (*domain.Organization, error) {
	var org domain.Organization
	result := s.db.First(&org, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &org, nil
}

// GetOrganizationByName 根据名称获取组织
func (s *userService) GetOrganizationByName(name string) (*domain.Organization, error) {
	var org domain.Organization
	result := s.db.Where("name = ?", name).First(&org)
	if result.Error != nil {
		return nil, result.Error
	}
	return &org, nil
}

// UpdateOrganization 更新组织信息
func (s *userService) UpdateOrganization(id int, updates map[string]interface{}) (*domain.Organization, error) {
	var org domain.Organization
	if err := s.db.First(&org, id).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&org).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &org, nil
}

// DeleteOrganization 删除组织
func (s *userService) DeleteOrganization(id int) error {
	return s.db.Delete(&domain.Organization{}, id).Error
}

// AddOrganizationMember 添加组织成员
func (s *userService) AddOrganizationMember(orgID, userID int, role string) (*domain.OrganizationMember, error) {
	member := &domain.OrganizationMember{
		OrganizationID: orgID,
		UserID:        userID,
		Role:          role,
	}

	if err := s.db.Create(member).Error; err != nil {
		return nil, err
	}

	return member, nil
}

// RemoveOrganizationMember 移除组织成员
func (s *userService) RemoveOrganizationMember(orgID, userID int) error {
	return s.db.Where("organization_id = ? AND user_id = ?", orgID, userID).Delete(&domain.OrganizationMember{}).Error
}

// UpdateOrganizationMemberRole 更新组织成员角色
func (s *userService) UpdateOrganizationMemberRole(orgID, userID int, role string) (*domain.OrganizationMember, error) {
	var member domain.OrganizationMember
	if err := s.db.Where("organization_id = ? AND user_id = ?", orgID, userID).First(&member).Error; err != nil {
		return nil, err
	}

	member.Role = role
	if err := s.db.Save(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// GetOrganizationMembers 获取组织成员
func (s *userService) GetOrganizationMembers(orgID int) ([]*domain.OrganizationMember, error) {
	var members []*domain.OrganizationMember
	result := s.db.Where("organization_id = ?", orgID).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

// AddRepositoryMember 添加仓库成员
func (s *userService) AddRepositoryMember(repoID, userID int, role string) (*domain.RepositoryMember, error) {
	member := &domain.RepositoryMember{
		RepositoryID: repoID,
		UserID:       userID,
		Role:         role,
	}

	if err := s.db.Create(member).Error; err != nil {
		return nil, err
	}

	return member, nil
}

// RemoveRepositoryMember 移除仓库成员
func (s *userService) RemoveRepositoryMember(repoID, userID int) error {
	return s.db.Where("repository_id = ? AND user_id = ?", repoID, userID).Delete(&domain.RepositoryMember{}).Error
}

// UpdateRepositoryMemberRole 更新仓库成员角色
func (s *userService) UpdateRepositoryMemberRole(repoID, userID int, role string) (*domain.RepositoryMember, error) {
	var member domain.RepositoryMember
	if err := s.db.Where("repository_id = ? AND user_id = ?", repoID, userID).First(&member).Error; err != nil {
		return nil, err
	}

	member.Role = role
	if err := s.db.Save(&member).Error; err != nil {
		return nil, err
	}

	return &member, nil
}

// GetRepositoryMembers 获取仓库成员
func (s *userService) GetRepositoryMembers(repoID int) ([]*domain.RepositoryMember, error) {
	var members []*domain.RepositoryMember
	result := s.db.Where("repository_id = ?", repoID).Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

// GetOrganizationMember 获取组织成员
func (s *userService) GetOrganizationMember(orgID, userID int) (*domain.OrganizationMember, error) {
	var member domain.OrganizationMember
	result := s.db.Where("organization_id = ? AND user_id = ?", orgID, userID).First(&member)
	if result.Error != nil {
		return nil, result.Error
	}
	return &member, nil
}

// GetRepositoryMember 获取仓库成员
func (s *userService) GetRepositoryMember(repoID, userID int) (*domain.RepositoryMember, error) {
	var member domain.RepositoryMember
	result := s.db.Where("repository_id = ? AND user_id = ?", repoID, userID).First(&member)
	if result.Error != nil {
		return nil, result.Error
	}
	return &member, nil
}

// 角色权限等级映射
var rolePermissionMap = map[string]int{
	"member":  1,
	"developer": 2,
	"maintainer": 3,
	"owner":   4,
}

// CheckOrganizationPermission 检查组织权限
func (s *userService) CheckOrganizationPermission(orgID, userID int, requiredRole string) (bool, error) {
	// 获取组织成员信息
	member, err := s.GetOrganizationMember(orgID, userID)
	if err != nil {
		return false, err
	}

	// 检查角色权限
	memberRoleLevel := rolePermissionMap[member.Role]
	requiredRoleLevel := rolePermissionMap[requiredRole]

	return memberRoleLevel >= requiredRoleLevel, nil
}

// CheckRepositoryPermission 检查仓库权限
func (s *userService) CheckRepositoryPermission(repoID, userID int, requiredRole string) (bool, error) {
	// 获取仓库成员信息
	member, err := s.GetRepositoryMember(repoID, userID)
	if err != nil {
		return false, err
	}

	// 检查角色权限
	memberRoleLevel := rolePermissionMap[member.Role]
	requiredRoleLevel := rolePermissionMap[requiredRole]

	return memberRoleLevel >= requiredRoleLevel, nil
}

// CheckUserPermission 检查用户权限
func (s *userService) CheckUserPermission(userID, targetUserID int) (bool, error) {
	// 只能访问自己的信息
	return userID == targetUserID, nil
}

// CheckRepositoryPermissionByPerm 检查仓库操作权限
func (s *userService) CheckRepositoryPermissionByPerm(repoID, userID int, permission string) (bool, error) {
	// 获取仓库成员信息
	var member domain.RepositoryMember
	result := s.db.Where("repository_id = ? AND user_id = ?", repoID, userID).First(&member)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 不是仓库成员，检查是否是仓库所有者或组织成员
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}

	// 检查权限
	return s.checkPermission(member.Role, permission)
}

// CheckOrganizationPermissionByPerm 检查组织操作权限
func (s *userService) CheckOrganizationPermissionByPerm(orgID, userID int, permission string) (bool, error) {
	// 获取组织成员信息
	var member domain.OrganizationMember
	result := s.db.Where("organization_id = ? AND user_id = ?", orgID, userID).First(&member)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 不是组织成员
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}

	// 检查权限
	return s.checkPermission(member.Role, permission)
}

// checkRolePermission 检查角色是否满足所需角色要求
func (s *userService) checkRolePermission(userRole, requiredRole string) (bool, error) {
	userRoleLevel := rolePermissionMap[userRole]
	requiredRoleLevel := rolePermissionMap[requiredRole]

	return userRoleLevel >= requiredRoleLevel, nil
}

// checkPermission 检查角色是否有指定权限
func (s *userService) checkPermission(role, permission string) (bool, error) {
	permissions, ok := domain.RolePermissions[role]
	if !ok {
		return false, fmt.Errorf("unknown role: %s", role)
	}

	for _, p := range permissions {
		if p == permission {
			return true, nil
		}
	}

	return false, nil
}


// calculateSSHKeyFingerprint 计算SSH密钥指纹
func calculateSSHKeyFingerprint(key string) (string, error) {
	// 解析SSH密钥
	parsedKey, err := ssh.ParsePublicKey([]byte(key))
	if err != nil {
		return "", err
	}

	// 计算MD5指纹
	md5Hash := md5.Sum(parsedKey.Marshal())
	fingerprint := base64.StdEncoding.EncodeToString(md5Hash[:])

	return fingerprint, nil
}

// AddSSHKey 添加SSH密钥
func (s *userService) AddSSHKey(userID int, key, title string) (*domain.SSHKey, error) {
	// 解析和验证SSH密钥
	_, err := ssh.ParsePublicKey([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("invalid SSH key: %w", err)
	}

	// 计算密钥指纹
	fingerprint, err := calculateSSHKeyFingerprint(key)
	if err != nil {
		return nil, err
	}

	// 检查指纹是否已存在
	var existingKey domain.SSHKey
	result := s.db.Where("fingerprint = ?", fingerprint).First(&existingKey)
	if result.Error == nil {
		return nil, errors.New("SSH key with this fingerprint already exists")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 创建新的SSH密钥记录
	sshKey := &domain.SSHKey{
		UserID:      userID,
		Key:         key,
		Fingerprint: fingerprint,
		Title:       title,
	}

	if err := s.db.Create(sshKey).Error; err != nil {
		return nil, err
	}

	return sshKey, nil
}

// RemoveSSHKey 删除SSH密钥
func (s *userService) RemoveSSHKey(id, userID int) error {
	// 确保密钥属于该用户
	var sshKey domain.SSHKey
	result := s.db.Where("id = ? AND user_id = ?", id, userID).First(&sshKey)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("SSH key not found")
		}
		return result.Error
	}

	// 删除密钥
	return s.db.Delete(&sshKey).Error
}

// GetSSHKeys 获取用户的所有SSH密钥
func (s *userService) GetSSHKeys(userID int) ([]*domain.SSHKey, error) {
	var keys []*domain.SSHKey
	result := s.db.Where("user_id = ?", userID).Find(&keys)
	if result.Error != nil {
		return nil, result.Error
	}
	return keys, nil
}

// GetSSHKeyByFingerprint 根据指纹获取SSH密钥
func (s *userService) GetSSHKeyByFingerprint(fingerprint string) (*domain.SSHKey, error) {
	var key domain.SSHKey
	result := s.db.Where("fingerprint = ?", fingerprint).First(&key)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("SSH key not found")
		}
		return nil, result.Error
	}
	return &key, nil
}
