package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"laima/internal/security/domain"

	"gorm.io/gorm"
)

// SecurityServiceImpl 安全扫描服务实现
type SecurityServiceImpl struct {
	db *gorm.DB
}

// NewSecurityService 创建安全扫描服务
func NewSecurityService(db *gorm.DB) SecurityService {
	return &SecurityServiceImpl{
		db: db,
	}
}

// CreateScan 创建安全扫描
func (s *SecurityServiceImpl) CreateScan(req *domain.CreateScanRequest) (*domain.SecurityScan, error) {
	scan := &domain.SecurityScan{
		RepositoryID: req.RepositoryID,
		ScanType:     req.ScanType,
		Status:       domain.ScanStatusPending,
		Branch:       req.Branch,
		Commit:       req.Commit,
	}

	if err := s.db.Create(scan).Error; err != nil {
		return nil, err
	}

	// 根据扫描类型启动扫描
	go func() {
		switch req.ScanType {
		case domain.ScanTypeDAST:
			s.RunDASTScan(scan.ID)
		case domain.ScanTypeContainer:
			s.RunContainerScan(scan.ID)
		}
	}()

	return scan, nil
}

// GetScan 获取扫描详情
func (s *SecurityServiceImpl) GetScan(scanID int) (*domain.SecurityScan, error) {
	var scan domain.SecurityScan
	if err := s.db.First(&scan, scanID).Error; err != nil {
		return nil, err
	}
	return &scan, nil
}

// ListScans 列出扫描列表
func (s *SecurityServiceImpl) ListScans(filter *domain.ScanFilter) ([]*domain.SecurityScan, int64, error) {
	var scans []*domain.SecurityScan
	var total int64

	query := s.db.Model(&domain.SecurityScan{})

	if filter.RepositoryID > 0 {
		query = query.Where("repository_id = ?", filter.RepositoryID)
	}
	if filter.ScanType != "" {
		query = query.Where("scan_type = ?", filter.ScanType)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Branch != "" {
		query = query.Where("branch = ?", filter.Branch)
	}
	if filter.StartDate != "" {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	if filter.EndDate != "" {
		query = query.Where("created_at <= ?", filter.EndDate)
	}

	// 计算总数
	query.Count(&total)

	// 分页
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Offset(offset).Limit(filter.PerPage).Order("created_at DESC").Find(&scans).Error; err != nil {
		return nil, 0, err
	}

	return scans, total, nil
}

// StopScan 停止扫描
func (s *SecurityServiceImpl) StopScan(scanID int) error {
	return s.db.Model(&domain.SecurityScan{}).Where("id = ?", scanID).Update("status", domain.ScanStatusFailed).Error
}

// GetScanFindings 获取扫描发现的问题
func (s *SecurityServiceImpl) GetScanFindings(scanID int, filter *domain.FindingFilter) ([]*domain.ScanFinding, int64, error) {
	var findings []*domain.ScanFinding
	var total int64

	query := s.db.Model(&domain.ScanFinding{}).Where("scan_id = ?", scanID)

	if filter.RepositoryID > 0 {
		query = query.Where("repository_id = ?", filter.RepositoryID)
	}
	if filter.Severity != "" {
		query = query.Where("severity = ?", filter.Severity)
	}
	if filter.CWE != "" {
		query = query.Where("cwe = ?", filter.CWE)
	}
	if filter.Filepath != "" {
		query = query.Where("filepath LIKE ?", "%"+filter.Filepath+"%")
	}
	if filter.Search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	// 计算总数
	query.Count(&total)

	// 分页
	offset := (filter.Page - 1) * filter.PerPage
	if err := query.Offset(offset).Limit(filter.PerPage).Order("severity DESC, created_at DESC").Find(&findings).Error; err != nil {
		return nil, 0, err
	}

	return findings, total, nil
}

// GetFinding 获取问题详情
func (s *SecurityServiceImpl) GetFinding(findingID int) (*domain.ScanFinding, error) {
	var finding domain.ScanFinding
	if err := s.db.First(&finding, findingID).Error; err != nil {
		return nil, err
	}
	return &finding, nil
}

// UpdateDASTConfig 更新DAST扫描配置
func (s *SecurityServiceImpl) UpdateDASTConfig(config *domain.DASTScanConfig) error {
	// 检查是否存在
	var existing domain.DASTScanConfig
	result := s.db.Where("repository_id = ?", config.RepositoryID).First(&existing)

	if result.Error == gorm.ErrRecordNotFound {
		// 创建新配置
		return s.db.Create(config).Error
	} else if result.Error != nil {
		return result.Error
	}

	// 更新现有配置
	config.ID = existing.ID
	return s.db.Save(config).Error
}

// GetDASTConfig 获取DAST扫描配置
func (s *SecurityServiceImpl) GetDASTConfig(repositoryID int) (*domain.DASTScanConfig, error) {
	var config domain.DASTScanConfig
	if err := s.db.Where("repository_id = ?", repositoryID).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回默认配置
			return &domain.DASTScanConfig{
				RepositoryID:   repositoryID,
				ScanDepth:      10,
				Concurrency:    5,
				EnableCrawling: true,
				EnableXSS:      true,
				EnableSQLi:     true,
				EnableCSRF:     true,
			}, nil
		}
		return nil, err
	}
	return &config, nil
}

// UpdateContainerConfig 更新容器扫描配置
func (s *SecurityServiceImpl) UpdateContainerConfig(config *domain.ContainerScanConfig) error {
	// 检查是否存在
	var existing domain.ContainerScanConfig
	result := s.db.Where("repository_id = ?", config.RepositoryID).First(&existing)

	if result.Error == gorm.ErrRecordNotFound {
		// 创建新配置
		return s.db.Create(config).Error
	} else if result.Error != nil {
		return result.Error
	}

	// 更新现有配置
	config.ID = existing.ID
	return s.db.Save(config).Error
}

// GetContainerConfig 获取容器扫描配置
func (s *SecurityServiceImpl) GetContainerConfig(repositoryID int) (*domain.ContainerScanConfig, error) {
	var config domain.ContainerScanConfig
	if err := s.db.Where("repository_id = ?", repositoryID).First(&config).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 返回默认配置
			return &domain.ContainerScanConfig{
				RepositoryID:  repositoryID,
				IncludeOS:     true,
				IncludeApps:   true,
				IncludeSecrets: true,
			}, nil
		}
		return nil, err
	}
	return &config, nil
}

// RunDASTScan 运行DAST扫描
func (s *SecurityServiceImpl) RunDASTScan(scanID int) error {
	// 更新扫描状态为运行中
	s.db.Model(&domain.SecurityScan{}).Where("id = ?", scanID).Updates(map[string]interface{}{
		"status":    domain.ScanStatusRunning,
		"start_time": time.Now(),
	})

	// 获取扫描信息
	var scan domain.SecurityScan
	if err := s.db.First(&scan, scanID).Error; err != nil {
		return err
	}

	// 获取DAST配置
	_, err = s.GetDASTConfig(scan.RepositoryID)
	if err != nil {
		return err
	}

	// 模拟DAST扫描过程
	time.Sleep(30 * time.Second) // 模拟扫描时间

	// 模拟发现的问题
	findings := []domain.ScanFinding{
		{
			ScanID:         scanID,
			RepositoryID:   scan.RepositoryID,
			Severity:       domain.SeverityHigh,
			Title:          "跨站脚本攻击 (XSS) 漏洞",
			Description:    "在用户输入处理中存在XSS漏洞，攻击者可以注入恶意脚本",
			Filepath:       "https://example.com/search",
			LineStart:      1,
			LineEnd:        1,
			CodeSnippet:    "?q=<script>alert('XSS')</script>",
			CWE:            "CWE-79",
			CVSS:           7.4,
			Recommendation: "使用HTML转义处理用户输入",
		},
		{
			ScanID:         scanID,
			RepositoryID:   scan.RepositoryID,
			Severity:       domain.SeverityMedium,
			Title:          "SQL注入漏洞",
			Description:    "在数据库查询中存在SQL注入漏洞",
			Filepath:       "https://example.com/login",
			LineStart:      1,
			LineEnd:        1,
			CodeSnippet:    "SELECT * FROM users WHERE username = 'admin' OR 1=1 --",
			CWE:            "CWE-89",
			CVSS:           6.8,
			Recommendation: "使用参数化查询",
		},
	}

	// 保存发现的问题
	for _, finding := range findings {
		s.db.Create(&finding)
	}

	// 更新扫描状态为完成
	endTime := time.Now()
	duration := int(endTime.Sub(scan.StartTime).Seconds())
	s.db.Model(&domain.SecurityScan{}).Where("id = ?", scanID).Updates(map[string]interface{}{
		"status":     domain.ScanStatusCompleted,
		"end_time":   endTime,
		"duration":   duration,
		"findings":   len(findings),
		"critical":   0,
		"high":       1,
		"medium":     1,
		"low":        0,
		"report_url": fmt.Sprintf("/api/security/scans/%d/report", scanID),
	})

	return nil
}

// RunContainerScan 运行容器扫描
func (s *SecurityServiceImpl) RunContainerScan(scanID int) error {
	// 更新扫描状态为运行中
	s.db.Model(&domain.SecurityScan{}).Where("id = ?", scanID).Updates(map[string]interface{}{
		"status":    domain.ScanStatusRunning,
		"start_time": time.Now(),
	})

	// 获取扫描信息
	var scan domain.SecurityScan
	if err := s.db.First(&scan, scanID).Error; err != nil {
		return err
	}

	// 获取容器扫描配置
	config, err := s.GetContainerConfig(scan.RepositoryID)
	if err != nil {
		return err
	}

	// 模拟容器扫描过程
	time.Sleep(20 * time.Second) // 模拟扫描时间

	// 模拟发现的问题
	findings := []domain.ScanFinding{
		{
			ScanID:         scanID,
			RepositoryID:   scan.RepositoryID,
			Severity:       domain.SeverityCritical,
			Title:          "严重的OS漏洞",
			Description:    "容器基础镜像中存在严重的OS漏洞",
			Filepath:       config.ImageName + ":" + config.ImageTag,
			LineStart:      1,
			LineEnd:        1,
			CodeSnippet:    "Ubuntu 20.04 LTS - CVE-2023-1234",
			CWE:            "CWE-200",
			CVSS:           9.8,
			Recommendation: "更新基础镜像到最新版本",
		},
		{
			ScanID:         scanID,
			RepositoryID:   scan.RepositoryID,
			Severity:       domain.SeverityMedium,
			Title:          "不安全的依赖包",
			Description:    "容器中存在有漏洞的依赖包",
			Filepath:       config.ImageName + ":" + config.ImageTag,
			LineStart:      1,
			LineEnd:        1,
			CodeSnippet:    "nodejs@14.17.0 - CVE-2023-5678",
			CWE:            "CWE-400",
			CVSS:           6.5,
			Recommendation: "更新依赖包到安全版本",
		},
	}

	// 保存发现的问题
	for _, finding := range findings {
		s.db.Create(&finding)
	}

	// 更新扫描状态为完成
	endTime := time.Now()
	duration := int(endTime.Sub(scan.StartTime).Seconds())
	s.db.Model(&domain.SecurityScan{}).Where("id = ?", scanID).Updates(map[string]interface{}{
		"status":     domain.ScanStatusCompleted,
		"end_time":   endTime,
		"duration":   duration,
		"findings":   len(findings),
		"critical":   1,
		"high":       0,
		"medium":     1,
		"low":        0,
		"report_url": fmt.Sprintf("/api/security/scans/%d/report", scanID),
	})

	return nil
}

// GenerateComplianceReport 生成合规报告
func (s *SecurityServiceImpl) GenerateComplianceReport(repositoryID int, scanTypes []string) (string, error) {
	// 生成报告目录
	reportDir := filepath.Join("/tmp", "laima-reports")
	os.MkdirAll(reportDir, 0755)

	// 生成报告文件名
	reportFile := filepath.Join(reportDir, fmt.Sprintf("compliance-report-%d-%s.json", repositoryID, time.Now().Format("20060102-150405")))

	// 收集扫描结果
	var scans []domain.SecurityScan
	query := s.db.Where("repository_id = ?", repositoryID)
	if len(scanTypes) > 0 {
		query = query.Where("scan_type IN ?", scanTypes)
	}
	query.Find(&scans)

	// 生成报告数据
	report := map[string]interface{}{
		"repository_id": repositoryID,
		"generated_at":  time.Now(),
		"scan_types":    scanTypes,
		"scans":         scans,
		"summary": map[string]int{
			"total_scans":  len(scans),
			"total_findings": 0,
			"critical":     0,
			"high":         0,
			"medium":       0,
			"low":          0,
		},
	}

	// 计算统计数据
	for _, scan := range scans {
		report["summary"].(map[string]int)["total_findings"] += scan.Findings
		report["summary"].(map[string]int)["critical"] += scan.Critical
		report["summary"].(map[string]int)["high"] += scan.High
		report["summary"].(map[string]int)["medium"] += scan.Medium
		report["summary"].(map[string]int)["low"] += scan.Low
	}

	// 写入报告文件
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(reportFile, data, 0644); err != nil {
		return "", err
	}

	return reportFile, nil
}
