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

// Runner 状态常量
const (
	RunnerStatusOnline  = "online"
	RunnerStatusOffline = "offline"
	RunnerStatusBusy    = "busy"
)

// Runner 执行器类型常量
const (
	RunnerExecutorShell     = "shell"
	RunnerExecutorDocker    = "docker"
	RunnerExecutorKubernetes = "kubernetes"
)

// Pipeline 流水线模型
type Pipeline struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	RepositoryID int       `json:"repository_id" gorm:"not null;index"`
	PRID         int       `json:"pr_id" gorm:"index"`
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
	PRID         int    `json:"pr_id"`
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

// Runner Runner模型
type Runner struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;unique;size:255"`
	Token       string    `json:"-" gorm:"not null;size:64"`
	Description string    `json:"description" gorm:"size:512"`
	Status      string    `json:"status" gorm:"not null;default:'offline'"`
	Executor    string    `json:"executor" gorm:"not null;default:'shell'"`
	Tags        string    `json:"tags" gorm:"type:text"`
	MaxJobs     int       `json:"max_jobs" gorm:"default:2"`
	ActiveJobs  int       `json:"active_jobs" gorm:"default:0"`
	Version     string    `json:"version" gorm:"size:32"`
	IPAddress   string    `json:"ip_address" gorm:"size:45"`
	LastSeenAt  time.Time `json:"last_seen_at"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// RunnerRegistration Runner注册请求
type RunnerRegistration struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Executor    string `json:"executor" binding:"required"`
	Tags        string `json:"tags"`
	MaxJobs     int    `json:"max_jobs"`
}

// RunnerHeartbeat Runner心跳请求
type RunnerHeartbeat struct {
	RunnerID int    `json:"runner_id" binding:"required"`
	Token    string `json:"token" binding:"required"`
	Version  string `json:"version"`
	Status   string `json:"status"`
}

// RunnerJobRequest Runner拉取任务请求
type RunnerJobRequest struct {
	RunnerID int    `json:"runner_id" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

// RunnerJobUpdate Runner更新任务状态请求
type RunnerJobUpdate struct {
	RunnerID int    `json:"runner_id" binding:"required"`
	Token    string `json:"token" binding:"required"`
	JobID    int    `json:"job_id" binding:"required"`
	Status   string `json:"status" binding:"required"`
	Log      string `json:"log"`
}

// RunnerResponse Runner响应
type RunnerResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Executor    string    `json:"executor"`
	Tags        string    `json:"tags"`
	MaxJobs     int       `json:"max_jobs"`
	ActiveJobs  int       `json:"active_jobs"`
	Version     string    `json:"version"`
	IPAddress   string    `json:"ip_address"`
	LastSeenAt  time.Time `json:"last_seen_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
