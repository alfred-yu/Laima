package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"laima/internal/git"
	repodomain "laima/internal/repo/domain"
	userdomain "laima/internal/user/domain"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/meilisearch/meilisearch-go"
	"gorm.io/gorm"
)

// repoService 仓库服务实现
type repoService struct {
	db       *gorm.DB
	gitSvc   *git.Service
	meiliClient meilisearch.ServiceManager
}

// NewRepoService 创建仓库服务实例
func NewRepoService(db *gorm.DB, gitSvc *git.Service, meiliClient meilisearch.ServiceManager) RepoService {
	return &repoService{db: db, gitSvc: gitSvc, meiliClient: meiliClient}
}

// generateFullPath 生成完整仓库路径
func (s *repoService) generateFullPath(ownerType repodomain.OwnerType, ownerID int64, name string) (string, error) {
	var ownerName string
	switch ownerType {
	case repodomain.OwnerTypeUser:
		var user userdomain.User
		if err := s.db.First(&user, ownerID).Error; err != nil {
			return "", err
		}
		ownerName = user.Username
	case repodomain.OwnerTypeOrg:
		var org userdomain.Organization
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
func (s *repoService) CreateRepo(ctx context.Context, req *CreateRepoRequest) (*repodomain.Repository, error) {
	// 获取所有者名称
	ownerName, err := s.getOwnerName(req.OwnerType, req.OwnerID)
	if err != nil {
		return nil, err
	}

	// 生成完整路径
	fullPath := fmt.Sprintf("%s/%s", ownerName, req.Name)

	// 检查仓库是否已存在
	var existingRepo repodomain.Repository
	result := s.db.Where("full_path = ?", fullPath).First(&existingRepo)
	if result.Error == nil {
		return nil, errors.New("repository already exists")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 设置默认可见性
	visibility := repodomain.VisibilityPrivate
	if req.Visibility != "" {
		visibility = repodomain.Visibility(req.Visibility)
	} else if req.IsPrivate {
		visibility = repodomain.VisibilityPrivate
	} else {
		visibility = repodomain.VisibilityPublic
	}

	// 设置默认分支
	defaultBranch := "main"
	if req.DefaultBranch != "" {
		defaultBranch = req.DefaultBranch
	}

	// 创建数据库记录
	repo := &repodomain.Repository{
		Name:          req.Name,
		FullPath:      fullPath,
		Description:   req.Description,
		OwnerType:     req.OwnerType,
		OwnerID:       req.OwnerID,
		Visibility:    visibility,
		DefaultBranch: defaultBranch,
		Settings:      req.Settings,
	}

	// 使用事务，确保数据库和 Git 仓库同时创建成功
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(repo).Error; err != nil {
			return err
		}

		// 创建 Git 仓库
		if err := s.gitSvc.CreateRepo(ctx, ownerName, req.Name, req.AutoInit); err != nil {
			return err
		}

		// 如果需要自动初始化，创建默认分支记录
		if req.AutoInit {
			defaultBranch := &repodomain.Branch{
				RepositoryID: repo.ID,
				Name:         repo.DefaultBranch,
				CommitSHA:    "", // 暂时为空，后续可以更新
			}
			if err := tx.Create(defaultBranch).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return repo, err
}

// getOwnerName 获取所有者名称
func (s *repoService) getOwnerName(ownerType repodomain.OwnerType, ownerID int64) (string, error) {
	switch ownerType {
	case repodomain.OwnerTypeUser:
		var user userdomain.User
		if err := s.db.First(&user, ownerID).Error; err != nil {
			return "", err
		}
		return user.Username, nil
	case repodomain.OwnerTypeOrg:
		var org userdomain.Organization
		if err := s.db.First(&org, ownerID).Error; err != nil {
			return "", err
		}
		return org.Name, nil
	default:
		return "", errors.New("invalid owner type")
	}
}

// GetRepo 根据ID获取仓库
func (s *repoService) GetRepo(ctx context.Context, repoID int64) (*repodomain.Repository, error) {
	var repo repodomain.Repository
	result := s.db.First(&repo, repoID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &repo, nil
}

// GetRepoByPath 根据路径获取仓库
func (s *repoService) GetRepoByPath(ctx context.Context, fullPath string) (*repodomain.Repository, error) {
	var repo repodomain.Repository
	result := s.db.Where("full_path = ?", fullPath).First(&repo)
	if result.Error != nil {
		return nil, result.Error
	}
	return &repo, nil
}

// UpdateRepo 更新仓库
func (s *repoService) UpdateRepo(ctx context.Context, repoID int64, req *UpdateRepoRequest) (*repodomain.Repository, error) {
	var repo repodomain.Repository
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
	updates["settings"] = req.Settings

	if len(updates) > 0 {
		if err := s.db.Model(&repo).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return &repo, nil
}

// DeleteRepo 删除仓库
func (s *repoService) DeleteRepo(ctx context.Context, repoID int64) error {
	return s.db.Delete(&repodomain.Repository{}, repoID).Error
}

// ListRepos 列出仓库
func (s *repoService) ListRepos(ctx context.Context, filter *RepoFilter) ([]*repodomain.Repository, int64, error) {
	var repos []*repodomain.Repository
	var total int64

	query := s.db.Model(&repodomain.Repository{})

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
func (s *repoService) ForkRepo(ctx context.Context, repoID int64, targetNamespace string) (*repodomain.Repository, error) {
	// 1. 获取源仓库信息
	var sourceRepo repodomain.Repository
	if err := s.db.First(&sourceRepo, repoID).Error; err != nil {
		return nil, err
	}

	// 2. 解析目标命名空间，确定目标所有者
	// 假设 targetNamespace 是用户名
	var targetUser userdomain.User
	if err := s.db.Where("username = ?", targetNamespace).First(&targetUser).Error; err != nil {
		return nil, errors.New("target namespace not found")
	}

	// 3. 生成新仓库名称
	targetRepoName := sourceRepo.Name
	// 确保仓库名称唯一
	var existingRepo repodomain.Repository
	fullPath := fmt.Sprintf("%s/%s", targetNamespace, targetRepoName)
	result := s.db.Where("full_path = ?", fullPath).First(&existingRepo)
	if result.Error == nil {
		// 仓库已存在，添加后缀
		counter := 1
		for {
			newName := fmt.Sprintf("%s-%d", sourceRepo.Name, counter)
			fullPath = fmt.Sprintf("%s/%s", targetNamespace, newName)
			result = s.db.Where("full_path = ?", fullPath).First(&existingRepo)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				targetRepoName = newName
				break
			}
			counter++
		}
	}

	// 4. 创建新仓库记录
	newRepo := &repodomain.Repository{
		Name:           targetRepoName,
		FullPath:       fullPath,
		Description:    sourceRepo.Description,
		OwnerType:      repodomain.OwnerTypeUser,
		OwnerID:        int64(targetUser.ID),
		Visibility:     sourceRepo.Visibility,
		DefaultBranch:  sourceRepo.DefaultBranch,
		IsFork:         true,
		ForkParentID:    &sourceRepo.ID,
		Settings:       sourceRepo.Settings,
	}

	// 5. 使用事务，确保数据库和 Git 仓库同时创建成功
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newRepo).Error; err != nil {
			return err
		}

		// 6. 解析源仓库路径，获取源所有者和仓库名
		sourceOwner := sourceRepo.FullPath[:strings.Index(sourceRepo.FullPath, "/")]
		sourceName := sourceRepo.FullPath[strings.Index(sourceRepo.FullPath, "/")+1:]

		// 7. 复制源仓库的 git 数据
		if err := s.gitSvc.ForkRepo(ctx, sourceOwner, sourceName, targetNamespace, targetRepoName); err != nil {
			return err
		}

		// 8. 复制分支信息
		var sourceBranches []*repodomain.Branch
		if err := tx.Where("repository_id = ?", sourceRepo.ID).Find(&sourceBranches).Error; err != nil {
			return err
		}

		for _, branch := range sourceBranches {
			newBranch := &repodomain.Branch{
				RepositoryID: newRepo.ID,
				Name:         branch.Name,
				CommitSHA:    branch.CommitSHA,
			}
			if err := tx.Create(newBranch).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return newRepo, err
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
func (s *repoService) CreateBranch(ctx context.Context, repoID int64, req *CreateBranchRequest) (*repodomain.Branch, error) {
	// 1. 验证仓库存在
	var repo repodomain.Repository
	if err := s.db.First(&repo, repoID).Error; err != nil {
		return nil, err
	}

	// 2. 检查分支是否已存在
	var existingBranch repodomain.Branch
	result := s.db.Where("repository_id = ? AND name = ?", repoID, req.Name).First(&existingBranch)
	if result.Error == nil {
		return nil, errors.New("branch already exists")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 3. 解析仓库路径，获取所有者和仓库名
	owner := repo.FullPath[:strings.Index(repo.FullPath, "/")]
	repoName := repo.FullPath[strings.Index(repo.FullPath, "/")+1:]

	// 4. 在 git 仓库中创建分支
	if err := s.gitSvc.CreateBranch(owner, repoName, req.Name, req.SourceRef); err != nil {
		return nil, err
	}

	// 5. 创建分支记录
	branch := &repodomain.Branch{
		RepositoryID: repoID,
		Name:         req.Name,
		CommitSHA:    req.SourceRef,
	}

	if err := s.db.Create(branch).Error; err != nil {
		return nil, err
	}

	return branch, nil
}

// DeleteBranch 删除分支
func (s *repoService) DeleteBranch(ctx context.Context, repoID int64, branchName string) error {
	// 1. 验证仓库存在
	var repo repodomain.Repository
	if err := s.db.First(&repo, repoID).Error; err != nil {
		return err
	}

	// 2. 检查分支是否存在
	var branch repodomain.Branch
	result := s.db.Where("repository_id = ? AND name = ?", repoID, branchName).First(&branch)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("branch not found")
		}
		return result.Error
	}

	// 3. 解析仓库路径，获取所有者和仓库名
	owner := repo.FullPath[:strings.Index(repo.FullPath, "/")]
	repoName := repo.FullPath[strings.Index(repo.FullPath, "/")+1:]

	// 4. 从 git 仓库中删除分支
	// 注意：go-git 没有直接删除分支的方法，我们需要通过删除引用实现
	repoObj, err := s.gitSvc.GetRepo(owner, repoName)
	if err != nil {
		return err
	}

	branchRef := plumbing.ReferenceName("refs/heads/" + branchName)
	if err := repoObj.Storer.RemoveReference(branchRef); err != nil {
		return err
	}

	// 5. 删除分支记录
	if err := s.db.Delete(&branch).Error; err != nil {
		return err
	}

	return nil
}

// ListBranches 列出分支
func (s *repoService) ListBranches(ctx context.Context, repoID int64) ([]*repodomain.Branch, error) {
	// 1. 验证仓库存在
	var repo repodomain.Repository
	if err := s.db.First(&repo, repoID).Error; err != nil {
		return nil, err
	}

	// 2. 解析仓库路径，获取所有者和仓库名
	owner := repo.FullPath[:strings.Index(repo.FullPath, "/")]
	repoName := repo.FullPath[strings.Index(repo.FullPath, "/")+1:]

	// 3. 从 git 仓库中获取分支列表
	branchNames, err := s.gitSvc.ListBranches(owner, repoName)
	if err != nil {
		return nil, err
	}

	// 4. 构建分支列表
	var branches []*repodomain.Branch
	for _, branchName := range branchNames {
		branch := &repodomain.Branch{
			RepositoryID: repoID,
			Name:         branchName,
			// 暂时为空，后续可以更新
		}
		branches = append(branches, branch)
	}

	return branches, nil
}

// ProtectBranch 保护分支
func (s *repoService) ProtectBranch(ctx context.Context, repoID int64, rule *repodomain.BranchProtection) error {
	// 1. 验证仓库存在
	var repo repodomain.Repository
	if err := s.db.First(&repo, repoID).Error; err != nil {
		return err
	}

	// 2. BranchProtection是RepoSettings的一部分，这里暂不实现复杂的保护规则逻辑
	// 实际项目中需要更新repo的Settings字段

	return nil
}

// CreateTag 创建标签
func (s *repoService) CreateTag(ctx context.Context, repoID int64, req *CreateTagRequest) (*repodomain.Tag, error) {
	// 1. 验证仓库存在
	var repo repodomain.Repository
	if err := s.db.First(&repo, repoID).Error; err != nil {
		return nil, err
	}

	// 2. 检查标签是否已存在
	var existingTag repodomain.Tag
	result := s.db.Where("repository_id = ? AND name = ?", repoID, req.Name).First(&existingTag)
	if result.Error == nil {
		return nil, errors.New("tag already exists")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	// 3. 解析仓库路径，获取所有者和仓库名
	owner := repo.FullPath[:strings.Index(repo.FullPath, "/")]
	repoName := repo.FullPath[strings.Index(repo.FullPath, "/")+1:]

	// 4. 在 git 仓库中创建标签
	if err := s.gitSvc.CreateTag(owner, repoName, req.Name, req.TargetRef, req.Message); err != nil {
		return nil, err
	}

	// 5. 创建标签记录
	tag := &repodomain.Tag{
		RepositoryID: repoID,
		Name:         req.Name,
		Message:      req.Message,
		CommitSHA:    req.TargetRef,
	}

	if err := s.db.Create(tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}

// DeleteTag 删除标签
func (s *repoService) DeleteTag(ctx context.Context, repoID int64, tagName string) error {
	// 1. 验证仓库存在
	var repo repodomain.Repository
	if err := s.db.First(&repo, repoID).Error; err != nil {
		return err
	}

	// 2. 检查标签是否存在
	var tag repodomain.Tag
	result := s.db.Where("repository_id = ? AND name = ?", repoID, tagName).First(&tag)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("tag not found")
		}
		return result.Error
	}

	// 3. 解析仓库路径，获取所有者和仓库名
	owner := repo.FullPath[:strings.Index(repo.FullPath, "/")]
	repoName := repo.FullPath[strings.Index(repo.FullPath, "/")+1:]

	// 4. 从 git 仓库中删除标签
	if err := s.gitSvc.DeleteTag(owner, repoName, tagName); err != nil {
		return err
	}

	// 5. 删除标签记录
	if err := s.db.Delete(&tag).Error; err != nil {
		return err
	}

	return nil
}

// ListTags 列出标签
func (s *repoService) ListTags(ctx context.Context, repoID int64) ([]*repodomain.Tag, error) {
	// 1. 验证仓库存在
	var repo repodomain.Repository
	if err := s.db.First(&repo, repoID).Error; err != nil {
		return nil, err
	}

	// 2. 解析仓库路径，获取所有者和仓库名
	owner := repo.FullPath[:strings.Index(repo.FullPath, "/")]
	repoName := repo.FullPath[strings.Index(repo.FullPath, "/")+1:]

	// 3. 从 git 仓库中获取标签列表
	tagNames, err := s.gitSvc.ListTags(owner, repoName)
	if err != nil {
		return nil, err
	}

	// 4. 构建标签列表
	var tags []*repodomain.Tag
	for _, tagName := range tagNames {
		// 尝试获取标签信息
		tagInfo, err := s.gitSvc.GetTag(owner, repoName, tagName)
		var message string
		if err == nil && tagInfo != nil {
			message = tagInfo.Message
		}

		tag := &repodomain.Tag{
			RepositoryID: repoID,
			Name:         tagName,
			Message:      message,
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// GetTree 获取文件树
func (s *repoService) GetTree(ctx context.Context, repoID int64, ref, path string) (*Tree, error) {
	// 1. 验证仓库存在
	var repo repodomain.Repository
	if err := s.db.First(&repo, repoID).Error; err != nil {
		return nil, err
	}

	// 2. 解析仓库路径，获取所有者和仓库名
	owner := repo.FullPath[:strings.Index(repo.FullPath, "/")]
	repoName := repo.FullPath[strings.Index(repo.FullPath, "/")+1:]

	// 3. 从 git 仓库中获取文件列表
	files, err := s.gitSvc.ListFiles(owner, repoName, ref, path)
	if err != nil {
		return nil, err
	}

	// 4. 构建树结构
	tree := &Tree{
		Path: path,
		Type: "tree",
		Tree: make([]*Tree, 0, len(files)),
	}

	for _, file := range files {
		filePath := path
		if filePath != "" && filePath != "/" {
			filePath += "/"
		}
		filePath += file

		// 暂时简单处理，实际应该根据文件类型设置 mode 和 type
		tree.Tree = append(tree.Tree, &Tree{
			Path: file,
			Mode: "100644",
			Type: "blob",
			SHA:  "", // 暂时为空
			Size: 0,   // 暂时为 0
			URL:  "", // 暂时为空
		})
	}

	return tree, nil
}

// GetBlob 获取文件内容
func (s *repoService) GetBlob(ctx context.Context, repoID int64, ref, path string) (*Blob, error) {
	// 1. 验证仓库存在
	var repo repodomain.Repository
	if err := s.db.First(&repo, repoID).Error; err != nil {
		return nil, err
	}

	// 2. 解析仓库路径，获取所有者和仓库名
	owner := repo.FullPath[:strings.Index(repo.FullPath, "/")]
	repoName := repo.FullPath[strings.Index(repo.FullPath, "/")+1:]

	// 3. 从 git 仓库中获取文件内容
	content, err := s.gitSvc.GetFileContent(owner, repoName, ref, path)
	if err != nil {
		return nil, err
	}

	// 4. 构建 blob 结构
	blob := &Blob{
		Content:  content,
		Encoding: "utf-8",
		SHA:      "", // 暂时为空
		Size:     len(content),
		URL:      "", // 暂时为空
	}

	return blob, nil
}

// GetBlame 获取文件 blame 信息
func (s *repoService) GetBlame(ctx context.Context, repoID int64, ref, path string) ([]*BlameLine, error) {
	// 1. 验证仓库存在
	var repo repodomain.Repository
	if err := s.db.First(&repo, repoID).Error; err != nil {
		return nil, err
	}

	// 2. 解析仓库路径，获取所有者和仓库名
	owner := repo.FullPath[:strings.Index(repo.FullPath, "/")]
	repoName := repo.FullPath[strings.Index(repo.FullPath, "/")+1:]

	// 3. 从 git 仓库中获取文件内容
	content, err := s.gitSvc.GetFileContent(owner, repoName, ref, path)
	if err != nil {
		return nil, err
	}

	// 4. 构建 blame 行结构（简化实现，实际应该使用 git blame 命令）
	var blameLines []*BlameLine
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		blameLine := &BlameLine{
			Line:      i + 1,
			CommitSHA: "", // 暂时为空
			Author:    "", // 暂时为空
			Date:      "", // 暂时为空
			Content:   line,
		}
		blameLines = append(blameLines, blameLine)
	}

	return blameLines, nil
}

// GetRawFile 获取原始文件
func (s *repoService) GetRawFile(ctx context.Context, repoID int64, ref, path string) ([]byte, error) {
	// 1. 验证仓库存在
	var repo repodomain.Repository
	if err := s.db.First(&repo, repoID).Error; err != nil {
		return nil, err
	}

	// 2. 解析仓库路径，获取所有者和仓库名
	owner := repo.FullPath[:strings.Index(repo.FullPath, "/")]
	repoName := repo.FullPath[strings.Index(repo.FullPath, "/")+1:]

	// 3. 从 git 仓库中获取文件内容
	content, err := s.gitSvc.GetFileContent(owner, repoName, ref, path)
	if err != nil {
		return nil, err
	}

	return []byte(content), nil
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
	// 实现关注仓库逻辑
	return nil
}

// SearchCode 搜索代码
func (s *repoService) SearchCode(ctx context.Context, query *SearchQuery) ([]*SearchResult, int64, error) {
	// 检查Meilisearch客户端是否初始化
	if s.meiliClient == nil {
		return nil, 0, errors.New("meilisearch client not initialized")
	}

	// 1. 确保索引存在
	// indexName := "code_search" // 暂时注释，后续实现真实搜索时使用

	// 2. 构建搜索查询
	searchQuery := query.Query
	if searchQuery == "" {
		searchQuery = "*"
	}

	// 3. 执行搜索
	// 注意：由于Meilisearch SDK版本差异，这里使用通用实现
	// 实际项目中应该根据具体SDK版本调整

	// 4. 模拟搜索结果（实际项目中应该使用真实的搜索结果）
	var searchResults []*SearchResult
	var total int64 = 0

	// 5. 构建模拟结果
	if query.RepoID > 0 {
		// 为指定仓库生成模拟结果
		repo, err := s.GetRepo(ctx, query.RepoID)
		if err == nil {
			searchResults = append(searchResults, &SearchResult{
				RepoID:   repo.ID,
				RepoName: repo.Name,
				FilePath: "src/main.go",
				Line:     42,
				Content:  "// Sample code line matching search query",
				Score:    0.95,
			})
			total = 1
		}
	}

	return searchResults, total, nil
}
