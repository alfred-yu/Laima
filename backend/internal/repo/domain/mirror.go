package domain

import "time"

type MirrorSyncStatus string

const (
	MirrorSyncStatusPending   MirrorSyncStatus = "pending"
	MirrorSyncStatusRunning   MirrorSyncStatus = "running"
	MirrorSyncStatusSuccess   MirrorSyncStatus = "success"
	MirrorSyncStatusFailed    MirrorSyncStatus = "failed"
	MirrorSyncStatusCanceled  MirrorSyncStatus = "canceled"
)

type MirrorSyncStrategy string

const (
	MirrorSyncStrategyFull      MirrorSyncStrategy = "full"
	MirrorSyncStrategyIncremental MirrorSyncStrategy = "incremental"
)

type Mirror struct {
	ID               int64            `json:"id" gorm:"primaryKey"`
	RepositoryID     int64            `json:"repository_id" gorm:"not null;unique;index"`
	SourceURL        string           `json:"source_url" gorm:"not null;size:512"`
	SourceBranch     string           `json:"source_branch" gorm:"size:255"`
	SourceUsername   string           `json:"source_username" gorm:"size:255"`
	SourcePassword   string           `json:"-" gorm:"size:255"`
	SourceToken      string           `json:"-" gorm:"size:255"`
	SyncStatus       MirrorSyncStatus `json:"sync_status" gorm:"not null;default:'pending'"`
	SyncStrategy     MirrorSyncStrategy `json:"sync_strategy" gorm:"not null;default:'incremental'"`
	SyncInterval     int              `json:"sync_interval" gorm:"default:300"`
	LastSyncAt       time.Time        `json:"last_sync_at"`
	LastSyncDuration int              `json:"last_sync_duration"`
	NextSyncAt       time.Time        `json:"next_sync_at"`
	LastError        string           `json:"last_error" gorm:"type:text"`
	Enabled          bool             `json:"enabled" gorm:"default:true"`
	CreatedAt        time.Time        `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt        time.Time        `json:"updated_at" gorm:"not null;default:now()"`
}

type MirrorConfig struct {
	SourceURL      string           `json:"source_url" binding:"required,url"`
	SourceBranch   string           `json:"source_branch"`
	Username       string           `json:"username"`
	Password       string           `json:"password"`
	Token          string           `json:"token"`
	SyncStrategy   MirrorSyncStrategy `json:"sync_strategy"`
	SyncInterval   int              `json:"sync_interval" binding:"min=60,max=86400"`
	Enabled        bool             `json:"enabled"`
}

type MirrorSyncLog struct {
	ID           int64            `json:"id" gorm:"primaryKey"`
	MirrorID     int64            `json:"mirror_id" gorm:"not null;index"`
	Status       MirrorSyncStatus `json:"status" gorm:"not null"`
	StartTime    time.Time        `json:"start_time" gorm:"not null"`
	EndTime      time.Time        `json:"end_time"`
	Duration     int              `json:"duration"`
	CommitsSynced int             `json:"commits_synced"`
	BranchesSynced int            `json:"branches_synced"`
	TagsSynced   int              `json:"tags_synced"`
	Error        string           `json:"error" gorm:"type:text"`
	CreatedAt    time.Time        `json:"created_at" gorm:"not null;default:now()"`
}

type MirrorCreateRequest struct {
	RepositoryID   int64            `json:"repository_id" binding:"required"`
	SourceURL      string           `json:"source_url" binding:"required,url"`
	SourceBranch   string           `json:"source_branch"`
	Username       string           `json:"username"`
	Password       string           `json:"password"`
	Token          string           `json:"token"`
	SyncStrategy   MirrorSyncStrategy `json:"sync_strategy"`
	SyncInterval   int              `json:"sync_interval" binding:"min=60,max=86400"`
}

type MirrorUpdateRequest struct {
	SourceURL    string             `json:"source_url" binding:"omitempty,url"`
	SourceBranch string             `json:"source_branch"`
	Username     string             `json:"username"`
	Password     string             `json:"password"`
	Token        string             `json:"token"`
	SyncStrategy MirrorSyncStrategy `json:"sync_strategy"`
	SyncInterval int                `json:"sync_interval" binding:"min=60,max=86400"`
	Enabled      bool               `json:"enabled"`
}

type MirrorResponse struct {
	ID               int64              `json:"id"`
	RepositoryID     int64              `json:"repository_id"`
	SourceURL        string             `json:"source_url"`
	SourceBranch     string             `json:"source_branch"`
	SyncStatus       MirrorSyncStatus   `json:"sync_status"`
	SyncStrategy     MirrorSyncStrategy `json:"sync_strategy"`
	SyncInterval     int                `json:"sync_interval"`
	LastSyncAt       time.Time          `json:"last_sync_at"`
	LastSyncDuration int                `json:"last_sync_duration"`
	NextSyncAt       time.Time          `json:"next_sync_at"`
	LastError        string             `json:"last_error"`
	Enabled          bool               `json:"enabled"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

type MirrorSyncResponse struct {
	ID             int64            `json:"id"`
	MirrorID       int64            `json:"mirror_id"`
	Status         MirrorSyncStatus `json:"status"`
	StartTime      time.Time        `json:"start_time"`
	EndTime        time.Time        `json:"end_time"`
	Duration       int              `json:"duration"`
	CommitsSynced  int              `json:"commits_synced"`
	BranchesSynced int              `json:"branches_synced"`
	TagsSynced     int              `json:"tags_synced"`
	Error          string           `json:"error"`
	CreatedAt      time.Time        `json:"created_at"`
}
