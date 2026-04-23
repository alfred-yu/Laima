
package api

import (
	"laima/internal/audit/app"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuditAPI 审计日志 API 结构体
type AuditAPI struct {
	auditService app.AuditService
	db *gorm.DB
}

// NewAuditAPI 创建审计日志 API 实例
func NewAuditAPI(db *gorm.DB) *AuditAPI {
	return &AuditAPI{
		auditService: app.NewAuditService(db),
		db: db,
	}
}

// RegisterRoutes 注册路由
func (api *AuditAPI) RegisterRoutes(r *gin.Engine) {
	audit := r.Group("/api/v1/audit")
	{
		audit.GET("/logs", api.GetAuditLogs)
		audit.GET("/logs/user/:user_id", api.GetAuditLogsByUser)
		audit.GET("/logs/resource/:resource", api.GetAuditLogsByResource)
		audit.GET("/logs/resource/:resource/:id", api.GetAuditLogsByResourceID)
	}
}

// GetAuditLogs 获取审计日志列表
func (api *AuditAPI) GetAuditLogs(c *gin.Context) {
	query := &app.AuditLogQuery{
		Action:    c.Query("action"),
		Resource:  c.Query("resource"),
		StartTime: c.Query("start_time"),
		EndTime:   c.Query("end_time"),
		Page:      api.getIntParam(c, "page", 1),
		PerPage:   api.getIntParam(c, "per_page", 20),
	}

	logs, total, err := api.auditService.GetAuditLogs(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  query.Page,
		"per_page": query.PerPage,
	})
}

// GetAuditLogsByUser 根据用户ID获取审计日志
func (api *AuditAPI) GetAuditLogsByUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	query := &app.AuditLogQuery{
		Page:    api.getIntParam(c, "page", 1),
		PerPage: api.getIntParam(c, "per_page", 20),
	}

	logs, total, err := api.auditService.GetAuditLogsByUserID(userID, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  query.Page,
		"per_page": query.PerPage,
	})
}

// GetAuditLogsByResource 根据资源类型获取审计日志
func (api *AuditAPI) GetAuditLogsByResource(c *gin.Context) {
	resource := c.Param("resource")
	
	logs, err := api.auditService.GetAuditLogsByResource(resource, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
	})
}

// GetAuditLogsByResourceID 根据资源类型和ID获取审计日志
func (api *AuditAPI) GetAuditLogsByResourceID(c *gin.Context) {
	resource := c.Param("resource")
	resourceID := c.Param("id")
	
	logs, err := api.auditService.GetAuditLogsByResource(resource, resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": logs,
	})
}

// getIntParam 获取整数参数
func (api *AuditAPI) getIntParam(c *gin.Context, key string, defaultValue int) int {
	val := c.Query(key)
	if val == "" {
		return defaultValue
	}
	
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return i
}
