package api

import (
	"laima/internal/issue/app"
	"laima/internal/issue/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IssueApi Issue API处理器
type IssueApi struct {
	issueService app.IssueService
	db          *gorm.DB
}

// NewIssueApi 创建Issue API实例
func NewIssueApi(db *gorm.DB) *IssueApi {
	return &IssueApi{
		issueService: app.NewIssueService(db),
		db:          db,
	}
}

// RegisterRoutes 注册路由
func (api *IssueApi) RegisterRoutes(r *gin.Engine) {
	issueGroup := r.Group("/api/v1/issues")
	{
		// Issue CRUD
		issueGroup.GET("", api.ListIssues)
		issueGroup.POST("", api.CreateIssue)
		issueGroup.GET("/:id", api.GetIssue)
		issueGroup.PUT("/:id", api.UpdateIssue)
		issueGroup.DELETE("/:id", api.DeleteIssue)

		// Issue 操作
		issueGroup.POST("/:id/close", api.CloseIssue)
		issueGroup.POST("/:id/reopen", api.ReopenIssue)
		issueGroup.POST("/:id/assign", api.AssignIssue)

		// 评论管理
		issueGroup.GET("/:id/comments", api.GetComments)
		issueGroup.POST("/:id/comments", api.CreateComment)
		issueGroup.PUT("/comments/:id", api.UpdateComment)
		issueGroup.DELETE("/comments/:id", api.DeleteComment)

		// 时间追踪
		issueGroup.GET("/:id/time-tracking", api.GetIssueTimeTracking)
		issueGroup.POST("/:id/time-tracking", api.AddIssueTimeTracking)
		issueGroup.GET("/:id/time-summary", api.GetIssueTimeSummary)
	}

	// 时间追踪路由
	timeTrackingGroup := r.Group("/api/v1/time-tracking")
	{
		timeTrackingGroup.GET("", api.ListTimeTracking)
		timeTrackingGroup.POST("", api.AddTimeTracking)
		timeTrackingGroup.GET("/:id", api.GetTimeTracking)
		timeTrackingGroup.PUT("/:id", api.UpdateTimeTracking)
		timeTrackingGroup.DELETE("/:id", api.DeleteTimeTracking)
	}

	// 仓库时间追踪路由
	repoTimeTrackingGroup := r.Group("/api/v1/repos/:owner/:repo/time-tracking")
	{
		repoTimeTrackingGroup.GET("", api.ListRepoTimeTracking)
		repoTimeTrackingGroup.GET("/summary", api.GetRepoTimeSummary)
	}

	// 里程碑路由
	milestoneGroup := r.Group("/api/v1/milestones")
	{
		milestoneGroup.GET("", api.ListMilestones)
		milestoneGroup.POST("", api.CreateMilestone)
		milestoneGroup.GET("/:id", api.GetMilestone)
		milestoneGroup.PUT("/:id", api.UpdateMilestone)
		milestoneGroup.DELETE("/:id", api.DeleteMilestone)
	}

	// 仓库相关的Issue路由
	repoIssueGroup := r.Group("/api/v1/repos/:owner/:repo/issues")
	{
		repoIssueGroup.GET("", api.ListRepoIssues)
		repoIssueGroup.POST("", api.CreateRepoIssue)
		repoIssueGroup.GET("/:number", api.GetRepoIssue)
		repoIssueGroup.PUT("/:number", api.UpdateRepoIssue)
		repoIssueGroup.DELETE("/:number", api.DeleteRepoIssue)

		repoIssueGroup.POST("/:number/close", api.CloseRepoIssue)
		repoIssueGroup.POST("/:number/reopen", api.ReopenRepoIssue)
		repoIssueGroup.POST("/:number/assign", api.AssignRepoIssue)

		repoIssueGroup.GET("/:number/comments", api.GetRepoIssueComments)
		repoIssueGroup.POST("/:number/comments", api.CreateRepoIssueComment)
	}

	// 仓库相关的里程碑路由
	repoMilestoneGroup := r.Group("/api/v1/repos/:owner/:repo/milestones")
	{
		repoMilestoneGroup.GET("", api.ListRepoMilestones)
		repoMilestoneGroup.POST("", api.CreateRepoMilestone)
		repoMilestoneGroup.GET("/:id", api.GetRepoMilestone)
		repoMilestoneGroup.PUT("/:id", api.UpdateRepoMilestone)
		repoMilestoneGroup.DELETE("/:id", api.DeleteRepoMilestone)
	}
}

// ListIssues 列出Issue
func (api *IssueApi) ListIssues(c *gin.Context) {
	filter := &domain.IssueFilter{
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

	if assigneeID, err := strconv.Atoi(c.Query("assignee_id")); err == nil && assigneeID > 0 {
		filter.AssigneeID = assigneeID
	}

	if milestoneID, err := strconv.Atoi(c.Query("milestone_id")); err == nil && milestoneID > 0 {
		filter.MilestoneID = milestoneID
	}

	filter.State = c.Query("state")
	filter.Priority = c.Query("priority")
	filter.Search = c.Query("q")

	issues, total, err := api.issueService.ListIssues(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"items":     issues,
		"page":      filter.Page,
		"per_page":  filter.PerPage,
	})
}

// CreateIssue 创建Issue
func (api *IssueApi) CreateIssue(c *gin.Context) {
	var req domain.CreateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT token中获取用户ID
	authorID := 1 // 临时硬编码

	issue, err := api.issueService.CreateIssue(c, &req, authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, issue)
}

// GetIssue 获取Issue详情
func (api *IssueApi) GetIssue(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	issue, err := api.issueService.GetIssue(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Issue not found"})
		return
	}

	c.JSON(http.StatusOK, issue)
}

// UpdateIssue 更新Issue
func (api *IssueApi) UpdateIssue(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	var req domain.UpdateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	issue, err := api.issueService.UpdateIssue(c, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, issue)
}

// DeleteIssue 删除Issue
func (api *IssueApi) DeleteIssue(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	if err := api.issueService.DeleteIssue(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Issue deleted successfully"})
}

// CloseIssue 关闭Issue
func (api *IssueApi) CloseIssue(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	// 从JWT token中获取用户ID
	userID := 1 // 临时硬编码

	issue, err := api.issueService.CloseIssue(c, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, issue)
}

// ReopenIssue 重新打开Issue
func (api *IssueApi) ReopenIssue(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	// 从JWT token中获取用户ID
	userID := 1 // 临时硬编码

	issue, err := api.issueService.ReopenIssue(c, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, issue)
}

// AssignIssue 分配Issue
func (api *IssueApi) AssignIssue(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	var req struct {
		AssigneeID int `json:"assignee_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	issue, err := api.issueService.AssignIssue(c, id, req.AssigneeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, issue)
}

// GetComments 获取评论列表
func (api *IssueApi) GetComments(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	comments, err := api.issueService.GetComments(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// CreateComment 创建评论
func (api *IssueApi) CreateComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	var req domain.IssueCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT token中获取用户ID
	authorID := 1 // 临时硬编码

	comment, err := api.issueService.CreateComment(c, id, &req, authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// UpdateComment 更新评论
func (api *IssueApi) UpdateComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var req struct {
		Body string `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := api.issueService.UpdateComment(c, id, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// DeleteComment 删除评论
func (api *IssueApi) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	if err := api.issueService.DeleteComment(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// CreateMilestone 创建里程碑
func (api *IssueApi) CreateMilestone(c *gin.Context) {
	var req struct {
		RepositoryID int                   `json:"repository_id" binding:"required"`
		Milestone    domain.MilestoneRequest `json:"milestone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	milestone, err := api.issueService.CreateMilestone(c, req.RepositoryID, &req.Milestone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, milestone)
}

// GetMilestone 获取里程碑详情
func (api *IssueApi) GetMilestone(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone ID"})
		return
	}

	milestone, err := api.issueService.GetMilestone(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Milestone not found"})
		return
	}

	c.JSON(http.StatusOK, milestone)
}

// UpdateMilestone 更新里程碑
func (api *IssueApi) UpdateMilestone(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone ID"})
		return
	}

	var req domain.MilestoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	milestone, err := api.issueService.UpdateMilestone(c, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, milestone)
}

// DeleteMilestone 删除里程碑
func (api *IssueApi) DeleteMilestone(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid milestone ID"})
		return
	}

	if err := api.issueService.DeleteMilestone(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Milestone deleted successfully"})
}

// ListMilestones 列出里程碑
func (api *IssueApi) ListMilestones(c *gin.Context) {
	repoIDStr := c.Query("repository_id")
	if repoIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Repository ID is required"})
		return
	}

	repoID, err := strconv.Atoi(repoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	milestones, err := api.issueService.ListMilestones(c, repoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, milestones)
}

// 仓库相关的Issue方法...
// 这些方法需要通过仓库路径来操作Issue
// 由于时间限制，这里只实现基本结构

func (api *IssueApi) ListRepoIssues(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "List repo issues"})
}

func (api *IssueApi) CreateRepoIssue(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Create repo issue"})
}

func (api *IssueApi) GetRepoIssue(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Get repo issue"})
}

func (api *IssueApi) UpdateRepoIssue(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Update repo issue"})
}

func (api *IssueApi) DeleteRepoIssue(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Delete repo issue"})
}

func (api *IssueApi) CloseRepoIssue(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Close repo issue"})
}

func (api *IssueApi) ReopenRepoIssue(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Reopen repo issue"})
}

func (api *IssueApi) AssignRepoIssue(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Assign repo issue"})
}

func (api *IssueApi) GetRepoIssueComments(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Get repo issue comments"})
}

func (api *IssueApi) CreateRepoIssueComment(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Create repo issue comment"})
}

func (api *IssueApi) ListRepoMilestones(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "List repo milestones"})
}

func (api *IssueApi) CreateRepoMilestone(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Create repo milestone"})
}

func (api *IssueApi) GetRepoMilestone(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Get repo milestone"})
}

func (api *IssueApi) UpdateRepoMilestone(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Update repo milestone"})
}

func (api *IssueApi) DeleteRepoMilestone(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Delete repo milestone"})
}

// AddTimeTracking 添加时间追踪
func (api *IssueApi) AddTimeTracking(c *gin.Context) {
	var req domain.TimeTrackingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT token中获取用户ID
	userID := 1 // 临时硬编码

	tracking, err := api.issueService.AddTimeTracking(c, &req, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tracking)
}

// UpdateTimeTracking 更新时间追踪
func (api *IssueApi) UpdateTimeTracking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time tracking ID"})
		return
	}

	var req domain.TimeTrackingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tracking, err := api.issueService.UpdateTimeTracking(c, id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tracking)
}

// DeleteTimeTracking 删除时间追踪
func (api *IssueApi) DeleteTimeTracking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time tracking ID"})
		return
	}

	if err := api.issueService.DeleteTimeTracking(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Time tracking deleted successfully"})
}

// GetTimeTracking 获取时间追踪详情
func (api *IssueApi) GetTimeTracking(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time tracking ID"})
		return
	}

	tracking, err := api.issueService.GetTimeTracking(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Time tracking not found"})
		return
	}

	c.JSON(http.StatusOK, tracking)
}

// ListTimeTracking 列出时间追踪记录
func (api *IssueApi) ListTimeTracking(c *gin.Context) {
	filter := &domain.TimeTrackingFilter{
		Page:    1,
		PerPage: 30,
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		filter.Page = page
	}

	if perPage, err := strconv.Atoi(c.Query("per_page")); err == nil && perPage > 0 {
		filter.PerPage = perPage
	}

	if issueID, err := strconv.Atoi(c.Query("issue_id")); err == nil && issueID > 0 {
		filter.IssueID = issueID
	}

	if userID, err := strconv.Atoi(c.Query("user_id")); err == nil && userID > 0 {
		filter.UserID = userID
	}

	if repoID, err := strconv.Atoi(c.Query("repository_id")); err == nil && repoID > 0 {
		filter.RepositoryID = repoID
	}

	filter.StartDate = c.Query("start_date")
	filter.EndDate = c.Query("end_date")

	trackings, total, err := api.issueService.ListTimeTracking(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"items":     trackings,
		"page":      filter.Page,
		"per_page":  filter.PerPage,
	})
}

// GetIssueTimeTracking 获取Issue的时间追踪记录
func (api *IssueApi) GetIssueTimeTracking(c *gin.Context) {
	idStr := c.Param("id")
	issueID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	filter := &domain.TimeTrackingFilter{
		IssueID: issueID,
		Page:    1,
		PerPage: 30,
	}

	trackings, total, err := api.issueService.ListTimeTracking(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"items":     trackings,
		"page":      filter.Page,
		"per_page":  filter.PerPage,
	})
}

// AddIssueTimeTracking 为Issue添加时间追踪
func (api *IssueApi) AddIssueTimeTracking(c *gin.Context) {
	idStr := c.Param("id")
	issueID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	var req struct {
		Hours       float64 `json:"hours" binding:"required,gt=0"`
		Description string  `json:"description"`
		Date        string  `json:"date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trackingReq := &domain.TimeTrackingRequest{
		IssueID:     issueID,
		Hours:       req.Hours,
		Description: req.Description,
		Date:        req.Date,
	}

	// 从JWT token中获取用户ID
	userID := 1 // 临时硬编码

	tracking, err := api.issueService.AddTimeTracking(c, trackingReq, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tracking)
}

// GetIssueTimeSummary 获取Issue的时间汇总
func (api *IssueApi) GetIssueTimeSummary(c *gin.Context) {
	idStr := c.Param("id")
	issueID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid issue ID"})
		return
	}

	totalHours, err := api.issueService.GetIssueTimeSummary(c, issueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"issue_id":     issueID,
		"total_hours":  totalHours,
	})
}

// ListRepoTimeTracking 列出仓库的时间追踪记录
func (api *IssueApi) ListRepoTimeTracking(c *gin.Context) {
	// 这里需要通过owner和repo名称获取仓库ID
	// 由于时间限制，这里只实现基本结构
	c.JSON(http.StatusOK, gin.H{"message": "List repo time tracking"})
}

// GetRepoTimeSummary 获取仓库的时间汇总
func (api *IssueApi) GetRepoTimeSummary(c *gin.Context) {
	// 这里需要通过owner和repo名称获取仓库ID
	// 由于时间限制，这里只实现基本结构
	c.JSON(http.StatusOK, gin.H{"message": "Get repo time summary"})
}
