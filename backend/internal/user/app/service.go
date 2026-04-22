package app

import (
	"laima/internal/user/domain"

	"gorm.io/gorm"
)

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

	// 仓库成员管理
	AddRepositoryMember(repoID, userID int, role string) (*domain.RepositoryMember, error)
	RemoveRepositoryMember(repoID, userID int) error
	UpdateRepositoryMemberRole(repoID, userID int, role string) (*domain.RepositoryMember, error)
	GetRepositoryMembers(repoID int) ([]*domain.RepositoryMember, error)
}

// userService 用户服务实现
type userService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB) UserService {
	return &userService{db: db}
}

// Login 用户登录
func (s *userService) Login(username, password string) (*domain.AuthResponse, error) {
	// 实现登录逻辑
	// 1. 查找用户
	// 2. 验证密码
	// 3. 生成 JWT token
	// 4. 返回认证响应
	return nil, nil
}

// Register 用户注册
func (s *userService) Register(req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// 实现注册逻辑
	// 1. 检查用户名和邮箱是否已存在
	// 2. 密码哈希
	// 3. 创建用户
	// 4. 生成 JWT token
	// 5. 返回认证响应
	return nil, nil
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
