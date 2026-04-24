package api

import (
	"net/http"
	"strconv"

	"laima/internal/security/app"
	"laima/internal/security/domain"

	"github.com/gin-gonic/gin"
)

// SecurityAPI 安全扫描API处理器
type SecurityAPI struct {
	securityService app.SecurityService
}

// NewSecurityAPI 创建安全扫描API处理器
func NewSecurityAPI(securityService app.SecurityService) *SecurityAPI {
	return &SecurityAPI{
		securityService: securityService,
	}
}

// RegisterRoutes 注册路由
func (api *SecurityAPI) RegisterRoutes(r *gin.Engine) {
	security := r.Group("/api/security")

	// 扫描管理
	security.POST("/scans", api.createScan)
	security.GET("/scans", api.listScans)
	security.GET("/scans/:id", api.getScan)
	security.POST("/scans/:id/stop", api.stopScan)

	// 扫描发现
	security.GET("/scans/:id/findings", api.getScanFindings)
	security.GET("/findings/:id", api.getFinding)

	// 扫描配置
	security.GET("/repos/:repo_id/dast/config", api.getDASTConfig)
	security.PUT("/repos/:repo_id/dast/config", api.updateDASTConfig)
	security.GET("/repos/:repo_id/container/config", api.getContainerConfig)
	security.PUT("/repos/:repo_id/container/config", api.updateContainerConfig)

	// 报告
	security.GET("/repos/:repo_id/reports/compliance", api.generateComplianceReport)
}

// createScan 创建安全扫描
func (api *SecurityAPI) createScan(c *gin.Context) {
	var req domain.CreateScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	scan, err := api.securityService.CreateScan(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, scan)
}

// listScans 列出扫描列表
func (api *SecurityAPI) listScans(c *gin.Context) {
	filter := &domain.ScanFilter{
		Page:    1,
		PerPage: 20,
	}

	if repoID := c.Query("repository_id"); repoID != "" {
		if id, err := strconv.Atoi(repoID); err == nil {
			filter.RepositoryID = id
		}
	}

	if scanType := c.Query("scan_type"); scanType != "" {
		filter.ScanType = scanType
	}

	if status := c.Query("status"); status != "" {
		filter.Status = status
	}

	if branch := c.Query("branch"); branch != "" {
		filter.Branch = branch
	}

	if startDate := c.Query("start_date"); startDate != "" {
		filter.StartDate = startDate
	}

	if endDate := c.Query("end_date"); endDate != "" {
		filter.EndDate = endDate
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filter.Page = p
		}
	}

	if perPage := c.Query("per_page"); perPage != "" {
		if pp, err := strconv.Atoi(perPage); err == nil && pp > 0 && pp <= 100 {
			filter.PerPage = pp
		}
	}

	scans, total, err := api.securityService.ListScans(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"scans": scans,
		"total": total,
		"page":  filter.Page,
		"per_page": filter.PerPage,
	})
}

// getScan 获取扫描详情
func (api *SecurityAPI) getScan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scan ID"})
		return
	}

	scan, err := api.securityService.GetScan(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Scan not found"})
		return
	}

	c.JSON(http.StatusOK, scan)
}

// stopScan 停止扫描
func (api *SecurityAPI) stopScan(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scan ID"})
		return
	}

	if err := api.securityService.StopScan(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Scan stopped"})
}

// getScanFindings 获取扫描发现的问题
func (api *SecurityAPI) getScanFindings(c *gin.Context) {
	scanID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scan ID"})
		return
	}

	filter := &domain.FindingFilter{
		ScanID:  scanID,
		Page:    1,
		PerPage: 20,
	}

	if repoID := c.Query("repository_id"); repoID != "" {
		if id, err := strconv.Atoi(repoID); err == nil {
			filter.RepositoryID = id
		}
	}

	if severity := c.Query("severity"); severity != "" {
		filter.Severity = severity
	}

	if cwe := c.Query("cwe"); cwe != "" {
		filter.CWE = cwe
	}

	if filepath := c.Query("filepath"); filepath != "" {
		filter.Filepath = filepath
	}

	if search := c.Query("search"); search != "" {
		filter.Search = search
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filter.Page = p
		}
	}

	if perPage := c.Query("per_page"); perPage != "" {
		if pp, err := strconv.Atoi(perPage); err == nil && pp > 0 && pp <= 100 {
			filter.PerPage = pp
		}
	}

	findings, total, err := api.securityService.GetScanFindings(scanID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"findings": findings,
		"total":    total,
		"page":     filter.Page,
		"per_page": filter.PerPage,
	})
}

// getFinding 获取问题详情
func (api *SecurityAPI) getFinding(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid finding ID"})
		return
	}

	finding, err := api.securityService.GetFinding(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Finding not found"})
		return
	}

	c.JSON(http.StatusOK, finding)
}

// getDASTConfig 获取DAST扫描配置
func (api *SecurityAPI) getDASTConfig(c *gin.Context) {
	repoID, err := strconv.Atoi(c.Param("repo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	config, err := api.securityService.GetDASTConfig(repoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// updateDASTConfig 更新DAST扫描配置
func (api *SecurityAPI) updateDASTConfig(c *gin.Context) {
	repoID, err := strconv.Atoi(c.Param("repo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	var config domain.DASTScanConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.RepositoryID = repoID

	if err := api.securityService.UpdateDASTConfig(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// getContainerConfig 获取容器扫描配置
func (api *SecurityAPI) getContainerConfig(c *gin.Context) {
	repoID, err := strconv.Atoi(c.Param("repo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	config, err := api.securityService.GetContainerConfig(repoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// updateContainerConfig 更新容器扫描配置
func (api *SecurityAPI) updateContainerConfig(c *gin.Context) {
	repoID, err := strconv.Atoi(c.Param("repo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	var config domain.ContainerScanConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.RepositoryID = repoID

	if err := api.securityService.UpdateContainerConfig(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// generateComplianceReport 生成合规报告
func (api *SecurityAPI) generateComplianceReport(c *gin.Context) {
	repoID, err := strconv.Atoi(c.Param("repo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	scanTypes := c.QueryArray("scan_types")

	reportFile, err := api.securityService.GenerateComplianceReport(repoID, scanTypes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"report_url": reportFile,
		"message":    "Compliance report generated",
	})
}
