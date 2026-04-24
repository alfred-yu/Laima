package domain

import (
	"time"
)

// OwnerType 所有者类型
type OwnerType string

const (
	OwnerTypeUser OwnerType = "user"
	OwnerTypeOrg  OwnerType = "org"
)

// Visibility 可见性
type Visibility string

const (
	VisibilityPublic   Visibility = "public"
	VisibilityPrivate  Visibility = "private"
	VisibilityInternal Visibility = "internal"
)

// MergeStrategy 合并策略
type MergeStrategy string

const (
	MergeStrategyMerge  MergeStrategy = "merge"
	MergeStrategySquash MergeStrategy = "squash"
	MergeStrategyRebase MergeStrategy = "rebase"
)

// ReviewMode 审查模式
type ReviewMode string

const (
	ReviewModeRelaxed  ReviewMode = "relaxed"
	ReviewModeStandard ReviewMode = "standard"
	ReviewModeStrict   ReviewMode = "strict"
)

// Repository 仓库实体
type Repository struct {
	ID              int64          `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name" gorm:"not null"`
	FullPath        string         `json:"full_path" gorm:"uniqueIndex;not null"`
	Description     string         `json:"description"`
	OwnerType       OwnerType      `json:"owner_type" gorm:"not null"`
	OwnerID         int64          `json:"owner_id" gorm:"not null"`
	Visibility      Visibility     `json:"visibility" gorm:"not null;default:'private'"`
	DefaultBranch   string         `json:"default_branch" gorm:"not null;default:'main'"`
	Size            int64          `json:"size" gorm:"not null;default:0"`
	IsFork          bool           `json:"is_fork" gorm:"not null;default:false"`
	ForkParentID    *int64         `json:"fork_parent_id"`
	IsMirror        bool           `json:"is_mirror" gorm:"not null;default:false"`
	MirrorURL       string         `json:"mirror_url"`
	Settings        RepoSettings   `json:"settings" gorm:"type:jsonb;default:'{}'"`
	CreatedAt       time.Time      `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"not null;default:now()"`
}

// RepoSettings 仓库设置
type RepoSettings struct {
	MergeStrategy      MergeStrategy      `json:"merge_strategy"`
	ReviewMode         ReviewMode         `json:"review_mode"`
	RequirePR          bool               `json:"require_pr"`
	RequireSignedCommit bool              `json:"require_signed"`
	EnableWiki         bool               `json:"enable_wiki"`
	EnableIssues       bool               `json:"enable_issues"`
	EnablePipeline     bool               `json:"enable_pipeline"`
	EnableAIReview     bool               `json:"enable_ai_review"`
	AIReviewRules      []string           `json:"ai_review_rules"`
	BranchProtection   []BranchProtection `json:"branch_protection"`
}

// BranchProtection 分支保护规则
type BranchProtection struct {
	BranchPattern     string   `json:"branch_pattern"`
	RequirePR         bool     `json:"require_pr"`
	RequiredApprovals int      `json:"required_approvals"`
	RequireCI         bool     `json:"require_ci"`
	RequireSigned     bool     `json:"require_signed"`
	AllowedRoles      []string `json:"allowed_roles"`
	CODEOWNERS        bool     `json:"codeowners"`
}

// Branch 分支实体
type Branch struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	RepositoryID int64     `json:"repository_id" gorm:"not null;index"`
	Name         string    `json:"name" gorm:"not null"`
	CommitSHA    string    `json:"commit_sha" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// Tag 标签实体
type Tag struct {
	ID           int64     `json:"id" gorm:"primaryKey"`
	RepositoryID int64     `json:"repository_id" gorm:"not null;index"`
	Name         string    `json:"name" gorm:"not null"`
	CommitSHA    string    `json:"commit_sha" gorm:"not null"`
	Message      string    `json:"message"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
}

// TableName 指定表名
func (Repository) TableName() string {
	return "repositories"
}

// TableName 指定表名
func (Branch) TableName() string {
	return "branches"
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}