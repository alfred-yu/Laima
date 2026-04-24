package app

import (
	"context"
	"strconv"
	"time"

	"laima/internal/pages/domain"

	"gorm.io/gorm"
)

// PagesService Pages服务接口
type PagesService interface {
	// Pages CRUD
	CreatePages(ctx context.Context, req *domain.CreatePagesRequest, authorID int) (*domain.Pages, error)
	GetPages(ctx context.Context, pagesID int) (*domain.Pages, error)
	GetPagesBySlug(ctx context.Context, repositoryID int, slug string) (*domain.Pages, error)
	UpdatePages(ctx context.Context, pagesID int, req *domain.UpdatePagesRequest, editorID int) (*domain.Pages, error)
	DeletePages(ctx context.Context, pagesID int) error
	ListPages(ctx context.Context, filter *domain.PagesFilter) ([]*domain.Pages, int64, error)

	// Pages 操作
	PublishPages(ctx context.Context, pagesID int, userID int) (*domain.Pages, error)
	ArchivePages(ctx context.Context, pagesID int, userID int) (*domain.Pages, error)
	UnpublishPages(ctx context.Context, pagesID int, userID int) (*domain.Pages, error)

	// 配置管理
	GetPagesConfig(ctx context.Context, repositoryID int) (*domain.PagesConfig, error)
	UpdatePagesConfig(ctx context.Context, repositoryID int, req *domain.UpdatePagesConfigRequest) (*domain.PagesConfig, error)

	// 静态站点生成
	GenerateStaticSite(ctx context.Context, repositoryID int) (string, error)
	GetStaticSiteURL(ctx context.Context, repositoryID int) (string, error)
}

// PagesServiceImpl Pages服务实现
type PagesServiceImpl struct {
	db *gorm.DB
}

// NewPagesService 创建Pages服务实例
func NewPagesService(db *gorm.DB) PagesService {
	return &PagesServiceImpl{db: db}
}

// CreatePages 创建Pages
func (s *PagesServiceImpl) CreatePages(ctx context.Context, req *domain.CreatePagesRequest, authorID int) (*domain.Pages, error) {
	// 验证slug是否已存在
	var existing domain.Pages
	result := s.db.Where("repository_id = ? AND slug = ?", req.RepositoryID, req.Slug).First(&existing)
	if result.Error == nil {
		return nil, gorm.ErrDuplicatedKey
	}

	// 设置默认状态
	status := req.Status
	if status == "" {
		status = domain.PagesStatusDraft
	}

	// 创建Pages记录
	pages := &domain.Pages{
		RepositoryID: req.RepositoryID,
		Slug:         req.Slug,
		Title:        req.Title,
		Content:      req.Content,
		Status:       status,
		AuthorID:     authorID,
		LastEditorID: authorID,
	}

	// 如果是发布状态，设置发布时间
	if status == domain.PagesStatusPublished {
		pages.PublishAt = time.Now()
	}

	if err := s.db.Create(pages).Error; err != nil {
		return nil, err
	}

	return pages, nil
}

// GetPages 根据ID获取Pages
func (s *PagesServiceImpl) GetPages(ctx context.Context, pagesID int) (*domain.Pages, error) {
	var pages domain.Pages
	result := s.db.First(&pages, pagesID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pages, nil
}

// GetPagesBySlug 根据slug获取Pages
func (s *PagesServiceImpl) GetPagesBySlug(ctx context.Context, repositoryID int, slug string) (*domain.Pages, error) {
	var pages domain.Pages
	result := s.db.Where("repository_id = ? AND slug = ?", repositoryID, slug).First(&pages)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pages, nil
}

// UpdatePages 更新Pages
func (s *PagesServiceImpl) UpdatePages(ctx context.Context, pagesID int, req *domain.UpdatePagesRequest, editorID int) (*domain.Pages, error) {
	var pages domain.Pages
	if err := s.db.First(&pages, pagesID).Error; err != nil {
		return nil, err
	}

	// 更新Pages信息
	updates := make(map[string]interface{})
	if req.Slug != "" && req.Slug != pages.Slug {
		// 验证新slug是否已存在
		var existing domain.Pages
		result := s.db.Where("repository_id = ? AND slug = ? AND id != ?", pages.RepositoryID, req.Slug, pagesID).First(&existing)
		if result.Error == nil {
			return nil, gorm.ErrDuplicatedKey
		}
		updates["slug"] = req.Slug
	}

	if req.Title != "" {
		updates["title"] = req.Title
	}

	if req.Content != "" {
		updates["content"] = req.Content
	}

	if req.Status != "" && req.Status != pages.Status {
		updates["status"] = req.Status
		if req.Status == domain.PagesStatusPublished {
			updates["publish_at"] = time.Now()
		}
	}

	updates["last_editor_id"] = editorID

	if len(updates) > 0 {
		if err := s.db.Model(&pages).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return &pages, nil
}

// DeletePages 删除Pages
func (s *PagesServiceImpl) DeletePages(ctx context.Context, pagesID int) error {
	return s.db.Delete(&domain.Pages{}, pagesID).Error
}

// ListPages 列出Pages
func (s *PagesServiceImpl) ListPages(ctx context.Context, filter *domain.PagesFilter) ([]*domain.Pages, int64, error) {
	var pagesList []*domain.Pages
	var total int64

	query := s.db.Model(&domain.Pages{})

	// 应用过滤条件
	if filter.RepositoryID > 0 {
		query = query.Where("repository_id = ?", filter.RepositoryID)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.AuthorID > 0 {
		query = query.Where("author_id = ?", filter.AuthorID)
	}

	if filter.Search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	// 计算总数
	query.Count(&total)

	// 分页
	offset := (filter.Page - 1) * filter.PerPage
	result := query.Offset(offset).Limit(filter.PerPage).Order("created_at DESC").Find(&pagesList)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return pagesList, total, nil
}

// PublishPages 发布Pages
func (s *PagesServiceImpl) PublishPages(ctx context.Context, pagesID int, userID int) (*domain.Pages, error) {
	var pages domain.Pages
	if err := s.db.First(&pages, pagesID).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"status":       domain.PagesStatusPublished,
		"publish_at":   time.Now(),
		"last_editor_id": userID,
	}

	if err := s.db.Model(&pages).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &pages, nil
}

// ArchivePages 归档Pages
func (s *PagesServiceImpl) ArchivePages(ctx context.Context, pagesID int, userID int) (*domain.Pages, error) {
	var pages domain.Pages
	if err := s.db.First(&pages, pagesID).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"status":       domain.PagesStatusArchived,
		"last_editor_id": userID,
	}

	if err := s.db.Model(&pages).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &pages, nil
}

// UnpublishPages 取消发布Pages
func (s *PagesServiceImpl) UnpublishPages(ctx context.Context, pagesID int, userID int) (*domain.Pages, error) {
	var pages domain.Pages
	if err := s.db.First(&pages, pagesID).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"status":       domain.PagesStatusDraft,
		"last_editor_id": userID,
	}

	if err := s.db.Model(&pages).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &pages, nil
}

// GetPagesConfig 获取Pages配置
func (s *PagesServiceImpl) GetPagesConfig(ctx context.Context, repositoryID int) (*domain.PagesConfig, error) {
	var config domain.PagesConfig
	result := s.db.Where("repository_id = ?", repositoryID).First(&config)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 创建默认配置
			config = domain.PagesConfig{
				RepositoryID: repositoryID,
				Theme:        "default",
				BasePath:     "/",
				EnableHTTPS:  true,
			}
			if err := s.db.Create(&config).Error; err != nil {
				return nil, err
			}
			return &config, nil
		}
		return nil, result.Error
	}
	return &config, nil
}

// UpdatePagesConfig 更新Pages配置
func (s *PagesServiceImpl) UpdatePagesConfig(ctx context.Context, repositoryID int, req *domain.UpdatePagesConfigRequest) (*domain.PagesConfig, error) {
	var config domain.PagesConfig
	result := s.db.Where("repository_id = ?", repositoryID).First(&config)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 创建新配置
			config = domain.PagesConfig{
				RepositoryID: repositoryID,
				Theme:        "default",
				BasePath:     "/",
				EnableHTTPS:  true,
			}
		} else {
			return nil, result.Error
		}
	}

	// 更新配置
	if req.CustomDomain != "" {
		config.CustomDomain = req.CustomDomain
	}

	if req.Theme != "" {
		config.Theme = req.Theme
	}

	if req.BasePath != "" {
		config.BasePath = req.BasePath
	}

	if req.EnableHTTPS != nil {
		config.EnableHTTPS = *req.EnableHTTPS
	}

	if result.Error == gorm.ErrRecordNotFound {
		if err := s.db.Create(&config).Error; err != nil {
			return nil, err
		}
	} else {
		if err := s.db.Save(&config).Error; err != nil {
			return nil, err
		}
	}

	return &config, nil
}

// GenerateStaticSite 生成静态站点
func (s *PagesServiceImpl) GenerateStaticSite(ctx context.Context, repositoryID int) (string, error) {
	// 这里实现静态站点生成逻辑
	// 1. 获取所有已发布的Pages
	// 2. 生成HTML文件
	// 3. 部署到静态文件服务器
	
	// 模拟实现
	return "/path/to/static/site", nil
}

// GetStaticSiteURL 获取静态站点URL
func (s *PagesServiceImpl) GetStaticSiteURL(ctx context.Context, repositoryID int) (string, error) {
	// 获取Pages配置
	config, err := s.GetPagesConfig(ctx, repositoryID)
	if err != nil {
		return "", err
	}

	// 如果有自定义域名，使用自定义域名
	if config.CustomDomain != "" {
		protocol := "https"
		if !config.EnableHTTPS {
			protocol = "http"
		}
		return protocol + "://" + config.CustomDomain + config.BasePath, nil
	}

	// 使用默认域名
	return "https://pages.laima.dev/repo/" + strconv.Itoa(repositoryID), nil
}
