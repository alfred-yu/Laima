package app

import (
	"context"
	"errors"
	"fmt"
	"laima/internal/repo/domain"
	"laima/internal/user/domain"
	"strings"

	"gorm.io/gorm"
)

// repoService 仓库服务实现
type repoService struct {
	db *gorm.DB
}

// NewRepoService 创建仓库服务实例
func NewRepoService(db *gorm.DB) RepoService {
	return &repoService{db: db}
}

// generateFullPath 生成完整仓库路径
func (s *repoService) generateFullPath(ownerType domain.OwnerType, ownerID int64, name string) (string, error) {
	var ownerName string
	switch ownerType {
	case domain.OwnerTypeUser:
		var user domain.User
		if err := s.db.First(&user, ownerID).Error; err != nil {
			return "", err
		}
		ownerName = user.Username
	case domain.OwnerTypeOrg:
		var org domain.Organization
		if err := s.db.First(&org, ownerID).Error; err != nil {
			return "", err
		}
		ownerName = org.Name
	default:
		return "", errors.New("invalid owner type")
	}
	return fmt.Sprintf("%s/%s", ownerName, name), nil
}

// CreateRepo 创建仓库
func (s *repoService) CreateRepo(ctx context.Context, req *CreateRepoRequest) (*domain.Repository, error) {
	// 生成完整路径
	fullPath, err := s.generateFullPath(req.OwnerType, req.OwnerID, req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to generate full path: %v", err)
	}

	// 检查仓库是否已存在
	var existingRepo domain.Repository
	result := s.db.Where("full_path = ?", fullPath).First(&existingRepo)
	if result.Error == nil {
		return nil, errors.New("repository already exists")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 设置默认可见性
	visibility := domain.VisibilityPrivate
	if req.Visibility != "" {
		visibility = domain.Visibility(req.Visibility)
	} else if req.IsPrivate {
		visibility = domain.VisibilityPrivate
	} else {
		visibility = domain.VisibilityPublic
	}

	// 设置默认分支
	defaultBranch := "main"
	if req.DefaultBranch != "" {
		defaultBranch = req.DefaultBranch
	}

	// 创建仓库记录
	repo := &domain.Repository{
		Name:          req.Name,
		FullPath:      fullPath,
		Description:   req.Description,
		OwnerType:     req.OwnerType,
		OwnerID:       req.OwnerID,
		Visibility:    visibility,
		DefaultBranch: defaultBranch,
		Settings:      req.Settings,
	}

	if err := s.db.Create(repo).Error; err != nil {
		return nil, err
	}

	// 如果需要自动初始化，创建默认分支
	if req.AutoInit {
		defaultBranch := &domain.Branch{
			RepositoryID: repo.ID,
			Name:         repo.DefaultBranch,
			CommitSHA:    "0000000000000000000000000000000000000000",
		}
		if err := s.db.Create(defaultBranch).Error; err != nil {
			return nil, err
		}
	}

	return repo, nil
}

// GetRepo 根据ID获取仓库
func (s *repoService) GetRepo(ctx context.Context, repoID int64) (*domain.Repository, error) {
	var repo domain.Repository
	result := s.db.First(&repo, repoID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &repo, nil
}

// GetRepoByPath 根据路径获取仓库
func (s *repoService) GetRepoByPath(ctx context.Context, fullPath string) (*domain.Repository, error) {
	var repo domain.Repository
	result := s.db.Where("full_path = ?", fullPath).First(&repo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &repo, nil
}

// UpdateRepo 更新仓库
func (s *repoService) UpdateRepo(ctx context.Context, repoID int64, req *UpdateRepoRequest) (*domain.Repository, error) {
	var repo domain.Repository
	if err := s.db.First(&repo, repoID).Error; err != nil {
		return nil, err
	}

	// 更新仓库信息
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Visibility != "" {
		updates["visibility"] = req.Visibility
	}
	if req.DefaultBranch != "" {
		updates["default_branch"] = req.DefaultBranch
	}
	if req.Settings != nil {
		updates["settings"] = req.Settings
	}

	if len(updates) > 0 {
		if err := s.db.Model(&repo).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return &repo, nil
}

// DeleteRepo 删除仓库
func (s *repoService) DeleteRepo(ctx context.Context, repoID int64) error {
	return s.db.Delete(&domain.Repository{}, repoID).Error
}

// ListRepos 列出仓库
func (s *repoService) ListRepos(ctx context.Context, filter *RepoFilter) ([]*domain.Repository, int64, error) {
	var repos []*domain.Repository
	var total int64

	query := s.db.Model(&domain.Repository{})

	// 应用过滤条件
	if filter.OwnerID > 0 {
		query = query.Where("owner_id = ?", filter.OwnerID)
	}
	if filter.OwnerType != "" {
		query = query.Where("owner_type = ?", filter.OwnerType)
	}
	if filter.Visibility != "" {
		query = query.Where("visibility = ?", filter.Visibility)
	}
	if filter.Search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	// 计算总数
	query.Count(&total)

	// 分页
	offset := (filter.Page - 1) * filter.PerPage
	result := query.Offset(offset).Limit(filter.PerPage).Find(&repos)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return repos, total, nil
}

// ForkRepo Fork 仓库
func (s *repoService) ForkRepo(ctx context.Context, repoID int64, targetNamespace string) (*domain.Repository, error) {
	// 实现 Fork 逻辑
	// 1. 获取源仓库信息
	// 2. 创建新仓库记录
	// 3. 复制源仓库的 git 数据
	// 4. 返回新仓库信息
	return nil, nil
}

// ImportRepo 导入仓库
func (s *repoService) ImportRepo(ctx context.Context, req *ImportRepoRequest) (*ImportTask, error) {
	// 实现导入逻辑
	// 1. 创建导入任务
	// 2. 异步执行导入操作
	// 3. 返回任务信息
	return nil, nil
}

// CreateBranch 创建分支
func (s *repoService) CreateBranch(ctx context.Context, repoID int64, req *CreateBranchRequest) (*domain.Branch, error) {
	// 实现创建分支逻辑
	// 1. 验证仓库存在
	// 2. 检查分支是否已存在
	// 3. 创建分支记录
	// 4. 在 git 仓库中创建分支
	// 5. 返回分支信息
	return nil, nil
}

// DeleteBranch 删除分支
func (s *repoService) DeleteBranch(ctx context.Context, repoID int64, branch string) error {
	// 实现删除分支逻辑
	// 1. 验证仓库存在
	// 2. 检查分支是否存在
	// 3. 从 git 仓库中删除分支
	// 4. 删除分支记录
	return nil
}

// ListBranches 列出分支
func (s *repoService) ListBranches(ctx context.Context, repoID int64) ([]*domain.Branch, error) {
	var branches []*domain.Branch
	result := s.db.Where("repository_id = ?", repoID).Find(&branches)
	if result.Error != nil {
		return nil, result.Error
	}
	return branches, nil
}

// ProtectBranch 保护分支
func (s *repoService) ProtectBranch(ctx context.Context, repoID int64, rule *domain.BranchProtection) error {
	// 实现分支保护逻辑
	// 1. 验证仓库存在
	// 2. 保存分支保护规则
	return nil
}

// CreateTag 创建标签
func (s *repoService) CreateTag(ctx context.Context, repoID int64, req *CreateTagRequest) (*domain.Tag, error) {
	// 实现创建标签逻辑
	// 1. 验证仓库存在
	// 2. 检查标签是否已存在
	// 3. 创建标签记录
	// 4. 在 git 仓库中创建标签
	// 5. 返回标签信息
	return nil, nil
}

// DeleteTag 删除标签
func (s *repoService) DeleteTag(ctx context.Context, repoID int64, tagName string) error {
	// 实现删除标签逻辑
	// 1. 验证仓库存在
	// 2. 检查标签是否存在
	// 3. 从 git 仓库中删除标签
	// 4. 删除标签记录
	return nil
}

// ListTags 列出标签
func (s *repoService) ListTags(ctx context.Context, repoID int64) ([]*domain.Tag, error) {
	var tags []*domain.Tag
	result := s.db.Where("repository_id = ?", repoID).Find(&tags)
	if result.Error != nil {
		return nil, result.Error
	}
	return tags, nil
}

// GetTree 获取文件树
func (s *repoService) GetTree(ctx context.Context, repoID int64, ref, path string) (*Tree, error) {
	// 实现获取文件树逻辑
	// 1. 验证仓库存在
	// 2. 从 git 仓库中获取文件树
	// 3. 构建树结构
	// 4. 返回文件树
	return nil, nil
}

// GetBlob 获取文件内容
func (s *repoService) GetBlob(ctx context.Context, repoID int64, ref, path string) (*Blob, error) {
	// 实现获取文件内容逻辑
	// 1. 验证仓库存在
	// 2. 从 git 仓库中获取文件内容
	// 3. 构建 blob 结构
	// 4. 返回文件内容
	return nil, nil
}

// GetBlame 获取文件 blame 信息
func (s *repoService) GetBlame(ctx context.Context, repoID int64, ref, path string) ([]*BlameLine, error) {
	// 实现获取 blame 信息逻辑
	// 1. 验证仓库存在
	// 2. 从 git 仓库中获取 blame 信息
	// 3. 构建 blame 行结构
	// 4. 返回 blame 信息
	return nil, nil
}

// GetRawFile 获取原始文件
func (s *repoService) GetRawFile(ctx context.Context, repoID int64, ref, path string) ([]byte, error) {
	// 实现获取原始文件逻辑
	// 1. 验证仓库存在
	// 2. 从 git 仓库中获取原始文件内容
	// 3. 返回文件内容
	return nil, nil
}

// SearchCode 搜索代码
func (s *repoService) SearchCode(ctx context.Context, query *SearchQuery) ([]*SearchResult, int64, error) {
	// 实现代码搜索逻辑
	// 1. 构建搜索查询
	// 2. 执行搜索
	// 3. 处理搜索结果
	// 4. 返回搜索结果
	return nil, 0, nil
}

// StarRepo 星标仓库
func (s *repoService) StarRepo(ctx context.Context, repoID int64) error {
	// 实现星标逻辑
	// 1. 验证仓库存在
	// 2. 检查是否已星标
	// 3. 添加星标记录
	return nil
}

// UnstarRepo 取消星标仓库
func (s *repoService) UnstarRepo(ctx context.Context, repoID int64) error {
	// 实现取消星标逻辑
	// 1. 验证仓库存在
	// 2. 检查是否已星标
	// 3. 删除星标记录
	return nil
}

// WatchRepo 关注仓库
func (s *repoService) WatchRepo(ctx context.Context, repoID int64) error {
	// 实现关注逻辑
	// 1. 验证仓库存在
	// 2. 检查是否已关注
	// 3. 添加关注记录
	return nil
}
