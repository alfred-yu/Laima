package domain

import (
	"time"
)

// PullRequest PR模型
type PullRequest struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	Number        int       `json:"number" gorm:"not null;index"`
	Title         string    `json:"title" gorm:"not null;size:512"`
	Description   string    `json:"description" gorm:"type:text"`
	RepositoryID  int       `json:"repository_id" gorm:"not null;index"`
	AuthorID      int       `json:"author_id" gorm:"index"`
	SourceRepoID  int       `json:"source_repo_id" gorm:"index"`
	SourceBranch  string    `json:"source_branch" gorm:"not null"`
	TargetBranch  string    `json:"target_branch" gorm:"not null"`
	State         string    `json:"state" gorm:"not null;default:'open'"`
	MergeState    string    `json:"merge_state" gorm:"not null;default:'checking'"`
	ReviewMode    string    `json:"review_mode" gorm:"not null;default:'standard'"`
	HeadCommitSHA string    `json:"head_commit_sha" gorm:"not null"`
	BaseCommitSHA string    `json:"base_commit_sha" gorm:"not null"`
	MergeCommitSHA string   `json:"merge_commit_sha"`
	MergedBy      int       `json:"merged_by" gorm:"index"`
	MergedAt      time.Time `json:"merged_at"`
	ClosedAt      time.Time `json:"closed_at"`
	IsDraft       bool      `json:"is_draft" gorm:"not null;default:false"`
	AI ReviewStatus string   `json:"ai_review_status" gorm:"not null;default:'pending'"`
	CreatedAt     time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// Review 审查模型
type Review struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	PullRequestID int       `json:"pull_request_id" gorm:"not null;index"`
	ReviewerID    int       `json:"reviewer_id" gorm:"index"`
	State         string    `json:"state" gorm:"not null"`
	Score         int       `json:"score"`
	Body          string    `json:"body" gorm:"type:text"`
	SubmittedAt   time.Time `json:"submitted_at" gorm:"not null;default:now()"`
}

// ReviewComment 审查评论模型
type ReviewComment struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	PullRequestID int       `json:"pull_request_id" gorm:"not null;index"`
	ReviewID      int       `json:"review_id" gorm:"index"`
	AuthorID      int       `json:"author_id" gorm:"index"`
	Type          string    `json:"type" gorm:"not null;default:'human'"`
	Path          string    `json:"path" gorm:"not null;size:512"`
	Line          int       `json:"line"`
	DiffHunk      string    `json:"diff_hunk" gorm:"type:text"`
	Body          string    `json:"body" gorm:"not null;type:text"`
	Resolution    string    `json:"resolution"`
	Suggestion    string    `json:"suggestion" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// CreatePRRequest 创建PR请求
type CreatePRRequest struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description"`
	RepositoryID int    `json:"repository_id" binding:"required"`
	SourceRepoID int    `json:"source_repo_id" binding:"required"`
	SourceBranch string `json:"source_branch" binding:"required"`
	TargetBranch string `json:"target_branch" binding:"required"`
	IsDraft      bool   `json:"is_draft"`
}

// UpdatePRRequest 更新PR请求
type UpdatePRRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	TargetBranch string `json:"target_branch"`
	IsDraft     bool   `json:"is_draft"`
}

// ReviewRequest 审查请求
type ReviewRequest struct {
	State string `json:"state" binding:"required,oneof=approve reject comment"`
	Score int    `json:"score"`
	Body  string `json:"body"`
}

// ReviewCommentRequest 审查评论请求
type ReviewCommentRequest struct {
	Path      string `json:"path" binding:"required"`
	Line      int    `json:"line"`
	DiffHunk  string `json:"diff_hunk"`
	Body      string `json:"body" binding:"required"`
	Suggestion string `json:"suggestion"`
}

// PRFilter PR过滤条件
type PRFilter struct {
	RepositoryID int    `json:"repository_id"`
	AuthorID     int    `json:"author_id"`
	State        string `json:"state"`
	Search       string `json:"search"`
	Page         int    `json:"page" binding:"min=1"`
	PerPage      int    `json:"per_page" binding:"min=1,max=100"`
}

// PRResponse PR响应
type PRResponse struct {
	ID            int       `json:"id"`
	Number        int       `json:"number"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	RepositoryID  int       `json:"repository_id"`
	AuthorID      int       `json:"author_id"`
	SourceRepoID  int       `json:"source_repo_id"`
	SourceBranch  string    `json:"source_branch"`
	TargetBranch  string    `json:"target_branch"`
	State         string    `json:"state"`
	MergeState    string    `json:"merge_state"`
	ReviewMode    string    `json:"review_mode"`
	HeadCommitSHA string    `json:"head_commit_sha"`
	BaseCommitSHA string    `json:"base_commit_sha"`
	MergeCommitSHA string   `json:"merge_commit_sha"`
	MergedBy      int       `json:"merged_by"`
	MergedAt      time.Time `json:"merged_at"`
	ClosedAt      time.Time `json:"closed_at"`
	IsDraft       bool      `json:"is_draft"`
	AI ReviewStatus string   `json:"ai_review_status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
