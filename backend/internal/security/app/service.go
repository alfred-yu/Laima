package app

import (
	"laima/internal/security/domain"
)

// SecurityService 安全扫描服务接口
type SecurityService interface {
	// CreateScan 创建安全扫描
	CreateScan(req *domain.CreateScanRequest) (*domain.SecurityScan, error)

	// GetScan 获取扫描详情
	GetScan(scanID int) (*domain.SecurityScan, error)

	// ListScans 列出扫描列表
	ListScans(filter *domain.ScanFilter) ([]*domain.SecurityScan, int64, error)

	// StopScan 停止扫描
	StopScan(scanID int) error

	// GetScanFindings 获取扫描发现的问题
	GetScanFindings(scanID int, filter *domain.FindingFilter) ([]*domain.ScanFinding, int64, error)

	// GetFinding 获取问题详情
	GetFinding(findingID int) (*domain.ScanFinding, error)

	// UpdateDASTConfig 更新DAST扫描配置
	UpdateDASTConfig(config *domain.DASTScanConfig) error

	// GetDASTConfig 获取DAST扫描配置
	GetDASTConfig(repositoryID int) (*domain.DASTScanConfig, error)

	// UpdateContainerConfig 更新容器扫描配置
	UpdateContainerConfig(config *domain.ContainerScanConfig) error

	// GetContainerConfig 获取容器扫描配置
	GetContainerConfig(repositoryID int) (*domain.ContainerScanConfig, error)

	// RunDASTScan 运行DAST扫描
	RunDASTScan(scanID int) error

	// RunContainerScan 运行容器扫描
	RunContainerScan(scanID int) error

	// GenerateComplianceReport 生成合规报告
	GenerateComplianceReport(repositoryID int, scanTypes []string) (string, error)
}
