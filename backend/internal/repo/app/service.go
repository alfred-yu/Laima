package app

import (
	"context"

	"laima/internal/repo/domain"
)

// RepoService 仓库应用服务接口
type RepoService interface {
	// 仓库 CRUD
	CreateRepo(ctx context.Context, req *CreateRepoRequest) (*domain.Repository, error)
	GetRepo(ctx context.Context, repoID int64) (*domain.Repository, error)
	GetRepoByPath(ctx context.Context, fullPath string) (*domain.Repository, error)
	UpdateRepo(ctx context.Context, repoID int64, req *UpdateRepoRequest) (*domain.Repository, error)
	DeleteRepo(ctx context.Context, repoID int64) error
	ListRepos(ctx context.Context, filter *RepoFilter) ([]*domain.Repository, int64, error)

	// Fork & 导入
	ForkRepo(ctx context.Context, repoID int64, targetNamespace string) (*domain.Repository, error)
	ImportRepo(ctx context.Context, req *ImportRepoRequest) (*ImportTask, error)

	// 分支操作
	CreateBranch(ctx context.Context, repoID int64, req *CreateBranchRequest) (*domain.Branch, error)
	DeleteBranch(ctx context.Context, repoID int64, branch string) error
	ListBranches(ctx context.Context, repoID int64) ([]*domain.Branch, error)
	ProtectBranch(ctx context.Context, repoID int64, rule *domain.BranchProtection) error

	// 标签操作
	CreateTag(ctx context.Context, repoID int64, req *CreateTagRequest) (*domain.Tag, error)
	DeleteTag(ctx context.Context, repoID int64, tagName string) error
	ListTags(ctx context.Context, repoID int64) ([]*domain.Tag, error)

	// 代码浏览
	GetTree(ctx context.Context, repoID int64, ref, path string) (*Tree, error)
	GetBlob(ctx context.Context, repoID int64, ref, path string) (*Blob, error)
	GetBlame(ctx context.Context, repoID int64, ref, path string) ([]*BlameLine, error)
	GetRawFile(ctx context.Context, repoID int64, ref, path string) ([]byte, error)

	// 代码搜索
	SearchCode(ctx context.Context, query *SearchQuery) ([]*SearchResult, int64, error)
	IndexRepoCode(ctx context.Context, repoID int64) error

	// 统计
	StarRepo(ctx context.Context, repoID int64) error
	UnstarRepo(ctx context.Context, repoID int64) error
	WatchRepo(ctx context.Context, repoID int64) error
}

// CreateRepoRequest 创建仓库请求
type CreateRepoRequest struct {
	Name        string            `json:"name" binding:"required"`
	Description string            `json:"description"`
	OwnerType   domain.OwnerType `json:"owner_type" binding:"required,oneof=user org"`
	OwnerID     int64             `json:"owner_id" binding:"required"`
	Visibility  string            `json:"visibility" binding:"omitempty,oneof=public private internal"`
	IsPrivate   bool              `json:"is_private"`
	AutoInit    bool              `json:"auto_init"`
	GitignoreTemplate string      `json:"gitignore_template"`
	LicenseTemplate string        `json:"license_template"`
	DefaultBranch string          `json:"default_branch"`
	Settings     domain.RepoSettings `json:"settings"`
}

// UpdateRepoRequest 更新仓库请求
type UpdateRepoRequest struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Visibility  string            `json:"visibility" binding:"omitempty,oneof=public private internal"`
	DefaultBranch string          `json:"default_branch"`
	Settings     domain.RepoSettings `json:"settings"`
}

// RepoFilter 仓库过滤条件
type RepoFilter struct {
	OwnerID     int64             `json:"owner_id"`
	OwnerType   domain.OwnerType `json:"owner_type"`
	Visibility  string            `json:"visibility"`
	Search      string            `json:"search"`
	Page        int               `json:"page" binding:"min=1"`
	PerPage     int               `json:"per_page" binding:"min=1,max=100"`
}

// CreateBranchRequest 创建分支请求
type CreateBranchRequest struct {
	Name      string `json:"name" binding:"required"`
	SourceRef string `json:"source_ref" binding:"required"` // 来源分支或提交 SHA
}

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name      string `json:"name" binding:"required"`
	TargetRef string `json:"target_ref" binding:"required"` // 目标分支或提交 SHA
	Message   string `json:"message"`
}

// ImportRepoRequest 导入仓库请求
type ImportRepoRequest struct {
	CloneURL    string            `json:"clone_url" binding:"required,url"`
	OwnerType   domain.OwnerType `json:"owner_type" binding:"required,oneof=user org"`
	OwnerID     int64             `json:"owner_id" binding:"required"`
	RepoName    string            `json:"repo_name"`
	Visibility  string            `json:"visibility" binding:"omitempty,oneof=public private internal"`
	AuthToken   string            `json:"auth_token"`
}

// ImportTask 导入任务
type ImportTask struct {
	ID        int64  `json:"id"`
	Status    string `json:"status"`
	Progress  int    `json:"progress"`
	Error     string `json:"error"`
	CreatedAt string `json:"created_at"`
}

// Tree 文件树
type Tree struct {
	Path     string     `json:"path"`
	Mode     string     `json:"mode"`
	Type     string     `json:"type"`
	SHA      string     `json:"sha"`
	Size     int        `json:"size"`
	URL      string     `json:"url"`
	Tree     []*Tree    `json:"tree,omitempty"`
}

// Blob 文件内容
type Blob struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
	SHA      string `json:"sha"`
	Size     int    `json:"size"`
	URL      string `json:"url"`
}

// BlameLine Blame 行信息
type BlameLine struct {
	Line      int    `json:"line"`
	CommitSHA string `json:"commit_sha"`
	Author    string `json:"author"`
	Date      string `json:"date"`
	Content   string `json:"content"`
}

// SearchQuery 搜索查询
type SearchQuery struct {
	Query     string `json:"query" binding:"required"`
	RepoID    int64  `json:"repo_id"`
	Language  string `json:"language"`
	FilePath  string `json:"file_path"`
	Page      int    `json:"page" binding:"min=1"`
	PerPage   int    `json:"per_page" binding:"min=1,max=100"`
}

// SearchResult 搜索结果
type SearchResult struct {
	RepoID    int64  `json:"repo_id"`
	RepoName  string `json:"repo_name"`
	FilePath  string `json:"file_path"`
	Line      int    `json:"line"`
	Content   string `json:"content"`
	Score     float64 `json:"score"`
}