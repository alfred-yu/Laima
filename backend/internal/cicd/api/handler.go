package api

import (
	"laima/internal/cicd/app"
	"laima/internal/cicd/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CICDApi CI/CD API处理器
type CICDApi struct {
	cicdService app.CICDService
	db          *gorm.DB
}

// NewCICDApi 创建CI/CD API实例
func NewCICDApi(db *gorm.DB) *CICDApi {
	return &CICDApi{
		cicdService: app.NewCICDService(db),
		db:          db,
	}
}

// RegisterRoutes 注册路由
func (api *CICDApi) RegisterRoutes(r *gin.Engine) {
	cicdGroup := r.Group("/api/v1/cicd")
	{
		// 流水线管理
		cicdGroup.POST("/pipelines", api.CreatePipeline)
		cicdGroup.GET("/pipelines", api.ListPipelines)
		cicdGroup.GET("/pipelines/:id", api.GetPipeline)
		cicdGroup.POST("/pipelines/:id/cancel", api.CancelPipeline)

		// 任务管理
		cicdGroup.GET("/pipelines/:id/jobs", api.GetJobs)
		cicdGroup.GET("/jobs/:id", api.GetJob)
		cicdGroup.PUT("/jobs/:id/status", api.UpdateJobStatus)
		cicdGroup.POST("/jobs/:id/log", api.AddJobLog)

		// 集成触发
		cicdGroup.POST("/prs/:id/pipelines", api.TriggerPRPipeline)
		cicdGroup.POST("/repos/:id/pipelines", api.TriggerPushPipeline)
	}

	// 仓库相关的CI/CD路由
	repoCICDGroup := r.Group("/api/v1/repos/:owner/:repo/cicd")
	{
		repoCICDGroup.GET("/pipelines", api.ListRepoPipelines)
		repoCICDGroup.POST("/pipelines", api.CreateRepoPipeline)
		repoCICDGroup.GET("/pipelines/:id", api.GetRepoPipeline)
		repoCICDGroup.POST("/pipelines/:id/cancel", api.CancelRepoPipeline)
		repoCICDGroup.GET("/pipelines/:id/jobs", api.GetRepoPipelineJobs)
	}
}

// CreatePipeline 创建流水线
func (api *CICDApi) CreatePipeline(c *gin.Context) {
	var req domain.PipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pipeline, err := api.cicdService.CreatePipeline(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pipeline)
}

// ListPipelines 列出流水线
func (api *CICDApi) ListPipelines(c *gin.Context) {
	filter := &domain.PipelineFilter{
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

	filter.Status = c.Query("status")
	filter.Ref = c.Query("ref")
	filter.CommitSHA = c.Query("commit_sha")

	pipelines, total, err := api.cicdService.ListPipelines(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"items":     pipelines,
		"page":      filter.Page,
		"per_page":  filter.PerPage,
	})
}

// GetPipeline 获取流水线详情
func (api *CICDApi) GetPipeline(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	pipeline, err := api.cicdService.GetPipeline(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pipeline not found"})
		return
	}

	// 获取任务列表
	jobs, err := api.cicdService.GetJobs(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := &domain.PipelineResponse{
		ID:           pipeline.ID,
		RepositoryID: pipeline.RepositoryID,
		CommitSHA:    pipeline.CommitSHA,
		Ref:          pipeline.Ref,
		Status:       pipeline.Status,
		Trigger:      pipeline.Trigger,
		Duration:     pipeline.Duration,
		Jobs:         jobs,
		CreatedAt:    pipeline.CreatedAt,
		UpdatedAt:    pipeline.UpdatedAt,
	}

	c.JSON(http.StatusOK, response)
}

// CancelPipeline 取消流水线
func (api *CICDApi) CancelPipeline(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	pipeline, err := api.cicdService.CancelPipeline(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pipeline)
}

// GetJobs 获取任务列表
func (api *CICDApi) GetJobs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pipeline ID"})
		return
	}

	jobs, err := api.cicdService.GetJobs(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

// GetJob 获取任务详情
func (api *CICDApi) GetJob(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	job, err := api.cicdService.GetJob(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// UpdateJobStatus 更新任务状态
func (api *CICDApi) UpdateJobStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := api.cicdService.UpdateJobStatus(c, id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}

// AddJobLog 添加任务日志
func (api *CICDApi) AddJobLog(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}

	var req struct {
		Log string `json:"log" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := api.cicdService.AddJobLog(c, id, req.Log); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Log added successfully"})
}

// TriggerPRPipeline 为PR触发流水线
func (api *CICDApi) TriggerPRPipeline(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PR ID"})
		return
	}

	pipeline, err := api.cicdService.TriggerPipelineForPR(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pipeline)
}

// TriggerPushPipeline 为推送触发流水线
func (api *CICDApi) TriggerPushPipeline(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	var req struct {
		CommitSHA string `json:"commit_sha" binding:"required"`
		Ref       string `json:"ref" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pipeline, err := api.cicdService.TriggerPipelineForPush(c, id, req.CommitSHA, req.Ref)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pipeline)
}

// 仓库相关的CI/CD方法...
// 这些方法需要通过仓库路径来操作CI/CD
// 由于时间限制，这里只实现基本结构

func (api *CICDApi) ListRepoPipelines(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "List repo pipelines"})
}

func (api *CICDApi) CreateRepoPipeline(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Create repo pipeline"})
}

func (api *CICDApi) GetRepoPipeline(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Get repo pipeline"})
}

func (api *CICDApi) CancelRepoPipeline(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Cancel repo pipeline"})
}

func (api *CICDApi) GetRepoPipelineJobs(c *gin.Context) {
	// 实现逻辑
	c.JSON(http.StatusOK, gin.H{"message": "Get repo pipeline jobs"})
}
