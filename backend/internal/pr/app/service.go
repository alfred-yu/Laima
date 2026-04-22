package app

import (
	"context"
	"errors"
	"fmt"
	"laima/internal/git"
	repodomain "laima/internal/repo/domain"
	"laima/internal/pr/domain"
	"time"

	"gorm.io/gorm"
)

// PRService PR服务接口
type PRService interface {
	// PR CRUD
	CreatePR(ctx context.Context, req *domain.CreatePRRequest, authorID int) (*domain.PullRequest, error)
	GetPR(ctx context.Context, prID int) (*domain.PullRequest, error)
	GetPRByNumber(ctx context.Context, repoID int, number int) (*domain.PullRequest, error)
	UpdatePR(ctx context.Context, prID int, req *domain.UpdatePRRequest) (*domain.PullRequest, error)
	DeletePR(ctx context.Context, prID int) error
	ListPRs(ctx context.Context, filter *domain.PRFilter) ([]*domain.PullRequest, int64, error)

	// PR 操作
	MergePR(ctx context.Context, prID int, userID int) (*domain.PullRequest, error)
	ClosePR(ctx context.Context, prID int, userID int) (*domain.PullRequest, error)
	ReopenPR(ctx context.Context, prID int, userID int) (*domain.PullRequest, error)

	// 审查操作
	CreateReview(ctx context.Context, prID int, req *domain.ReviewRequest, reviewerID int) (*domain.Review, error)
	GetReviews(ctx context.Context, prID int) ([]*domain.Review, error)
	CreateReviewComment(ctx context.Context, prID int, req *domain.ReviewCommentRequest, authorID int) (*domain.ReviewComment, error)
	GetReviewComments(ctx context.Context, prID int) ([]*domain.ReviewComment, error)

	// 状态检查
	CheckMergeability(ctx context.Context, prID int) (bool, error)
	UpdateMergeState(ctx context.Context, prID int) error
}

// prService PR服务实现
type prService struct {
	db     *gorm.DB
	gitSvc *git.Service
}

// NewPRService 创建PR服务实例
func NewPRService(db *gorm.DB, gitSvc *git.Service) PRService {
	return &prService{db: db, gitSvc: gitSvc}
}

// CreatePR 创建PR
func (s *prService) CreatePR(ctx context.Context, req *domain.CreatePRRequest, authorID int) (*domain.PullRequest, error) {
	// 验证目标仓库存在
	var targetRepo repodomain.Repository
	if err := s.db.First(&targetRepo, req.RepositoryID).Error; err != nil {
		return nil, fmt.Errorf("target repository not found: %w", err)
	}

	// 验证源仓库存在
	var sourceRepo repodomain.Repository
	if err := s.db.First(&sourceRepo, req.SourceRepoID).Error; err != nil {
		return nil, fmt.Errorf("source repository not found: %w", err)
	}

	// 生成PR编号
	var maxNumber int
	s.db.Model(&domain.PullRequest{}).Where("repository_id = ?", req.RepositoryID).Select("COALESCE(MAX(number), 0)").Scan(&maxNumber)
	prNumber := maxNumber + 1

	// 创建PR记录
	pr := &domain.PullRequest{
		Number:          prNumber,
		Title:           req.Title,
		Description:     req.Description,
		RepositoryID:    req.RepositoryID,
		AuthorID:        authorID,
		SourceRepoID:    req.SourceRepoID,
		SourceBranch:    req.SourceBranch,
		TargetBranch:    req.TargetBranch,
		State:           "open",
		MergeState:      "checking",
		ReviewMode:      "standard",
		HeadCommitSHA:   "placeholder",
		BaseCommitSHA:   "placeholder",
		IsDraft:         req.IsDraft,
		AIReviewStatus:  "pending",
	}

	if err := s.db.Create(pr).Error; err != nil {
		return nil, err
	}

	// 触发合并状态检查
	go s.UpdateMergeState(context.Background(), pr.ID)

	return pr, nil
}

// GetPR 根据ID获取PR
func (s *prService) GetPR(ctx context.Context, prID int) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	result := s.db.First(&pr, prID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pr, nil
}

// GetPRByNumber 根据编号获取PR
func (s *prService) GetPRByNumber(ctx context.Context, repoID int, number int) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	result := s.db.Where("repository_id = ? AND number = ?", repoID, number).First(&pr)
	if result.Error != nil {
		return nil, result.Error
	}
	return &pr, nil
}

// UpdatePR 更新PR
func (s *prService) UpdatePR(ctx context.Context, prID int, req *domain.UpdatePRRequest) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	if err := s.db.First(&pr, prID).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.TargetBranch != "" {
		updates["target_branch"] = req.TargetBranch
	}
	updates["is_draft"] = req.IsDraft

	if len(updates) > 0 {
		if err := s.db.Model(&pr).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return &pr, nil
}

// DeletePR 删除PR
func (s *prService) DeletePR(ctx context.Context, prID int) error {
	return s.db.Delete(&domain.PullRequest{}, prID).Error
}

// ListPRs 列出PR
func (s *prService) ListPRs(ctx context.Context, filter *domain.PRFilter) ([]*domain.PullRequest, int64, error) {
	var prs []*domain.PullRequest
	var total int64

	query := s.db.Model(&domain.PullRequest{})

	if filter.RepositoryID > 0 {
		query = query.Where("repository_id = ?", filter.RepositoryID)
	}
	if filter.AuthorID > 0 {
		query = query.Where("author_id = ?", filter.AuthorID)
	}
	if filter.State != "" {
		query = query.Where("state = ?", filter.State)
	}
	if filter.Search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	query.Count(&total)

	offset := (filter.Page - 1) * filter.PerPage
	result := query.Offset(offset).Limit(filter.PerPage).Order("created_at DESC").Find(&prs)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return prs, total, nil
}

// MergePR 合并PR
func (s *prService) MergePR(ctx context.Context, prID int, userID int) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	if err := s.db.First(&pr, prID).Error; err != nil {
		return nil, err
	}

	if pr.State != "open" {
		return nil, errors.New("PR is not open")
	}

	// 检查可合并性
	mergeable, err := s.CheckMergeability(ctx, prID)
	if err != nil {
		return nil, err
	}
	if !mergeable {
		return nil, errors.New("PR is not mergeable")
	}

	// 更新PR状态
	pr.State = "merged"
	pr.MergedBy = userID
	pr.MergedAt = time.Now()
	pr.MergeCommitSHA = "placeholder_merge_sha"

	if err := s.db.Save(&pr).Error; err != nil {
		return nil, err
	}

	return &pr, nil
}

// ClosePR 关闭PR
func (s *prService) ClosePR(ctx context.Context, prID int, userID int) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	if err := s.db.First(&pr, prID).Error; err != nil {
		return nil, err
	}

	if pr.State != "open" {
		return nil, errors.New("PR is not open")
	}

	pr.State = "closed"
	pr.ClosedAt = time.Now()

	if err := s.db.Save(&pr).Error; err != nil {
		return nil, err
	}

	return &pr, nil
}

// ReopenPR 重新打开PR
func (s *prService) ReopenPR(ctx context.Context, prID int, userID int) (*domain.PullRequest, error) {
	var pr domain.PullRequest
	if err := s.db.First(&pr, prID).Error; err != nil {
		return nil, err
	}

	if pr.State != "closed" {
		return nil, errors.New("PR is not closed")
	}

	pr.State = "open"
	pr.ClosedAt = time.Time{}

	if err := s.db.Save(&pr).Error; err != nil {
		return nil, err
	}

	return &pr, nil
}

// CreateReview 创建审查
func (s *prService) CreateReview(ctx context.Context, prID int, req *domain.ReviewRequest, reviewerID int) (*domain.Review, error) {
	var pr domain.PullRequest
	if err := s.db.First(&pr, prID).Error; err != nil {
		return nil, err
	}

	review := &domain.Review{
		PullRequestID: prID,
		ReviewerID:    reviewerID,
		State:         req.State,
		Score:         req.Score,
		Body:          req.Body,
		SubmittedAt:   time.Now(),
	}

	if err := s.db.Create(review).Error; err != nil {
		return nil, err
	}

	return review, nil
}

// GetReviews 获取审查列表
func (s *prService) GetReviews(ctx context.Context, prID int) ([]*domain.Review, error) {
	var reviews []*domain.Review
	result := s.db.Where("pull_request_id = ?", prID).Order("submitted_at DESC").Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}
	return reviews, nil
}

// CreateReviewComment 创建审查评论
func (s *prService) CreateReviewComment(ctx context.Context, prID int, req *domain.ReviewCommentRequest, authorID int) (*domain.ReviewComment, error) {
	var pr domain.PullRequest
	if err := s.db.First(&pr, prID).Error; err != nil {
		return nil, err
	}

	comment := &domain.ReviewComment{
		PullRequestID: prID,
		AuthorID:      authorID,
		Type:          "human",
		Path:          req.Path,
		Line:          req.Line,
		DiffHunk:      req.DiffHunk,
		Body:          req.Body,
		Suggestion:    req.Suggestion,
	}

	if err := s.db.Create(comment).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

// GetReviewComments 获取审查评论列表
func (s *prService) GetReviewComments(ctx context.Context, prID int) ([]*domain.ReviewComment, error) {
	var comments []*domain.ReviewComment
	result := s.db.Where("pull_request_id = ?", prID).Order("created_at ASC").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

// CheckMergeability 检查PR是否可合并
func (s *prService) CheckMergeability(ctx context.Context, prID int) (bool, error) {
	var pr domain.PullRequest
	if err := s.db.First(&pr, prID).Error; err != nil {
		return false, err
	}

	if pr.State != "open" {
		return false, nil
	}

	if pr.IsDraft {
		return false, nil
	}

	// 检查是否有批准的审查
	var approveCount int64
	s.db.Model(&domain.Review{}).Where("pull_request_id = ? AND state = ?", prID, "approve").Count(&approveCount)
	
	if approveCount == 0 {
		return false, nil
	}

	return true, nil
}

// UpdateMergeState 更新合并状态
func (s *prService) UpdateMergeState(ctx context.Context, prID int) error {
	mergeable, err := s.CheckMergeability(ctx, prID)
	if err != nil {
		return err
	}

	var pr domain.PullRequest
	if err := s.db.First(&pr, prID).Error; err != nil {
		return err
	}

	if mergeable {
		pr.MergeState = "clean"
	} else {
		pr.MergeState = "checking"
	}

	return s.db.Save(&pr).Error
}
