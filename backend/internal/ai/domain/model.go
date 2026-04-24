package domain

import (
	"time"
)

// AI Review 状态常量
const (
	AIReviewStatusPending   = "pending"
	AIReviewStatusRunning   = "running"
	AIReviewStatusCompleted = "completed"
	AIReviewStatusFailed    = "failed"
)

// AI Review 严重程度常量
const (
	AIReviewSeverityLow      = "low"
	AIReviewSeverityMedium   = "medium"
	AIReviewSeverityHigh     = "high"
	AIReviewSeverityCritical = "critical"
)

// AIReview AI审查模型
type AIReview struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	PullRequestID int       `json:"pull_request_id" gorm:"not null;index"`
	RepositoryID  int       `json:"repository_id" gorm:"not null;index"`
	Status        string    `json:"status" gorm:"not null;default:'pending'"`
	Score         float64   `json:"score"`
	Summary       string    `json:"summary" gorm:"type:text"`
	StartedAt     time.Time `json:"started_at"`
	CompletedAt   time.Time `json:"completed_at"`
	Error         string    `json:"error" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// AIReviewIssue AI审查问题模型
type AIReviewIssue struct {
	ID            int       `json:"id" gorm:"primaryKey"`
	AIReviewID    int       `json:"ai_review_id" gorm:"not null;index"`
	PullRequestID int       `json:"pull_request_id" gorm:"not null;index"`
	Severity      string    `json:"severity" gorm:"not null"`
	Category      string    `json:"category" gorm:"not null"`
	Title         string    `json:"title" gorm:"not null;size:512"`
	Description   string    `json:"description" gorm:"type:text"`
	Path          string    `json:"path" gorm:"not null;size:512"`
	Line          int       `json:"line"`
	Suggestion    string    `json:"suggestion" gorm:"type:text"`
	Confidence    float64   `json:"confidence"`
	CreatedAt     time.Time `json:"created_at" gorm:"not null;default:now()"`
}

// AIReviewRequest AI审查请求
type AIReviewRequest struct {
	PullRequestID int    `json:"pull_request_id" binding:"required"`
	RepositoryID  int    `json:"repository_id" binding:"required"`
	HeadCommitSHA string `json:"head_commit_sha" binding:"required"`
	BaseCommitSHA string `json:"base_commit_sha" binding:"required"`
}

// AIReviewResponse AI审查响应
type AIReviewResponse struct {
	ID            int             `json:"id"`
	PullRequestID int             `json:"pull_request_id"`
	RepositoryID  int             `json:"repository_id"`
	Status        string          `json:"status"`
	Score         float64         `json:"score"`
	Summary       string          `json:"summary"`
	Issues        []*AIReviewIssue `json:"issues,omitempty"`
	StartedAt     time.Time       `json:"started_at"`
	CompletedAt   time.Time       `json:"completed_at"`
	Error         string          `json:"error"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

// AIReviewFilter AI审查过滤条件
type AIReviewFilter struct {
	PullRequestID int    `json:"pull_request_id"`
	RepositoryID  int    `json:"repository_id"`
	Status        string `json:"status"`
	Page          int    `json:"page" binding:"min=1"`
	PerPage       int    `json:"per_page" binding:"min=1,max=100"`
}
