
package domain

import (
	"time"
)

// AuditLog 审计日志模型
type AuditLog struct {
	ID         int       `json:"id" gorm:"primaryKey"`
	UserID     int       `json:"user_id" gorm:"index"`
	Action     string    `json:"action" gorm:"index;not null"`
	Resource   string    `json:"resource" gorm:"index;not null"`
	ResourceID string    `json:"resource_id"`
	BeforeData   string    `json:"before_data" gorm:"type:jsonb"`
	AfterData    string    `json:"after_data" gorm:"type:jsonb"`
	IPAddress   string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	Metadata     string    `json:"metadata" gorm:"type:jsonb"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
}
