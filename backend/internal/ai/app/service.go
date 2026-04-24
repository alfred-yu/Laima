package app

import (
	"context"
	"fmt"
	"strings"
	"laima/internal/ai/domain"
	"laima/internal/git"
	prapp "laima/internal/pr/app"
	prdomain "laima/internal/pr/domain"
	repodomain "laima/internal/repo/domain"

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
	db     *gorm.DB
	gitSvc *git.Service
	prSvc  prapp.PRService
}

// NewAIService 创建AI服务实例
func NewAIService(db *gorm.DB, gitSvc *git.Service, prSvc prapp.PRService) AIService {
	return &aiService{
		db:     db,
		gitSvc: gitSvc,
		prSvc:  prSvc,
	}
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

	// 异步执行AI审查
	go func() {
		ctx := context.Background()
		_ = s.UpdateAIReviewStatus(ctx, review.ID, domain.AIReviewStatusRunning)

		// 获取仓库信息
		var repo repodomain.Repository
		if err := s.db.First(&repo, int64(req.RepositoryID)).Error; err != nil {
			_ = s.FailAIReview(ctx, review.ID, fmt.Sprintf("获取仓库信息失败: %v", err))
			return
		}

		// 解析所有者和仓库名
		parts := strings.Split(repo.FullPath, "/")
		if len(parts) < 2 {
			_ = s.FailAIReview(ctx, review.ID, "无效的仓库路径")
			return
		}
		owner := parts[0]
		repoName := parts[1]

		// 获取PR信息以获取分支
		var pr prdomain.PullRequest
		if err := s.db.First(&pr, req.PullRequestID).Error; err != nil {
			_ = s.FailAIReview(ctx, review.ID, fmt.Sprintf("获取PR信息失败: %v", err))
			return
		}

		// 获取变更diff
		diff, err := s.gitSvc.GetDiff(owner, repoName, pr.TargetBranch, pr.SourceBranch)
		if err != nil {
			_ = s.FailAIReview(ctx, review.ID, fmt.Sprintf("获取变更diff失败: %v", err))
			return
		}

		// 解析diff并提取变更（简化实现）
		issues := s.extractIssuesFromDiff(diff, review.ID, req.PullRequestID)

		// 完成审查
		_ = s.CompleteAIReview(ctx, review.ID, 0.88, "代码整体质量良好，有少量改进建议", issues)
	}()

	return review, nil
}

// extractIssuesFromDiff 从diff中提取审查问题（简化实现）
func (s *aiService) extractIssuesFromDiff(diff string, reviewID, prID int) []*domain.AIReviewIssue {
	var issues []*domain.AIReviewIssue

	// 简单解析diff，查找变更的文件
	lines := strings.Split(diff, "\n")
	var currentFile string
	var lineNumber int

	for _, line := range lines {
		if strings.HasPrefix(line, "diff --git") {
			// 新文件开始
			parts := strings.Split(line, " ")
			if len(parts) >= 4 {
				currentFile = strings.TrimPrefix(parts[2], "a/")
			}
			lineNumber = 0
			continue
		}

		if strings.HasPrefix(line, "@@") {
			// 解析行号信息
			continue
		}

		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			lineNumber++
			// 简化的问题检测逻辑
			content := strings.TrimSpace(strings.TrimPrefix(line, "+"))
			if len(content) > 0 {
				// 检查一些常见问题（简化示例）
				if strings.Contains(content, "TODO") || strings.Contains(content, "FIXME") {
					issues = append(issues, &domain.AIReviewIssue{
						AIReviewID:    reviewID,
						PullRequestID: prID,
						Severity:      domain.AIReviewSeverityLow,
						Category:      "code_quality",
						Title:         "待办事项标记",
						Description:   "代码中包含待办事项标记，建议在合并前完成",
						Path:          currentFile,
						Line:          lineNumber,
						Suggestion:    "完成相关功能或移除待办标记",
						Confidence:    0.90,
					})
				}

				if len(content) > 120 {
					issues = append(issues, &domain.AIReviewIssue{
						AIReviewID:    reviewID,
						PullRequestID: prID,
						Severity:      domain.AIReviewSeverityLow,
						Category:      "code_style",
						Title:         "行过长",
						Description:   "代码行过长，建议拆分为多行提高可读性",
						Path:          currentFile,
						Line:          lineNumber,
						Suggestion:    "将长行拆分为多行",
						Confidence:    0.85,
					})
				}
			}
		} else if !strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			lineNumber++
		}
	}

	return issues
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
	// 获取PR信息
	pr, err := s.prSvc.GetPR(ctx, prID)
	if err != nil {
		return nil, fmt.Errorf("获取PR信息失败: %w", err)
	}

	// 构建审查请求
	req := &domain.AIReviewRequest{
		PullRequestID: prID,
		RepositoryID:  pr.RepositoryID,
		HeadCommitSHA: pr.HeadCommitSHA,
		BaseCommitSHA: pr.BaseCommitSHA,
	}

	// 触发审查
	return s.TriggerAIReview(ctx, req)
}
