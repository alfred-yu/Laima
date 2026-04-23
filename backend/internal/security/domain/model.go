package domain

import (
	"time"
)

// 安全扫描类型常量
const (
	ScanTypeSAST    = "sast"    // 静态代码扫描
	ScanTypeDAST    = "dast"    // 动态应用安全测试
	ScanTypeSCA     = "sca"     // 依赖扫描
	ScanTypeSecret  = "secret"  // 密钥检测
	ScanTypeLicense = "license" // 许可证检查
	ScanTypeContainer = "container" // 容器扫描
)

// 扫描状态常量
const (
	ScanStatusPending   = "pending"   // 待处理
	ScanStatusRunning   = "running"   // 运行中
	ScanStatusCompleted = "completed" // 完成
	ScanStatusFailed    = "failed"    // 失败
)

// 漏洞严重程度常量
const (
	SeverityLow      = "low"      // 低
	SeverityMedium   = "medium"   // 中
	SeverityHigh     = "high"     // 高
	SeverityCritical = "critical" // 严重
)

// SecurityScan 安全扫描模型
type SecurityScan struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	RepositoryID int       `json:"repository_id" gorm:"not null;index"`
	ScanType     string    `json:"scan_type" gorm:"not null"`
	Status       string    `json:"status" gorm:"not null;default:'pending'"`
	Branch       string    `json:"branch" gorm:"not null"`
	Commit       string    `json:"commit" gorm:"not null"`
	StartTime    time.Time `json:"start_time" gorm:"index"`
	EndTime      time.Time `json:"end_time" gorm:"index"`
	Duration     int       `json:"duration"` // 扫描持续时间（秒）
	Findings     int       `json:"findings"` // 发现的问题数量
	Critical     int       `json:"critical"` // 严重问题数量
	High         int       `json:"high"`     // 高危问题数量
	Medium       int       `json:"medium"`   // 中危问题数量
	Low          int       `json:"low"`      // 低危问题数量
	ReportURL    string    `json:"report_url" gorm:"size:512"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// ScanFinding 扫描发现的问题模型
type ScanFinding struct {
	ID             int       `json:"id" gorm:"primaryKey"`
	ScanID         int       `json:"scan_id" gorm:"not null;index"`
	RepositoryID   int       `json:"repository_id" gorm:"not null;index"`
	Severity       string    `json:"severity" gorm:"not null"`
	Title          string    `json:"title" gorm:"not null;size:512"`
	Description    string    `json:"description" gorm:"type:text"`
	Filepath       string    `json:"filepath" gorm:"size:512"`
	LineStart      int       `json:"line_start"`
	LineEnd        int       `json:"line_end"`
	CodeSnippet    string    `json:"code_snippet" gorm:"type:text"`
	CWE            string    `json:"cwe" gorm:"size:64"` // Common Weakness Enumeration
	CVSS           float64   `json:"cvss"` // Common Vulnerability Scoring System
	Recommendation string    `json:"recommendation" gorm:"type:text"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null;default:now()"`
}

// DASTScanConfig DAST扫描配置
type DASTScanConfig struct {
	ID              int    `json:"id" gorm:"primaryKey"`
	RepositoryID    int    `json:"repository_id" gorm:"not null;uniqueIndex"`
	TargetURL       string `json:"target_url" gorm:"not null;size:512"`
	LoginURL        string `json:"login_url" gorm:"size:512"`
	Username        string `json:"username" gorm:"size:256"`
	Password        string `json:"password" gorm:"size:256"`
	CustomHeaders   string `json:"custom_headers" gorm:"type:jsonb"`
	AllowedHosts    string `json:"allowed_hosts" gorm:"type:jsonb"`
	ExcludePaths    string `json:"exclude_paths" gorm:"type:jsonb"`
	ScanDepth       int    `json:"scan_depth" gorm:"default:10"`
	Concurrency     int    `json:"concurrency" gorm:"default:5"`
	EnableCrawling  bool   `json:"enable_crawling" gorm:"default:true"`
	EnableXSS       bool   `json:"enable_xss" gorm:"default:true"`
	EnableSQLi      bool   `json:"enable_sqli" gorm:"default:true"`
	EnableCSRF      bool   `json:"enable_csrf" gorm:"default:true"`
	CreatedAt       time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// ContainerScanConfig 容器扫描配置
type ContainerScanConfig struct {
	ID              int       `json:"id" gorm:"primaryKey"`
	RepositoryID    int       `json:"repository_id" gorm:"not null;uniqueIndex"`
	ImageName       string    `json:"image_name" gorm:"not null;size:512"`
	ImageTag        string    `json:"image_tag" gorm:"not null;size:128"`
	RegistryURL     string    `json:"registry_url" gorm:"size:512"`
	RegistryUsername string   `json:"registry_username" gorm:"size:256"`
	RegistryPassword string   `json:"registry_password" gorm:"size:256"`
	IncludeOS       bool      `json:"include_os" gorm:"default:true"`
	IncludeApps     bool      `json:"include_apps" gorm:"default:true"`
	IncludeSecrets  bool      `json:"include_secrets" gorm:"default:true"`
	CreatedAt       time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

// CreateScanRequest 创建扫描请求
type CreateScanRequest struct {
	RepositoryID int    `json:"repository_id" binding:"required"`
	ScanType     string `json:"scan_type" binding:"required,oneof=sast dast sca secret license container"`
	Branch       string `json:"branch" binding:"required"`
	Commit       string `json:"commit" binding:"required"`
	Config       map[string]interface{} `json:"config"`
}

// ScanFilter 扫描过滤条件
type ScanFilter struct {
	RepositoryID int    `json:"repository_id"`
	ScanType     string `json:"scan_type"`
	Status       string `json:"status"`
	Branch       string `json:"branch"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Page         int    `json:"page" binding:"min=1"`
	PerPage      int    `json:"per_page" binding:"min=1,max=100"`
}

// FindingFilter 扫描发现过滤条件
type FindingFilter struct {
	ScanID       int    `json:"scan_id"`
	RepositoryID int    `json:"repository_id"`
	Severity     string `json:"severity"`
	CWE          string `json:"cwe"`
	Filepath     string `json:"filepath"`
	Search       string `json:"search"`
	Page         int    `json:"page" binding:"min=1"`
	PerPage      int    `json:"per_page" binding:"min=1,max=100"`
}
