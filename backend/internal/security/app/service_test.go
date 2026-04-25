package app

import (
	"testing"

	"laima/internal/security/domain"

	"github.com/stretchr/testify/assert"
)

func TestScanStatus(t *testing.T) {
	t.Run("scan status values", func(t *testing.T) {
		assert.Equal(t, "pending", domain.ScanStatusPending)
		assert.Equal(t, "running", domain.ScanStatusRunning)
		assert.Equal(t, "completed", domain.ScanStatusCompleted)
		assert.Equal(t, "failed", domain.ScanStatusFailed)
	})
}

func TestScanType(t *testing.T) {
	t.Run("scan type values", func(t *testing.T) {
		assert.Equal(t, "sast", domain.ScanTypeSAST)
		assert.Equal(t, "dependency", domain.ScanTypeDependency)
		assert.Equal(t, "secret", domain.ScanTypeSecret)
		assert.Equal(t, "license", domain.ScanTypeLicense)
		assert.Equal(t, "dast", domain.ScanTypeDAST)
		assert.Equal(t, "container", domain.ScanTypeContainer)
	})
}

func TestSeverity(t *testing.T) {
	t.Run("severity values", func(t *testing.T) {
		assert.Equal(t, "critical", domain.SeverityCritical)
		assert.Equal(t, "high", domain.SeverityHigh)
		assert.Equal(t, "medium", domain.SeverityMedium)
		assert.Equal(t, "low", domain.SeverityLow)
		assert.Equal(t, "info", domain.SeverityInfo)
	})
}

func TestSecurityScan(t *testing.T) {
	t.Run("create security scan", func(t *testing.T) {
		scan := &domain.SecurityScan{
			RepositoryID: 1,
			ScanType:     domain.ScanTypeSAST,
			Status:       domain.ScanStatusPending,
		}

		assert.Equal(t, 1, scan.RepositoryID)
		assert.Equal(t, domain.ScanTypeSAST, scan.ScanType)
		assert.Equal(t, domain.ScanStatusPending, scan.Status)
	})
}

func TestScanFinding(t *testing.T) {
	t.Run("create scan finding", func(t *testing.T) {
		finding := &domain.ScanFinding{
			ScanID:         1,
			RepositoryID:   1,
			Severity:       domain.SeverityHigh,
			Title:          "SQL Injection",
			Description:    "Potential SQL injection vulnerability",
			Filepath:       "src/db/query.go",
			LineStart:      42,
			LineEnd:        45,
			CodeSnippet:    "db.Query(\"SELECT * FROM users WHERE id = \" + userID)",
			CWE:            "CWE-89",
			CVSS:           8.5,
			Recommendation: "Use parameterized queries",
		}

		assert.Equal(t, 1, finding.ScanID)
		assert.Equal(t, domain.SeverityHigh, finding.Severity)
		assert.Equal(t, "SQL Injection", finding.Title)
		assert.Equal(t, "CWE-89", finding.CWE)
		assert.Equal(t, 8.5, finding.CVSS)
	})
}

func TestDASTConfig(t *testing.T) {
	t.Run("create DAST config", func(t *testing.T) {
		config := &domain.DASTConfig{
			RepositoryID: 1,
			TargetURL:    "https://example.com",
			ScanType:     "full",
			Enabled:      true,
		}

		assert.Equal(t, 1, config.RepositoryID)
		assert.Equal(t, "https://example.com", config.TargetURL)
		assert.Equal(t, "full", config.ScanType)
		assert.True(t, config.Enabled)
	})
}

func TestContainerScanConfig(t *testing.T) {
	t.Run("create container scan config", func(t *testing.T) {
		config := &domain.ContainerScanConfig{
			RepositoryID:  1,
			ImageName:     "myapp:latest",
			RegistryURL:   "https://registry.example.com",
			ScanType:      "vulnerability",
			Enabled:       true,
		}

		assert.Equal(t, 1, config.RepositoryID)
		assert.Equal(t, "myapp:latest", config.ImageName)
		assert.Equal(t, "https://registry.example.com", config.RegistryURL)
		assert.Equal(t, "vulnerability", config.ScanType)
		assert.True(t, config.Enabled)
	})
}

func TestComplianceReport(t *testing.T) {
	t.Run("create compliance report", func(t *testing.T) {
		report := &domain.ComplianceReport{
			RepositoryID: 1,
			ReportType:   "security",
			Status:       "passed",
			Score:        85.5,
		}

		assert.Equal(t, 1, report.RepositoryID)
		assert.Equal(t, "security", report.ReportType)
		assert.Equal(t, "passed", report.Status)
		assert.Equal(t, 85.5, report.Score)
	})
}
