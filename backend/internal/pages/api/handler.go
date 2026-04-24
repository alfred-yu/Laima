package api

import (
	"laima/internal/pages/app"
	"laima/internal/pages/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PagesApi Pages API处理器
type PagesApi struct {
	pagesService app.PagesService
	db          *gorm.DB
}

// NewPagesApi 创建Pages API实例
func NewPagesApi(db *gorm.DB) *PagesApi {
	return &PagesApi{
		pagesService: app.NewPagesService(db),
		db:          db,
	}
}

// RegisterRoutes 注册路由
func (api *PagesApi) RegisterRoutes(r *gin.Engine) {
	pagesGroup := r.Group("/api/v1/pages")
	{
		// Pages CRUD
		pagesGroup.GET("", api.ListPages)
		pagesGroup.POST("", api.CreatePages)
		pagesGroup.GET("/:id", api.GetPages)
		pagesGroup.PUT("/:id", api.UpdatePages)
		pagesGroup.DELETE("/:id", api.DeletePages)

		// Pages 操作
		pagesGroup.POST("/:id/publish", api.PublishPages)
		pagesGroup.POST("/:id/archive", api.ArchivePages)
		pagesGroup.POST("/:id/unpublish", api.UnpublishPages)

		// 配置管理
		pagesGroup.GET("/config", api.GetPagesConfig)
		pagesGroup.PUT("/config", api.UpdatePagesConfig)

		// 静态站点
		pagesGroup.POST("/generate", api.GenerateStaticSite)
		pagesGroup.GET("/url", api.GetStaticSiteURL)
	}

	// 仓库相关的Pages路由
	repoPagesGroup := r.Group("/api/v1/repos/:owner/:repo/pages")
	{
		repoPagesGroup.GET("", api.ListRepoPages)
		repoPagesGroup.POST("", api.CreateRepoPages)
		repoPagesGroup.GET("/:slug", api.GetRepoPages)
		repoPagesGroup.PUT("/:slug", api.UpdateRepoPages)
		repoPagesGroup.DELETE("/:slug", api.DeleteRepoPages)

		repoPagesGroup.POST("/:slug/publish", api.PublishRepoPages)
		repoPagesGroup.POST("/:slug/archive", api.ArchiveRepoPages)
		repoPagesGroup.POST("/:slug/unpublish", api.UnpublishRepoPages)

		repoPagesGroup.GET("/config", api.GetRepoPagesConfig)
		repoPagesGroup.PUT("/config", api.UpdateRepoPagesConfig)

		repoPagesGroup.POST("/generate", api.GenerateRepoStaticSite)
		repoPagesGroup.GET("/url", api.GetRepoStaticSiteURL)
	}

	// 公开访问的Pages路由
	publicPagesGroup := r.Group("/pages")
	{
		publicPagesGroup.GET("/:owner/:repo/*slug", api.GetPublicPages)
	}
}

// ListPages 列出Pages
func (api *PagesApi) ListPages(c *gin.Context) {
	filter := &domain.PagesFilter{
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

	if status := c.Query("status"); status != "" {
		filter.Status = status
	}

	if authorID, err := strconv.Atoi(c.Query("author_id")); err == nil && authorID > 0 {
		filter.AuthorID = authorID
	}

	filter.Search = c.Query("q")

	pages, total, err := api.pagesService.ListPages(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":     total,
		"items":     pages,
		"page":      filter.Page,
		"per_page":  filter.PerPage,
	})
}

// CreatePages 创建Pages
func (api *PagesApi) CreatePages(c *gin.Context) {
	var req domain.CreatePagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT token中获取用户ID
	authorID := 1 // 临时硬编码

	pages, err := api.pagesService.CreatePages(c, &req, authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pages)
}

// GetPages 获取Pages详情
func (api *PagesApi) GetPages(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pages ID"})
		return
	}

	pages, err := api.pagesService.GetPages(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pages not found"})
		return
	}

	c.JSON(http.StatusOK, pages)
}

// UpdatePages 更新Pages
func (api *PagesApi) UpdatePages(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pages ID"})
		return
	}

	var req domain.UpdatePagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从JWT token中获取用户ID
	editorID := 1 // 临时硬编码

	pages, err := api.pagesService.UpdatePages(c, id, &req, editorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pages)
}

// DeletePages 删除Pages
func (api *PagesApi) DeletePages(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pages ID"})
		return
	}

	if err := api.pagesService.DeletePages(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pages deleted successfully"})
}

// PublishPages 发布Pages
func (api *PagesApi) PublishPages(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pages ID"})
		return
	}

	// 从JWT token中获取用户ID
	userID := 1 // 临时硬编码

	pages, err := api.pagesService.PublishPages(c, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pages)
}

// ArchivePages 归档Pages
func (api *PagesApi) ArchivePages(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pages ID"})
		return
	}

	// 从JWT token中获取用户ID
	userID := 1 // 临时硬编码

	pages, err := api.pagesService.ArchivePages(c, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pages)
}

// UnpublishPages 取消发布Pages
func (api *PagesApi) UnpublishPages(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pages ID"})
		return
	}

	// 从JWT token中获取用户ID
	userID := 1 // 临时硬编码

	pages, err := api.pagesService.UnpublishPages(c, id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pages)
}

// GetPagesConfig 获取Pages配置
func (api *PagesApi) GetPagesConfig(c *gin.Context) {
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

	config, err := api.pagesService.GetPagesConfig(c, repoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdatePagesConfig 更新Pages配置
func (api *PagesApi) UpdatePagesConfig(c *gin.Context) {
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

	var req domain.UpdatePagesConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config, err := api.pagesService.UpdatePagesConfig(c, repoID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// GenerateStaticSite 生成静态站点
func (api *PagesApi) GenerateStaticSite(c *gin.Context) {
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

	path, err := api.pagesService.GenerateStaticSite(c, repoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path":    path,
		"message": "Static site generated successfully",
	})
}

// GetStaticSiteURL 获取静态站点URL
func (api *PagesApi) GetStaticSiteURL(c *gin.Context) {
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

	url, err := api.pagesService.GetStaticSiteURL(c, repoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

// GetPublicPages 获取公开访问的Pages
func (api *PagesApi) GetPublicPages(c *gin.Context) {
	// 这里需要通过owner和repo名称获取仓库ID
	// 然后根据slug获取Pages
	// 由于时间限制，这里只实现基本结构
	c.JSON(http.StatusOK, gin.H{"message": "Public pages"})
}

// 仓库相关的Pages方法...
// 这些方法需要通过仓库路径来操作Pages
// 由于时间限制，这里只实现基本结构

func (api *PagesApi) ListRepoPages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "List repo pages"})
}

func (api *PagesApi) CreateRepoPages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Create repo pages"})
}

func (api *PagesApi) GetRepoPages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get repo pages"})
}

func (api *PagesApi) UpdateRepoPages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update repo pages"})
}

func (api *PagesApi) DeleteRepoPages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Delete repo pages"})
}

func (api *PagesApi) PublishRepoPages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Publish repo pages"})
}

func (api *PagesApi) ArchiveRepoPages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Archive repo pages"})
}

func (api *PagesApi) UnpublishRepoPages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Unpublish repo pages"})
}

func (api *PagesApi) GetRepoPagesConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get repo pages config"})
}

func (api *PagesApi) UpdateRepoPagesConfig(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Update repo pages config"})
}

func (api *PagesApi) GenerateRepoStaticSite(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Generate repo static site"})
}

func (api *PagesApi) GetRepoStaticSiteURL(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get repo static site URL"})
}
