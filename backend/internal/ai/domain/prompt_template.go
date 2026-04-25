package domain

import (
	"fmt"
	"strings"
	"sync"
	"text/template"
)

type PromptCategory string

const (
	PromptCategoryBugDetection PromptCategory = "bug_detection"
	PromptCategorySecurity     PromptCategory = "security"
	PromptCategoryPerformance  PromptCategory = "performance"
	PromptCategoryStyle        PromptCategory = "style"
	PromptCategoryBestPractice PromptCategory = "best_practice"
	PromptCategorySummary      PromptCategory = "summary"
)

type PromptTemplate struct {
	ID          string         `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Category    PromptCategory `json:"category" gorm:"not null"`
	Language    string         `json:"language"`
	Description string         `json:"description" gorm:"type:text"`
	Template    string         `json:"template" gorm:"type:text;not null"`
	Variables   []string       `json:"variables" gorm:"type:text;serializer:json"`
	Enabled     bool           `json:"enabled" gorm:"default:true"`
	CreatedAt   string         `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt   string         `json:"updated_at" gorm:"not null;default:now()"`
}

type PromptTemplateManager interface {
	GetTemplate(id string) (*PromptTemplate, error)
	GetTemplatesByCategory(category PromptCategory) ([]*PromptTemplate, error)
	GetTemplatesByLanguage(language string) ([]*PromptTemplate, error)
	RenderTemplate(id string, variables map[string]interface{}) (string, error)
	AddTemplate(template *PromptTemplate) error
	UpdateTemplate(template *PromptTemplate) error
	DeleteTemplate(id string) error
}

type promptTemplateManager struct {
	templates map[string]*PromptTemplate
	mu        sync.RWMutex
}

func NewPromptTemplateManager() PromptTemplateManager {
	manager := &promptTemplateManager{
		templates: make(map[string]*PromptTemplate),
	}

	manager.loadDefaultTemplates()

	return manager
}

func (m *promptTemplateManager) GetTemplate(id string) (*PromptTemplate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	template, exists := m.templates[id]
	if !exists {
		return nil, fmt.Errorf("template not found: %s", id)
	}

	return template, nil
}

func (m *promptTemplateManager) GetTemplatesByCategory(category PromptCategory) ([]*PromptTemplate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*PromptTemplate
	for _, t := range m.templates {
		if t.Category == category && t.Enabled {
			result = append(result, t)
		}
	}

	return result, nil
}

func (m *promptTemplateManager) GetTemplatesByLanguage(language string) ([]*PromptTemplate, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*PromptTemplate
	for _, t := range m.templates {
		if (t.Language == "" || t.Language == language) && t.Enabled {
			result = append(result, t)
		}
	}

	return result, nil
}

func (m *promptTemplateManager) RenderTemplate(id string, variables map[string]interface{}) (string, error) {
	promptTemplate, err := m.GetTemplate(id)
	if err != nil {
		return "", err
	}

	tmpl, err := template.New("prompt").Parse(promptTemplate.Template)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, variables); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}

func (m *promptTemplateManager) AddTemplate(template *PromptTemplate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.templates[template.ID]; exists {
		return fmt.Errorf("template already exists: %s", template.ID)
	}

	m.templates[template.ID] = template
	return nil
}

func (m *promptTemplateManager) UpdateTemplate(template *PromptTemplate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.templates[template.ID]; !exists {
		return fmt.Errorf("template not found: %s", template.ID)
	}

	m.templates[template.ID] = template
	return nil
}

func (m *promptTemplateManager) DeleteTemplate(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.templates[id]; !exists {
		return fmt.Errorf("template not found: %s", id)
	}

	delete(m.templates, id)
	return nil
}

func (m *promptTemplateManager) loadDefaultTemplates() {
	templates := []*PromptTemplate{
		{
			ID:          "bug_detection_default",
			Name:        "Bug Detection - Default",
			Category:    PromptCategoryBugDetection,
			Language:    "",
			Description: "Default template for bug detection in code changes",
			Template: `你是一位经验丰富的代码审查专家。请审查以下代码变更，检测潜在的 Bug。

## 项目信息
- 语言: {{.Language}}
- 框架: {{.Framework}}
- 项目描述: {{.ProjectDescription}}

## 代码变更
{{.Diff}}

## 审查要求
1. 重点关注：空指针、边界条件、竞态条件、资源泄漏、类型错误
2. 不要报告代码风格问题（由专门的分析器处理）
3. 每个发现必须包含：文件路径、行号、问题描述、修复建议
4. 仅报告你确信是 Bug 的问题，不要猜测

## 输出格式 (JSON)
{
  "findings": [
    {
      "path": "文件路径",
      "line": 行号,
      "message": "问题描述",
      "suggestion": "修复建议",
      "confidence": 0.0-1.0
    }
  ]
}`,
			Variables: []string{"Language", "Framework", "ProjectDescription", "Diff"},
			Enabled:   true,
		},
		{
			ID:          "security_detection_default",
			Name:        "Security Detection - Default",
			Category:    PromptCategorySecurity,
			Language:    "",
			Description: "Default template for security vulnerability detection",
			Template: `你是一位安全审计专家。请审查以下代码变更，检测安全漏洞。

## 代码变更
{{.Diff}}

## 审查要求
1. 重点关注：SQL 注入、XSS、CSRF、路径遍历、硬编码密钥、不安全的反序列化
2. 参考 OWASP Top 10 进行检查
3. 每个发现必须包含 CWE 编号（如适用）

## 输出格式 (JSON)
{
  "findings": [
    {
      "path": "文件路径",
      "line": 行号,
      "message": "问题描述",
      "cwe": "CWE-XXX",
      "suggestion": "修复建议",
      "confidence": 0.0-1.0
    }
  ]
}`,
			Variables: []string{"Diff"},
			Enabled:   true,
		},
		{
			ID:          "performance_analysis_default",
			Name:        "Performance Analysis - Default",
			Category:    PromptCategoryPerformance,
			Language:    "",
			Description: "Default template for performance issue detection",
			Template: `你是一位性能优化专家。请审查以下代码变更，识别性能问题。

## 代码变更
{{.Diff}}

## 审查要求
1. 重点关注：算法复杂度、内存泄漏、不必要的循环、数据库查询优化、缓存使用
2. 提供具体的优化建议
3. 评估性能影响的严重程度

## 输出格式 (JSON)
{
  "findings": [
    {
      "path": "文件路径",
      "line": 行号,
      "message": "问题描述",
      "suggestion": "优化建议",
      "impact": "high/medium/low",
      "confidence": 0.0-1.0
    }
  ]
}`,
			Variables: []string{"Diff"},
			Enabled:   true,
		},
		{
			ID:          "summary_generation_default",
			Name:        "Summary Generation - Default",
			Category:    PromptCategorySummary,
			Language:    "",
			Description: "Default template for generating PR summary",
			Template: `请为以下 Pull Request 生成简洁的变更摘要。

## PR 标题
{{.PRTitle}}

## PR 描述
{{.PRDescription}}

## 变更文件列表
{{.ChangedFiles}}

## 变更统计
- 新增 {{.Additions}} 行
- 删除 {{.Deletions}} 行
- 涉及 {{.FileCount}} 个文件

## 要求
用 2-3 句话概括本次变更的目的和影响。使用中文。`,
			Variables: []string{"PRTitle", "PRDescription", "ChangedFiles", "Additions", "Deletions", "FileCount"},
			Enabled:   true,
		},
	}

	for _, t := range templates {
		m.templates[t.ID] = t
	}
}
