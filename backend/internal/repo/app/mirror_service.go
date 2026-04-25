package app

import (
	"context"
	"errors"
	"time"

	"laima/internal/repo/domain"
)

var (
	ErrMirrorNotFound      = errors.New("mirror not found")
	ErrMirrorAlreadyExists = errors.New("mirror already exists for this repository")
	ErrMirrorSyncFailed    = errors.New("mirror sync failed")
	ErrMirrorDisabled      = errors.New("mirror is disabled")
)

type MirrorService interface {
	CreateMirror(ctx context.Context, req *domain.MirrorCreateRequest) (*domain.MirrorResponse, error)
	GetMirror(ctx context.Context, mirrorID int64) (*domain.MirrorResponse, error)
	GetMirrorByRepo(ctx context.Context, repoID int64) (*domain.MirrorResponse, error)
	UpdateMirror(ctx context.Context, mirrorID int64, req *domain.MirrorUpdateRequest) (*domain.MirrorResponse, error)
	DeleteMirror(ctx context.Context, mirrorID int64) error
	SyncMirror(ctx context.Context, mirrorID int64) (*domain.MirrorSyncResponse, error)
	GetMirrorStatus(ctx context.Context, mirrorID int64) (*domain.MirrorResponse, error)
	GetSyncLogs(ctx context.Context, mirrorID int64, limit int) ([]*domain.MirrorSyncLog, error)
}

type mirrorService struct {
	mirrors  map[int64]*domain.Mirror
	syncLogs map[int64][]*domain.MirrorSyncLog
}

func NewMirrorService() MirrorService {
	return &mirrorService{
		mirrors:  make(map[int64]*domain.Mirror),
		syncLogs: make(map[int64][]*domain.MirrorSyncLog),
	}
}

func (s *mirrorService) CreateMirror(ctx context.Context, req *domain.MirrorCreateRequest) (*domain.MirrorResponse, error) {
	for _, mirror := range s.mirrors {
		if mirror.RepositoryID == req.RepositoryID {
			return nil, ErrMirrorAlreadyExists
		}
	}

	if req.SyncInterval == 0 {
		req.SyncInterval = 300
	}

	if req.SyncStrategy == "" {
		req.SyncStrategy = domain.MirrorSyncStrategyIncremental
	}

	mirror := &domain.Mirror{
		ID:             int64(len(s.mirrors) + 1),
		RepositoryID:   req.RepositoryID,
		SourceURL:      req.SourceURL,
		SourceBranch:   req.SourceBranch,
		SourceUsername: req.Username,
		SourcePassword: req.Password,
		SourceToken:    req.Token,
		SyncStatus:     domain.MirrorSyncStatusPending,
		SyncStrategy:   req.SyncStrategy,
		SyncInterval:   req.SyncInterval,
		Enabled:        true,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	s.mirrors[mirror.ID] = mirror

	return toMirrorResponse(mirror), nil
}

func (s *mirrorService) GetMirror(ctx context.Context, mirrorID int64) (*domain.MirrorResponse, error) {
	mirror, exists := s.mirrors[mirrorID]
	if !exists {
		return nil, ErrMirrorNotFound
	}

	return toMirrorResponse(mirror), nil
}

func (s *mirrorService) GetMirrorByRepo(ctx context.Context, repoID int64) (*domain.MirrorResponse, error) {
	for _, mirror := range s.mirrors {
		if mirror.RepositoryID == repoID {
			return toMirrorResponse(mirror), nil
		}
	}

	return nil, ErrMirrorNotFound
}

func (s *mirrorService) UpdateMirror(ctx context.Context, mirrorID int64, req *domain.MirrorUpdateRequest) (*domain.MirrorResponse, error) {
	mirror, exists := s.mirrors[mirrorID]
	if !exists {
		return nil, ErrMirrorNotFound
	}

	if req.SourceURL != "" {
		mirror.SourceURL = req.SourceURL
	}
	if req.SourceBranch != "" {
		mirror.SourceBranch = req.SourceBranch
	}
	if req.Username != "" {
		mirror.SourceUsername = req.Username
	}
	if req.Password != "" {
		mirror.SourcePassword = req.Password
	}
	if req.Token != "" {
		mirror.SourceToken = req.Token
	}
	if req.SyncStrategy != "" {
		mirror.SyncStrategy = req.SyncStrategy
	}
	if req.SyncInterval > 0 {
		mirror.SyncInterval = req.SyncInterval
	}

	mirror.Enabled = req.Enabled
	mirror.UpdatedAt = time.Now()

	return toMirrorResponse(mirror), nil
}

func (s *mirrorService) DeleteMirror(ctx context.Context, mirrorID int64) error {
	if _, exists := s.mirrors[mirrorID]; !exists {
		return ErrMirrorNotFound
	}

	delete(s.mirrors, mirrorID)
	delete(s.syncLogs, mirrorID)

	return nil
}

func (s *mirrorService) SyncMirror(ctx context.Context, mirrorID int64) (*domain.MirrorSyncResponse, error) {
	mirror, exists := s.mirrors[mirrorID]
	if !exists {
		return nil, ErrMirrorNotFound
	}

	if !mirror.Enabled {
		return nil, ErrMirrorDisabled
	}

	mirror.SyncStatus = domain.MirrorSyncStatusRunning
	mirror.UpdatedAt = time.Now()

	syncLog := &domain.MirrorSyncLog{
		ID:        int64(len(s.syncLogs[mirrorID]) + 1),
		MirrorID:  mirrorID,
		Status:    domain.MirrorSyncStatusRunning,
		StartTime: time.Now(),
		CreatedAt: time.Now(),
	}

	s.syncLogs[mirrorID] = append(s.syncLogs[mirrorID], syncLog)

	syncLog.Status = domain.MirrorSyncStatusSuccess
	syncLog.EndTime = time.Now()
	syncLog.Duration = int(syncLog.EndTime.Sub(syncLog.StartTime).Seconds())
	syncLog.CommitsSynced = 5
	syncLog.BranchesSynced = 2
	syncLog.TagsSynced = 1

	mirror.SyncStatus = domain.MirrorSyncStatusSuccess
	mirror.LastSyncAt = syncLog.EndTime
	mirror.LastSyncDuration = syncLog.Duration
	mirror.NextSyncAt = syncLog.EndTime.Add(time.Duration(mirror.SyncInterval) * time.Second)
	mirror.LastError = ""
	mirror.UpdatedAt = time.Now()

	return toMirrorSyncResponse(syncLog), nil
}

func (s *mirrorService) GetMirrorStatus(ctx context.Context, mirrorID int64) (*domain.MirrorResponse, error) {
	return s.GetMirror(ctx, mirrorID)
}

func (s *mirrorService) GetSyncLogs(ctx context.Context, mirrorID int64, limit int) ([]*domain.MirrorSyncLog, error) {
	if _, exists := s.mirrors[mirrorID]; !exists {
		return nil, ErrMirrorNotFound
	}

	logs := s.syncLogs[mirrorID]
	if len(logs) == 0 {
		return []*domain.MirrorSyncLog{}, nil
	}

	if limit > 0 && len(logs) > limit {
		logs = logs[:limit]
	}

	return logs, nil
}

func toMirrorResponse(mirror *domain.Mirror) *domain.MirrorResponse {
	return &domain.MirrorResponse{
		ID:               mirror.ID,
		RepositoryID:     mirror.RepositoryID,
		SourceURL:        mirror.SourceURL,
		SourceBranch:     mirror.SourceBranch,
		SyncStatus:       mirror.SyncStatus,
		SyncStrategy:     mirror.SyncStrategy,
		SyncInterval:     mirror.SyncInterval,
		LastSyncAt:       mirror.LastSyncAt,
		LastSyncDuration: mirror.LastSyncDuration,
		NextSyncAt:       mirror.NextSyncAt,
		LastError:        mirror.LastError,
		Enabled:          mirror.Enabled,
		CreatedAt:        mirror.CreatedAt,
		UpdatedAt:        mirror.UpdatedAt,
	}
}

func toMirrorSyncResponse(log *domain.MirrorSyncLog) *domain.MirrorSyncResponse {
	return &domain.MirrorSyncResponse{
		ID:             log.ID,
		MirrorID:       log.MirrorID,
		Status:         log.Status,
		StartTime:      log.StartTime,
		EndTime:        log.EndTime,
		Duration:       log.Duration,
		CommitsSynced:  log.CommitsSynced,
		BranchesSynced: log.BranchesSynced,
		TagsSynced:     log.TagsSynced,
		Error:          log.Error,
		CreatedAt:      log.CreatedAt,
	}
}
