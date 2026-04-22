package app

import (
	"context"
	"laima/internal/issue/domain"

	"gorm.io/gorm"
)

// IssueService Issue服务接口
type IssueService interface {
	// Issue CRUD
	CreateIssue(ctx context.Context, req *domain.CreateIssueRequest, authorID int) (*domain.Issue, error)
	GetIssue(ctx context.Context, issueID int) (*domain.Issue, error)
	GetIssueByNumber(ctx context.Context, repoID int, number int) (*domain.Issue, error)
	UpdateIssue(ctx context.Context, issueID int, req *domain.UpdateIssueRequest) (*domain.Issue, error)
	DeleteIssue(ctx context.Context, issueID int) error
	ListIssues(ctx context.Context, filter *domain.IssueFilter) ([]*domain.Issue, int64, error)

	// Issue 操作
	CloseIssue(ctx context.Context, issueID int, userID int) (*domain.Issue, error)
	ReopenIssue(ctx context.Context, issueID int, userID int) (*domain.Issue, error)
	AssignIssue(ctx context.Context, issueID int, assigneeID int) (*domain.Issue, error)

	// 评论管理
	CreateComment(ctx context.Context, issueID int, req *domain.IssueCommentRequest, authorID int) (*domain.IssueComment, error)
	GetComments(ctx context.Context, issueID int) ([]*domain.IssueComment, error)
	UpdateComment(ctx context.Context, commentID int, body string) (*domain.IssueComment, error)
	DeleteComment(ctx context.Context, commentID int) error

	// 里程碑管理
	CreateMilestone(ctx context.Context, repoID int, req *domain.MilestoneRequest) (*domain.Milestone, error)
	GetMilestone(ctx context.Context, milestoneID int) (*domain.Milestone, error)
	UpdateMilestone(ctx context.Context, milestoneID int, req *domain.MilestoneRequest) (*domain.Milestone, error)
	DeleteMilestone(ctx context.Context, milestoneID int) error
	ListMilestones(ctx context.Context, repoID int) ([]*domain.Milestone, error)
}

// issueService Issue服务实现
type issueService struct {
	db *gorm.DB
}

// NewIssueService 创建Issue服务实例
func NewIssueService(db *gorm.DB) IssueService {
	return &issueService{db: db}
}

// CreateIssue 创建Issue
func (s *issueService) CreateIssue(ctx context.Context, req *domain.CreateIssueRequest, authorID int) (*domain.Issue, error) {
	// 实现创建Issue逻辑
	// 1. 验证请求参数
	// 2. 生成Issue编号
	// 3. 创建Issue记录
	// 4. 返回Issue信息
	return nil, nil
}

// GetIssue 根据ID获取Issue
func (s *issueService) GetIssue(ctx context.Context, issueID int) (*domain.Issue, error) {
	var issue domain.Issue
	result := s.db.First(&issue, issueID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &issue, nil
}

// GetIssueByNumber 根据编号获取Issue
func (s *issueService) GetIssueByNumber(ctx context.Context, repoID int, number int) (*domain.Issue, error) {
	var issue domain.Issue
	result := s.db.Where("repository_id = ? AND number = ?", repoID, number).First(&issue)
	if result.Error != nil {
		return nil, result.Error
	}
	return &issue, nil
}

// UpdateIssue 更新Issue
func (s *issueService) UpdateIssue(ctx context.Context, issueID int, req *domain.UpdateIssueRequest) (*domain.Issue, error) {
	var issue domain.Issue
	if err := s.db.First(&issue, issueID).Error; err != nil {
		return nil, err
	}

	// 更新Issue信息
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.AssigneeID > 0 {
		updates["assignee_id"] = req.AssigneeID
	}
	if req.State != "" {
		updates["state"] = req.State
	}
	if req.MilestoneID > 0 {
		updates["milestone_id"] = req.MilestoneID
	}
	if req.Labels != nil {
		// 处理标签
	}
	if req.Priority != "" {
		updates["priority"] = req.Priority
	}

	if len(updates) > 0 {
		if err := s.db.Model(&issue).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return &issue, nil
}

// DeleteIssue 删除Issue
func (s *issueService) DeleteIssue(ctx context.Context, issueID int) error {
	return s.db.Delete(&domain.Issue{}, issueID).Error
}

// ListIssues 列出Issue
func (s *issueService) ListIssues(ctx context.Context, filter *domain.IssueFilter) ([]*domain.Issue, int64, error) {
	var issues []*domain.Issue
	var total int64

	query := s.db.Model(&domain.Issue{})

	// 应用过滤条件
	if filter.RepositoryID > 0 {
		query = query.Where("repository_id = ?", filter.RepositoryID)
	}
	if filter.AuthorID > 0 {
		query = query.Where("author_id = ?", filter.AuthorID)
	}
	if filter.AssigneeID > 0 {
		query = query.Where("assignee_id = ?", filter.AssigneeID)
	}
	if filter.State != "" {
		query = query.Where("state = ?", filter.State)
	}
	if filter.MilestoneID > 0 {
		query = query.Where("milestone_id = ?", filter.MilestoneID)
	}
	if filter.Priority != "" {
		query = query.Where("priority = ?", filter.Priority)
	}
	if filter.Search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	// 计算总数
	query.Count(&total)

	// 分页
	offset := (filter.Page - 1) * filter.PerPage
	result := query.Offset(offset).Limit(filter.PerPage).Order("created_at DESC").Find(&issues)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return issues, total, nil
}

// CloseIssue 关闭Issue
func (s *issueService) CloseIssue(ctx context.Context, issueID int, userID int) (*domain.Issue, error) {
	var issue domain.Issue
	if err := s.db.First(&issue, issueID).Error; err != nil {
		return nil, err
	}

	issue.State = domain.IssueStatusClosed
	if err := s.db.Save(&issue).Error; err != nil {
		return nil, err
	}

	return &issue, nil
}

// ReopenIssue 重新打开Issue
func (s *issueService) ReopenIssue(ctx context.Context, issueID int, userID int) (*domain.Issue, error) {
	var issue domain.Issue
	if err := s.db.First(&issue, issueID).Error; err != nil {
		return nil, err
	}

	issue.State = domain.IssueStatusOpen
	if err := s.db.Save(&issue).Error; err != nil {
		return nil, err
	}

	return &issue, nil
}

// AssignIssue 分配Issue
func (s *issueService) AssignIssue(ctx context.Context, issueID int, assigneeID int) (*domain.Issue, error) {
	var issue domain.Issue
	if err := s.db.First(&issue, issueID).Error; err != nil {
		return nil, err
	}

	issue.AssigneeID = assigneeID
	if err := s.db.Save(&issue).Error; err != nil {
		return nil, err
	}

	return &issue, nil
}

// CreateComment 创建评论
func (s *issueService) CreateComment(ctx context.Context, issueID int, req *domain.IssueCommentRequest, authorID int) (*domain.IssueComment, error) {
	// 实现创建评论逻辑
	// 1. 验证Issue存在
	// 2. 创建评论记录
	// 3. 返回评论信息
	return nil, nil
}

// GetComments 获取评论列表
func (s *issueService) GetComments(ctx context.Context, issueID int) ([]*domain.IssueComment, error) {
	var comments []*domain.IssueComment
	result := s.db.Where("issue_id = ?", issueID).Order("created_at ASC").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

// UpdateComment 更新评论
func (s *issueService) UpdateComment(ctx context.Context, commentID int, body string) (*domain.IssueComment, error) {
	var comment domain.IssueComment
	if err := s.db.First(&comment, commentID).Error; err != nil {
		return nil, err
	}

	comment.Body = body
	if err := s.db.Save(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

// DeleteComment 删除评论
func (s *issueService) DeleteComment(ctx context.Context, commentID int) error {
	return s.db.Delete(&domain.IssueComment{}, commentID).Error
}

// CreateMilestone 创建里程碑
func (s *issueService) CreateMilestone(ctx context.Context, repoID int, req *domain.MilestoneRequest) (*domain.Milestone, error) {
	// 实现创建里程碑逻辑
	// 1. 验证请求参数
	// 2. 创建里程碑记录
	// 3. 返回里程碑信息
	return nil, nil
}

// GetMilestone 根据ID获取里程碑
func (s *issueService) GetMilestone(ctx context.Context, milestoneID int) (*domain.Milestone, error) {
	var milestone domain.Milestone
	result := s.db.First(&milestone, milestoneID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &milestone, nil
}

// UpdateMilestone 更新里程碑
func (s *issueService) UpdateMilestone(ctx context.Context, milestoneID int, req *domain.MilestoneRequest) (*domain.Milestone, error) {
	var milestone domain.Milestone
	if err := s.db.First(&milestone, milestoneID).Error; err != nil {
		return nil, err
	}

	// 更新里程碑信息
	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if !req.DueDate.IsZero() {
		updates["due_date"] = req.DueDate
	}
	if req.State != "" {
		updates["state"] = req.State
	}

	if len(updates) > 0 {
		if err := s.db.Model(&milestone).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return &milestone, nil
}

// DeleteMilestone 删除里程碑
func (s *issueService) DeleteMilestone(ctx context.Context, milestoneID int) error {
	return s.db.Delete(&domain.Milestone{}, milestoneID).Error
}

// ListMilestones 列出里程碑
func (s *issueService) ListMilestones(ctx context.Context, repoID int) ([]*domain.Milestone, error) {
	var milestones []*domain.Milestone
	result := s.db.Where("repository_id = ?", repoID).Order("created_at DESC").Find(&milestones)
	if result.Error != nil {
		return nil, result.Error
	}
	return milestones, nil
}
