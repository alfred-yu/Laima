package app

import (
	"context"
	"laima/internal/ai/domain"

	"gorm.io/gorm"
)

// AIService AI服务接口
type AIService interface {
	// AI审查
	TriggerAIReview(ctx context.Context, req *domain.AIReviewRequest) (*domain.AIReview, error)
	GetAIReview(ctx context.Context, reviewID int) (*domain.AIReview, error)
	GetAIReviewByPR(ctx context.Context, prID int) (*domain.AIReview, error)
	ListAIReviews(ctx context.Context, filter *domain.AIReviewFilter) ([]*domain.AIReview, int64, error)

	// 审查问题管理
	GetAIReviewIssues(ctx context.Context, reviewID int) ([]*domain.AIReviewIssue, error)
	GetAIReviewIssuesByPR(ctx context.Context, prID int) ([]*domain.AIReviewIssue, error)

	// 审查状态更新
	UpdateAIReviewStatus(ctx context.Context, reviewID int, status string) error
	CompleteAIReview(ctx context.Context, reviewID int, score float64, summary string, issues []*domain.AIReviewIssue) error
	FailAIReview(ctx context.Context, reviewID int, errorMsg string) error

	// 集成触发
	TriggerAIReviewForPR(ctx context.Context, prID int) (*domain.AIReview, error)
}

// aiService AI服务实现
type aiService struct {
	db *gorm.DB
}

// NewAIService 创建AI服务实例
func NewAIService(db *gorm.DB) AIService {
	return &aiService{db: db}
}

// TriggerAIReview 触发AI审查
func (s *aiService) TriggerAIReview(ctx context.Context, req *domain.AIReviewRequest) (*domain.AIReview, error) {
	// 创建审查记录
	review := &domain.AIReview{
		PullRequestID: req.PullRequestID,
		RepositoryID:  req.RepositoryID,
		Status:        domain.AIReviewStatusPending,
	}

	if err := s.db.Create(review).Error; err != nil {
		return nil, err
	}

	// 异步执行AI审查（模拟）
	go func() {
		ctx := context.Background()
		_ = s.UpdateAIReviewStatus(ctx, review.ID, domain.AIReviewStatusRunning)

		// 模拟AI审查过程
		issues := []*domain.AIReviewIssue{
			{
				AIReviewID:    review.ID,
				PullRequestID: req.PullRequestID,
				Severity:      domain.AIReviewSeverityMedium,
				Category:      "code_quality",
				Title:         "代码可读性改进建议",
				Description:   "建议增加注释提高代码可读性",
				Path:          "src/main.go",
				Line:          42,
				Suggestion:    "添加函数和关键逻辑的注释",
				Confidence:    0.85,
			},
		}

		_ = s.CompleteAIReview(ctx, review.ID, 0.88, "代码整体质量良好，有少量改进建议", issues)
	}()

	return review, nil
}

// GetAIReview 根据ID获取AI审查
func (s *aiService) GetAIReview(ctx context.Context, reviewID int) (*domain.AIReview, error) {
	var review domain.AIReview
	result := s.db.First(&review, reviewID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &review, nil
}

// GetAIReviewByPR 根据PR ID获取AI审查
func (s *aiService) GetAIReviewByPR(ctx context.Context, prID int) (*domain.AIReview, error) {
	var review domain.AIReview
	result := s.db.Where("pull_request_id = ?", prID).Order("created_at DESC").First(&review)
	if result.Error != nil {
		return nil, result.Error
	}
	return &review, nil
}

// ListAIReviews 列出AI审查
func (s *aiService) ListAIReviews(ctx context.Context, filter *domain.AIReviewFilter) ([]*domain.AIReview, int64, error) {
	var reviews []*domain.AIReview
	var total int64

	query := s.db.Model(&domain.AIReview{})

	// 应用过滤条件
	if filter.PullRequestID > 0 {
		query = query.Where("pull_request_id = ?", filter.PullRequestID)
	}
	if filter.RepositoryID > 0 {
		query = query.Where("repository_id = ?", filter.RepositoryID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// 计算总数
	query.Count(&total)

	// 分页
	offset := (filter.Page - 1) * filter.PerPage
	result := query.Offset(offset).Limit(filter.PerPage).Order("created_at DESC").Find(&reviews)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return reviews, total, nil
}

// GetAIReviewIssues 获取AI审查问题
func (s *aiService) GetAIReviewIssues(ctx context.Context, reviewID int) ([]*domain.AIReviewIssue, error) {
	var issues []*domain.AIReviewIssue
	result := s.db.Where("ai_review_id = ?", reviewID).Order("severity DESC, created_at ASC").Find(&issues)
	if result.Error != nil {
		return nil, result.Error
	}
	return issues, nil
}

// GetAIReviewIssuesByPR 根据PR ID获取AI审查问题
func (s *aiService) GetAIReviewIssuesByPR(ctx context.Context, prID int) ([]*domain.AIReviewIssue, error) {
	var issues []*domain.AIReviewIssue
	result := s.db.Where("pull_request_id = ?", prID).Order("severity DESC, created_at ASC").Find(&issues)
	if result.Error != nil {
		return nil, result.Error
	}
	return issues, nil
}

// UpdateAIReviewStatus 更新AI审查状态
func (s *aiService) UpdateAIReviewStatus(ctx context.Context, reviewID int, status string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": s.db.NowFunc(),
	}

	if status == domain.AIReviewStatusRunning {
		updates["started_at"] = s.db.NowFunc()
	}

	return s.db.Model(&domain.AIReview{}).Where("id = ?", reviewID).Updates(updates).Error
}

// CompleteAIReview 完成AI审查
func (s *aiService) CompleteAIReview(ctx context.Context, reviewID int, score float64, summary string, issues []*domain.AIReviewIssue) error {
	// 开始事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 更新审查状态
		if err := tx.Model(&domain.AIReview{}).Where("id = ?", reviewID).Updates(map[string]interface{}{
			"status":       domain.AIReviewStatusCompleted,
			"score":        score,
			"summary":      summary,
			"completed_at": s.db.NowFunc(),
			"updated_at":   s.db.NowFunc(),
		}).Error; err != nil {
			return err
		}

		// 保存审查问题
		for _, issue := range issues {
			issue.AIReviewID = reviewID
			if err := tx.Create(issue).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// FailAIReview 失败AI审查
func (s *aiService) FailAIReview(ctx context.Context, reviewID int, errorMsg string) error {
	return s.db.Model(&domain.AIReview{}).Where("id = ?", reviewID).Updates(map[string]interface{}{
		"status":     domain.AIReviewStatusFailed,
		"error":      errorMsg,
		"updated_at": s.db.NowFunc(),
	}).Error
}

// TriggerAIReviewForPR 为PR触发AI审查
func (s *aiService) TriggerAIReviewForPR(ctx context.Context, prID int) (*domain.AIReview, error) {
	// 模拟获取PR信息
	req := &domain.AIReviewRequest{
		PullRequestID: prID,
		RepositoryID:  1,
		HeadCommitSHA: "abc123",
		BaseCommitSHA: "def456",
	}

	// 触发审查
	return s.TriggerAIReview(ctx, req)
}
