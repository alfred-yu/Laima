package domain

import (
	"time"
)

// Issue 状态常量
const (
	IssueStatusOpen   = "open"
	IssueStatusClosed = "closed"
)

// Issue 模型
type Issue struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Number       int       `json:"number" gorm:"not null;index"`
	Title        string    `json:"title" gorm:"not null;size:512"`
	Description  string    `json:"description" gorm:"type:text"`
	RepositoryID int       `json:"repository_id" gorm:"not null;index"`
	AuthorID     int       `json:"author_id" gorm:"index"`
	AssigneeID   int       `json:"assignee_id" gorm:"index"`
	State        string    `json:"state" gorm:"not null;default:'open'"`
	MilestoneID  int       `json:"milestone_id" gorm:"index"`
	Labels       string    `json:"labels" gorm:"type:jsonb;default:'[]'"`
	Priority     string    `json:"priority"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// Milestone 里程碑模型
type Milestone struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	RepositoryID int       `json:"repository_id" gorm:"not null;index"`
	Title        string    `json:"title" gorm:"not null"`
	Description  string    `json:"description" gorm:"type:text"`
	State        string    `json:"state" gorm:"not null;default:'open'"`
	DueDate      time.Time `json:"due_date"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// IssueComment Issue评论模型
type IssueComment struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	IssueID      int       `json:"issue_id" gorm:"not null;index"`
	AuthorID     int       `json:"author_id" gorm:"index"`
	Body         string    `json:"body" gorm:"not null;type:text"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// CreateIssueRequest 创建Issue请求
type CreateIssueRequest struct {
	Title        string   `json:"title" binding:"required"`
	Description  string   `json:"description"`
	RepositoryID int      `json:"repository_id" binding:"required"`
	AssigneeID   int      `json:"assignee_id"`
	MilestoneID  int      `json:"milestone_id"`
	Labels       []string `json:"labels"`
	Priority     string   `json:"priority"`
}

// UpdateIssueRequest 更新Issue请求
type UpdateIssueRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	AssigneeID  int      `json:"assignee_id"`
	State       string   `json:"state" binding:"omitempty,oneof=open closed"`
	MilestoneID int      `json:"milestone_id"`
	Labels      []string `json:"labels"`
	Priority    string   `json:"priority"`
}

// IssueFilter Issue过滤条件
type IssueFilter struct {
	RepositoryID int    `json:"repository_id"`
	AuthorID     int    `json:"author_id"`
	AssigneeID   int    `json:"assignee_id"`
	State        string `json:"state"`
	MilestoneID  int    `json:"milestone_id"`
	Labels       string `json:"labels"`
	Priority     string `json:"priority"`
	Search       string `json:"search"`
	Page         int    `json:"page" binding:"min=1"`
	PerPage      int    `json:"per_page" binding:"min=1,max=100"`
}

// MilestoneRequest 里程碑请求
type MilestoneRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	State       string    `json:"state" binding:"omitempty,oneof=open closed"`
}

// IssueCommentRequest Issue评论请求
type IssueCommentRequest struct {
	Body string `json:"body" binding:"required"`
}
