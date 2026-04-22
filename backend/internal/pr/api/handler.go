package api

import (
	"laima/internal/git"
	"laima/internal/pr/app"
	"laima/internal/pr/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PRAPI PR API处理器
type PRAPI struct {
	prService app.PRService
	db        *gorm.DB
}

// NewPRAPI 创建PR API实例
func NewPRAPI(db *gorm.DB, gitSvc *git.Service) *PRAPI {
	return &PRAPI{
		prService: app.NewPRService(db, gitSvc),
		db:        db,
	}
}

// RegisterRoutes 注册路由
func (api *PRAPI) RegisterRoutes(r *gin.Engine) {
	prGroup := r.Group("/api/v1/prs")
	{
		// PR CRUD
		prGroup.GET("", api.ListPRs)
		prGroup.POST("", api.CreatePR)
		prGroup.GET("/:id", api.GetPR)
		prGroup.PUT("/:id", api.UpdatePR)
		prGroup.DELETE("/:id", api.DeletePR)

		// PR 操作
		prGroup.POST("/:id/merge", api.MergePR)
		prGroup.POST("/:id/close", api.ClosePR)
		prGroup.POST("/:id/reopen", api.ReopenPR)

		// 审查操作
		prGroup.GET("/:id/reviews", api.GetReviews)
		prGroup.POST("/:id/reviews", api.CreateReview)
		prGroup.GET("/:id/comments", api.GetReviewComments)
		prGroup.POST("/:id/comments", api.CreateReviewComment)

		// 状态检查
		prGroup.GET("/:id/mergeability", api.CheckMergeability)
	}

	// 仓库相关的PR路由
	repoPRGroup := r.Group("/api/v1/repos/:owner/:repo/pulls")
	{
		repoPRGroup.GET("", api.ListRepoPRs)
		repoPRGroup.POST("", api.CreateRepoPR)
		repoPRGroup.GET("/:number", api.GetRepoPR)
		repoPRGroup.PUT("/:number", api.UpdateRepoPR)
		repoPRGroup.DELETE("/:number", api.DeleteRepoPR)

		repoPRGroup.POST("/:number/merge", api.MergeRepoPR)
		repoPRGroup.POST("/:number/close", api.CloseRepoPR)
		repoPRGroup.POST("/:number/reopen", api.ReopenRepoPR)

		repoPRGroup.GET("/:number/reviews", api.GetRepoPRReviews)
		repoPRGroup.POST("/:number/reviews", api.CreateRepoPRReview)
		repoPRGroup.GET("/:number/comments", api.GetRepoPRReviewComments)
		repoPRGroup.POST("/:number/comments", api.CreateRepoPRReviewComment)

		repoPRGroup.GET("/:number/mergeability", api.CheckRepoPRMergeability)
	}
}

// ListPRs 列出PR
func (api *PRAPI) ListPRs(c *gin.Context) {
	filter := &domain.PRFilter{
		Page:    1,
		PerPage: 30,
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		filter.Page = page
	}

	if perPage, err := strconv.Atoi(c.Query("per_page")); err == nil && perPage > 0 {
		filter.PerPage = perPage
	}

	if repoID, err := strconv.Atoi(c.Query("repository_id")); err == nil && repoID > 0 {
		filter.RepositoryID = repoID
	}

	if authorID, err := strconv.Atoi(c.Query("author_id")); err == nil && authorID > 0 {
		filter.AuthorID = authorID
	}

	filter.State = c.Query("state")
	filter.Search = c.Query("q")

	prs, total, err := api.prService.ListPRs(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"items":     prs,
		"page":      filter.Page,
		"per_page":  filter.PerPage,
	})
}

// CreatePR 创建PR
func (api *PRAPI) CreatePR(c *gin.Context) {
	var req domain.CreatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT token中获取用户ID
	userID := 1 // 临时硬编码

	pr, err := api.prService.CreatePR(c, &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pr)
}

// GetPR 获取PR详情
func (api *PRAPI) GetPR(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	pr, err := api.prService.GetPR(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "PR not found"})
		return
	}

	c.JSON(http.StatusOK, pr)
}

// UpdatePR 更新PR
func (api *PRAPI) UpdatePR(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	var req domain.UpdatePRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pr, err := api.prService.UpdatePR(c, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pr)
}

// DeletePR 删除PR
func (api *PRAPI) DeletePR(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	if err := api.prService.DeletePR(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PR deleted successfully"})
}

// MergePR 合并PR
func (api *PRAPI) MergePR(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	// 从JWT token中获取用户ID
	userID := 1 // 临时硬编码

	pr, err := api.prService.MergePR(c, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pr)
}

// ClosePR 关闭PR
func (api *PRAPI) ClosePR(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	// 从JWT token中获取用户ID
	userID := 1 // 临时硬编码

	pr, err := api.prService.ClosePR(c, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pr)
}

// ReopenPR 重新打开PR
func (api *PRAPI) ReopenPR(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	// 从JWT token中获取用户ID
	userID := 1 // 临时硬编码

	pr, err := api.prService.ReopenPR(c, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pr)
}

// GetReviews 获取审查列表
func (api *PRAPI) GetReviews(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	reviews, err := api.prService.GetReviews(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// CreateReview 创建审查
func (api *PRAPI) CreateReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	var req domain.ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT token中获取用户ID
	reviewerID := 1 // 临时硬编码

	review, err := api.prService.CreateReview(c, id, &req, reviewerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// GetReviewComments 获取审查评论列表
func (api *PRAPI) GetReviewComments(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	comments, err := api.prService.GetReviewComments(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// CreateReviewComment 创建审查评论
func (api *PRAPI) CreateReviewComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	var req domain.ReviewCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT token中获取用户ID
	authorID := 1 // 临时硬编码

	comment, err := api.prService.CreateReviewComment(c, id, &req, authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// CheckMergeability 检查PR是否可合并
func (api *PRAPI) CheckMergeability(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	mergeable, err := api.prService.CheckMergeability(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mergeable": mergeable})
}

// 仓库相关的PR方法...
// 这些方法需要通过仓库路径和PR编号来操作PR
// 由于时间限制，这里只实现基本结构

func (api *PRAPI) ListRepoPRs(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "List repo PRs"})
}

func (api *PRAPI) CreateRepoPR(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Create repo PR"})
}

func (api *PRAPI) GetRepoPR(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Get repo PR"})
}

func (api *PRAPI) UpdateRepoPR(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Update repo PR"})
}

func (api *PRAPI) DeleteRepoPR(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Delete repo PR"})
}

func (api *PRAPI) MergeRepoPR(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Merge repo PR"})
}

func (api *PRAPI) CloseRepoPR(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Close repo PR"})
}

func (api *PRAPI) ReopenRepoPR(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Reopen repo PR"})
}

func (api *PRAPI) GetRepoPRReviews(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Get repo PR reviews"})
}

func (api *PRAPI) CreateRepoPRReview(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Create repo PR review"})
}

func (api *PRAPI) GetRepoPRReviewComments(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Get repo PR review comments"})
}

func (api *PRAPI) CreateRepoPRReviewComment(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Create repo PR review comment"})
}

func (api *PRAPI) CheckRepoPRMergeability(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Check repo PR mergeability"})
}
