
package app

import (
	"laima/internal/audit/domain"
	"time"

	"gorm.io/gorm"
)

// AuditLogQuery 查询条件
type AuditLogQuery struct {
	UserID    int    `json:"user_id"`
	Action    string `json:"action"`
	Resource  string `json:"resource"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Page      int    `json:"page"`
	PerPage   int    `json:"per_page"`
}

// AuditService 审计日志服务接口
type AuditService interface {
	LogAudit(log *domain.AuditLog) error
	GetAuditLogs(query *AuditLogQuery) ([]*domain.AuditLog, int64, error)
	GetAuditLogsByUserID(userID int, query *AuditLogQuery) ([]*domain.AuditLog, int64, error)
	GetAuditLogsByResource(resource, resourceID string) ([]*domain.AuditLog, error)
}

// auditService 审计日志服务实现
type auditService struct {
	db *gorm.DB
}

// NewAuditService 创建审计日志服务实例
func NewAuditService(db *gorm.DB) AuditService {
	return &auditService{db: db}
}

// LogAudit 记录审计日志
func (s *auditService) LogAudit(log *domain.AuditLog) error {
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	return s.db.Create(log).Error
}

// GetAuditLogs 获取审计日志列表
func (s *auditService) GetAuditLogs(query *AuditLogQuery) ([]*domain.AuditLog, int64, error) {
	var logs []*domain.AuditLog
	var total int64

	db := s.db.Model(&domain.AuditLog{})

	if query.UserID > 0 {
		db = db.Where("user_id = ?", query.UserID)
	}
	if query.Action != "" {
		db = db.Where("action = ?", query.Action)
	}
	if query.Resource != "" {
		db = db.Where("resource = ?", query.Resource)
	}
	if query.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, query.StartTime)
		if err == nil {
			db = db.Where("created_at >= ?", startTime)
		}
	}
	if query.EndTime != "" {
		endTime, err := time.Parse(time.RFC3339, query.EndTime)
		if err == nil {
			db = db.Where("created_at <= ?", endTime)
		}
	}

	// 计算总数
	db.Count(&total)

	// 分页
	page := query.Page
	if page <= 0 {
		page = 1
	}
	perPage := query.PerPage
	if perPage <= 0 || perPage > 100 {
		perPage = 20
	}

	offset := (page - 1) * perPage

	result := db.Offset(offset).Limit(perPage).Order("created_at DESC").Find(&logs)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return logs, total, nil
}

// GetAuditLogsByUserID 根据用户ID获取审计日志
func (s *auditService) GetAuditLogsByUserID(userID int, query *AuditLogQuery) ([]*domain.AuditLog, int64, error) {
	if query == nil {
		query = &AuditLogQuery{}
	}
	query.UserID = userID
	return s.GetAuditLogs(query)
}

// GetAuditLogsByResource 根据资源获取审计日志
func (s *auditService) GetAuditLogsByResource(resource, resourceID string) ([]*domain.AuditLog, error) {
	var logs []*domain.AuditLog
	db := s.db.Model(&domain.AuditLog{}).Where("resource = ?", resource)
	
	if resourceID != "" {
		db = db.Where("resource_id = ?", resourceID)
	}

	result := db.Order("created_at DESC").Find(&logs)
	if result.Error != nil {
		return nil, result.Error
	}

	return logs, nil
}
