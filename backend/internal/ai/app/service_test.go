package app

import (
	"context"
	"testing"

	"laima/internal/ai/domain"
	"laima/internal/ai/infrastructure"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLLMProvider struct {
	mock.Mock
}

func (m *MockLLMProvider) Name() infrastructure.ProviderType {
	args := m.Called()
	return args.Get(0).(infrastructure.ProviderType)
}

func (m *MockLLMProvider) ChatCompletion(ctx context.Context, req *infrastructure.ChatRequest) (*infrastructure.ChatResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*infrastructure.ChatResponse), args.Error(1)
}

func (m *MockLLMProvider) StreamChatCompletion(ctx context.Context, req *infrastructure.ChatRequest) (<-chan *infrastructure.ChatChunk, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(<-chan *infrastructure.ChatChunk), args.Error(1)
}

func (m *MockLLMProvider) IsAvailable(ctx context.Context) bool {
	args := m.Called(ctx)
	return args.Bool(0)
}

func (m *MockLLMProvider) GetModel() string {
	args := m.Called()
	return args.String(0)
}

type MockAIReviewRepository struct {
	mock.Mock
}

func (m *MockAIReviewRepository) Create(review *domain.AIReview) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockAIReviewRepository) GetByID(id int) (*domain.AIReview, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.AIReview), args.Error(1)
}

func (m *MockAIReviewRepository) Update(review *domain.AIReview) error {
	args := m.Called(review)
	return args.Error(0)
}

func TestPromptTemplateManager(t *testing.T) {
	manager := domain.NewPromptTemplateManager()

	t.Run("get bug detection template", func(t *testing.T) {
		template, err := manager.GetTemplate("bug_detection")

		assert.NoError(t, err)
		assert.NotNil(t, template)
		assert.Equal(t, domain.PromptCategoryBugDetection, template.Category)
	})

	t.Run("get security template", func(t *testing.T) {
		template, err := manager.GetTemplate("security")

		assert.NoError(t, err)
		assert.NotNil(t, template)
		assert.Equal(t, domain.PromptCategorySecurity, template.Category)
	})

	t.Run("render template", func(t *testing.T) {
		variables := map[string]interface{}{
			"Language":    "Go",
			"CodeDiff":    "+ func test() {}",
			"FileName":    "test.go",
			"Context":     "Test context",
		}

		rendered, err := manager.RenderTemplate("bug_detection", variables)

		assert.NoError(t, err)
		assert.Contains(t, rendered, "Go")
		assert.Contains(t, rendered, "test.go")
	})

	t.Run("template not found", func(t *testing.T) {
		template, err := manager.GetTemplate("nonexistent")

		assert.Error(t, err)
		assert.Nil(t, template)
	})
}

func TestResultAggregator(t *testing.T) {
	aggregator := NewResultAggregator()

	t.Run("aggregate findings", func(t *testing.T) {
		findings := []*domain.AIReviewIssue{
			{
				ID:          1,
				AIReviewID:  1,
				Path:        "test.go",
				Line:        10,
				Severity:    domain.AIReviewSeverityHigh,
				Category:    "bug",
				Description: "Potential nil pointer dereference",
				Suggestion:  "Add nil check",
				Confidence:  0.85,
			},
			{
				ID:          2,
				AIReviewID:  1,
				Path:        "test.go",
				Line:        20,
				Severity:    domain.AIReviewSeverityMedium,
				Category:    "style",
				Description: "Unused variable",
				Suggestion:  "Remove unused variable",
				Confidence:  0.90,
			},
		}

		aggregated, err := aggregator.Aggregate(findings)

		assert.NoError(t, err)
		assert.Len(t, aggregated, 2)
		assert.Equal(t, domain.AIReviewSeverityHigh, aggregated[0].Severity)
	})

	t.Run("deduplicate findings", func(t *testing.T) {
		findings := []*domain.AIReviewIssue{
			{
				ID:          1,
				Path:        "test.go",
				Line:        10,
				Category:    "bug",
				Description: "Issue 1",
			},
			{
				ID:          2,
				Path:        "test.go",
				Line:        10,
				Category:    "bug",
				Description: "Issue 2",
			},
		}

		deduplicated := aggregator.Deduplicate(findings)

		assert.Len(t, deduplicated, 1)
	})

	t.Run("sort by priority", func(t *testing.T) {
		findings := []*AggregatedFinding{
			{
				ID:         "1",
				Path:       "test.go",
				Line:       10,
				Severity:   domain.AIReviewSeverityLow,
				Confidence: 0.8,
			},
			{
				ID:         "2",
				Path:       "test.go",
				Line:       20,
				Severity:   domain.AIReviewSeverityCritical,
				Confidence: 0.9,
			},
			{
				ID:         "3",
				Path:       "test.go",
				Line:       30,
				Severity:   domain.AIReviewSeverityHigh,
				Confidence: 0.85,
			},
		}

		sorted := aggregator.SortByPriority(findings)

		assert.Equal(t, domain.AIReviewSeverityCritical, sorted[0].Severity)
		assert.Equal(t, domain.AIReviewSeverityHigh, sorted[1].Severity)
		assert.Equal(t, domain.AIReviewSeverityLow, sorted[2].Severity)
	})

	t.Run("generate summary", func(t *testing.T) {
		findings := []*AggregatedFinding{
			{
				ID:         "1",
				Path:       "test.go",
				Line:       10,
				Severity:   domain.AIReviewSeverityHigh,
				Message:    "Potential nil pointer dereference",
				Confidence: 0.85,
			},
			{
				ID:         "2",
				Path:       "test.go",
				Line:       20,
				Severity:   domain.AIReviewSeverityMedium,
				Message:    "Unused variable",
				Confidence: 0.90,
			},
		}

		summary := aggregator.GenerateSummary(findings)

		assert.Contains(t, summary, "1个高危问题")
		assert.Contains(t, summary, "1个中危问题")
	})
}

func TestLLMProviderFactory(t *testing.T) {
	factory := infrastructure.NewProviderFactory()

	t.Run("get available providers", func(t *testing.T) {
		providers := factory.GetAvailableProviders()

		assert.Len(t, providers, 4)
		assert.Contains(t, providers, infrastructure.ProviderTypeOpenAI)
		assert.Contains(t, providers, infrastructure.ProviderTypeOllama)
		assert.Contains(t, providers, infrastructure.ProviderTypeQwen)
		assert.Contains(t, providers, infrastructure.ProviderTypeDeepSeek)
	})
}
