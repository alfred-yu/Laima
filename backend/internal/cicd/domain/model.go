package domain

import (
	"time"
)

// Pipeline 状态常量
const (
	PipelineStatusPending   = "pending"
	PipelineStatusRunning   = "running"
	PipelineStatusSuccess   = "success"
	PipelineStatusFailed    = "failed"
	PipelineStatusCanceled  = "canceled"
	PipelineStatusSkipped   = "skipped"
)

// Job 状态常量
const (
	JobStatusPending   = "pending"
	JobStatusRunning   = "running"
	JobStatusSuccess   = "success"
	JobStatusFailed    = "failed"
	JobStatusCanceled  = "canceled"
	JobStatusSkipped   = "skipped"
)

// Pipeline 流水线模型
type Pipeline struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	RepositoryID int       `json:"repository_id" gorm:"not null;index"`
	CommitSHA    string    `json:"commit_sha" gorm:"not null;size:40"`
	Ref          string    `json:"ref" gorm:"not null"`
	Status       string    `json:"status" gorm:"not null;default:'pending'"`
	Trigger      string    `json:"trigger" gorm:"not null"`
	YAMLContent  string    `json:"yaml_content" gorm:"type:text"`
	Duration     int       `json:"duration"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// Job 任务模型
type Job struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	PipelineID int       `json:"pipeline_id" gorm:"not null;index"`
	Name       string    `json:"name" gorm:"not null"`
	Status     string    `json:"status" gorm:"not null;default:'pending'"`
	Stage      string    `json:"stage" gorm:"not null"`
	Duration   int       `json:"duration"`
	Log        string    `json:"log" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// PipelineRequest 流水线请求
type PipelineRequest struct {
	RepositoryID int    `json:"repository_id" binding:"required"`
	CommitSHA    string `json:"commit_sha" binding:"required"`
	Ref          string `json:"ref" binding:"required"`
	Trigger      string `json:"trigger" binding:"required"`
}

// PipelineFilter 流水线过滤条件
type PipelineFilter struct {
	RepositoryID int    `json:"repository_id"`
	Status       string `json:"status"`
	Ref          string `json:"ref"`
	CommitSHA    string `json:"commit_sha"`
	Page         int    `json:"page" binding:"min=1"`
	PerPage      int    `json:"per_page" binding:"min=1,max=100"`
}

// PipelineResponse 流水线响应
type PipelineResponse struct {
	ID           int       `json:"id"`
	RepositoryID int       `json:"repository_id"`
	CommitSHA    string    `json:"commit_sha"`
	Ref          string    `json:"ref"`
	Status       string    `json:"status"`
	Trigger      string    `json:"trigger"`
	Duration     int       `json:"duration"`
	Jobs         []*Job    `json:"jobs,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// JobResponse 任务响应
type JobResponse struct {
	ID         int       `json:"id"`
	PipelineID int       `json:"pipeline_id"`
	Name       string    `json:"name"`
	Status     string    `json:"status"`
	Stage      string    `json:"stage"`
	Duration   int       `json:"duration"`
	Log        string    `json:"log"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
