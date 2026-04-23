package api

import (
	"laima/internal/ai/app"
	"laima/internal/ai/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AIApi AI API处理器
type AIApi struct {
	aiService app.AIService
	db        *gorm.DB
}

// NewAIApi 创建AI API实例
func NewAIApi(db *gorm.DB, aiService app.AIService) *AIApi {
	return &AIApi{
		aiService: aiService,
		db:        db,
	}
}

// RegisterRoutes 注册路由
func (api *AIApi) RegisterRoutes(r *gin.Engine) {
	aiGroup := r.Group("/api/v1/ai")
	{
		// AI审查管理
		aiGroup.POST("/reviews", api.TriggerAIReview)
		aiGroup.GET("/reviews", api.ListAIReviews)
		aiGroup.GET("/reviews/:id", api.GetAIReview)
		aiGroup.GET("/reviews/:id/issues", api.GetAIReviewIssues)

		// PR相关的AI审查
		aiGroup.POST("/prs/:id/reviews", api.TriggerPRReview)
		aiGroup.GET("/prs/:id/reviews", api.GetPRReview)
		aiGroup.GET("/prs/:id/issues", api.GetPRReviewIssues)
	}
}

// TriggerAIReview 触发AI审查
func (api *AIApi) TriggerAIReview(c *gin.Context) {
	var req domain.AIReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := api.aiService.TriggerAIReview(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// ListAIReviews 列出AI审查
func (api *AIApi) ListAIReviews(c *gin.Context) {
	filter := &domain.AIReviewFilter{
		Page:    1,
		PerPage: 30,
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		filter.Page = page
	}

	if perPage, err := strconv.Atoi(c.Query("per_page")); err == nil && perPage > 0 {
		filter.PerPage = perPage
	}

	if prID, err := strconv.Atoi(c.Query("pull_request_id")); err == nil && prID > 0 {
		filter.PullRequestID = prID
	}

	if repoID, err := strconv.Atoi(c.Query("repository_id")); err == nil && repoID > 0 {
		filter.RepositoryID = repoID
	}

	filter.Status = c.Query("status")

	reviews, total, err := api.aiService.ListAIReviews(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"items":     reviews,
		"page":      filter.Page,
		"per_page":  filter.PerPage,
	})
}

// GetAIReview 获取AI审查详情
func (api *AIApi) GetAIReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	review, err := api.aiService.GetAIReview(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	// 获取审查问题
	issues, err := api.aiService.GetAIReviewIssues(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := &domain.AIReviewResponse{
		ID:            review.ID,
		PullRequestID: review.PullRequestID,
		RepositoryID:  review.RepositoryID,
		Status:        review.Status,
		Score:         review.Score,
		Summary:       review.Summary,
		Issues:        issues,
		StartedAt:     review.StartedAt,
		CompletedAt:   review.CompletedAt,
		Error:         review.Error,
		CreatedAt:     review.CreatedAt,
		UpdatedAt:     review.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// GetAIReviewIssues 获取AI审查问题
func (api *AIApi) GetAIReviewIssues(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	issues, err := api.aiService.GetAIReviewIssues(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, issues)
}

// TriggerPRReview 为PR触发AI审查
func (api *AIApi) TriggerPRReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	review, err := api.aiService.TriggerAIReviewForPR(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// GetPRReview 获取PR的AI审查
func (api *AIApi) GetPRReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	review, err := api.aiService.GetAIReviewByPR(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	// 获取审查问题
	issues, err := api.aiService.GetAIReviewIssues(c, review.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := &domain.AIReviewResponse{
		ID:            review.ID,
		PullRequestID: review.PullRequestID,
		RepositoryID:  review.RepositoryID,
		Status:        review.Status,
		Score:         review.Score,
		Summary:       review.Summary,
		Issues:        issues,
		StartedAt:     review.StartedAt,
		CompletedAt:   review.CompletedAt,
		Error:         review.Error,
		CreatedAt:     review.CreatedAt,
		UpdatedAt:     review.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// GetPRReviewIssues 获取PR的AI审查问题
func (api *AIApi) GetPRReviewIssues(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	issues, err := api.aiService.GetAIReviewIssuesByPR(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, issues)
}
