
package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"laima/internal/audit/app"
	"laima/internal/audit/domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuditMiddleware 审计日志中间件
type AuditMiddleware struct {
	auditService app.AuditService
}

// NewAuditMiddleware 创建审计日志中间件
func NewAuditMiddleware(auditService app.AuditService) *AuditMiddleware {
	return &AuditMiddleware{auditService: auditService}
}

// Handler 中间件处理函数
func (m *AuditMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户信息
		var userID int
		userIDVal, exists := c.Get("user_id")
		if exists {
			userID = userIDVal.(int)
		}

		// 获取请求信息
		method := c.Request.Method
		path := c.Request.URL.Path
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 判断是否需要审计
		if m.shouldAudit(method, path) {
			// 获取资源类型
			resource := m.getResourceFromPath(path)
			resourceID := m.getResourceIDFromPath(path)
			action := m.getActionFromMethod(method)

			// 记录审计日志
			log := &domain.AuditLog{
				UserID:    userID,
				Action:    action,
				Resource:  resource,
				ResourceID: resourceID,
				IPAddress:  ip,
				UserAgent:  userAgent,
			}

			// 异步记录日志
			go func() {
				_ = m.auditService.LogAudit(log)
			}()
		}

		c.Next()
	}
}

// shouldAudit 判断是否需要审计
func (m *AuditMiddleware) shouldAudit(method, path string) bool {
	// 只审计写操作
	if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
		return false
	}

	// 忽略一些路径
	if strings.HasPrefix(path, "/health") || 
	   strings.HasPrefix(path, "/metrics") ||
	   strings.HasPrefix(path, "/static") {
		return false
	}

	return true
}

// getResourceFromPath 从路径中获取资源类型
func (m *AuditMiddleware) getResourceFromPath(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "api" && i+1 < len(parts) && i+2 < len(parts) {
			return parts[i+2]
		}
	}
	return "unknown"
}

// getResourceIDFromPath 从路径中获取资源ID
func (m *AuditMiddleware) getResourceIDFromPath(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "api" && i+1 < len(parts) && i+2 < len(parts) && i+3 < len(parts) {
			// 检查是否是数字ID
			if _, err := fmt.Sscanf(parts[i+3], "%d", new(int)); err == nil {
				return parts[i+3]
			}
		}
	}
	return ""
}

// getActionFromMethod 从HTTP方法获取动作
func (m *AuditMiddleware) getActionFromMethod(method string) string {
	switch method {
	case http.MethodPost:
		return "create"
	case http.MethodPut:
		return "update"
	case http.MethodPatch:
		return "update"
	case http.MethodDelete:
		return "delete"
	default:
		return "other"
	}
}

// AuditContextKey 审计上下文键
type AuditContextKey string

const (
	// AuditLogKey 审计日志键
	AuditLogKey AuditContextKey = "audit_log"
)

// GetAuditLogFromContext 从上下文中获取审计日志
func GetAuditLogFromContext(ctx context.Context) *domain.AuditLog {
	log, ok := ctx.Value(AuditLogKey).(*domain.AuditLog)
	if !ok {
		return nil
	}
	return log
}

// SetAuditLogInContext 设置审计日志到上下文
func SetAuditLogInContext(ctx context.Context, log *domain.AuditLog) context.Context {
	return context.WithValue(ctx, AuditLogKey, log)
}

// UpdateAuditLogInContext 更新审计日志
func UpdateAuditLogInContext(ctx context.Context, beforeData, afterData interface{}) {
	log := GetAuditLogFromContext(ctx)
	if log == nil {
		return
	}

	if beforeData != nil {
		if bytes, err := json.Marshal(beforeData); err == nil {
			log.BeforeData = string(bytes)
		}
	}

	if afterData != nil {
		if bytes, err := json.Marshal(afterData); err == nil {
			log.AfterData = string(bytes)
		}
	}
}
