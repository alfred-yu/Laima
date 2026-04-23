package api

import (
	"net/http"
	"strconv"

	"laima/internal/git"
	"laima/internal/middleware"
	"laima/internal/repo/app"
	"laima/internal/repo/domain"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
	"github.com/minio/minio-go/v7"
	"github.com/meilisearch/meilisearch-go"
)

// RepoAPI 仓库 API 结构体
type RepoAPI struct {
	repoService app.RepoService
	db         *gorm.DB
	redis      *redis.Client
	minio      *minio.Client
	meili      meilisearch.ServiceManager
	gitSvc     *git.Service
}

// NewRepoAPI 创建仓库 API 实例
func NewRepoAPI(db *gorm.DB, redis *redis.Client, minio *minio.Client, meili meilisearch.ServiceManager, gitSvc *git.Service) *RepoAPI {
	return &RepoAPI{
		repoService: app.NewRepoService(db, gitSvc, meili),
		db:    db,
		redis: redis,
		minio: minio,
		meili: meili,
		gitSvc: gitSvc,
	}
}

// getCurrentUserID 获取当前用户 ID
func getCurrentUserID(c *gin.Context) int64 {
	userID, _ := c.Get("user_id")
	return int64(userID.(int))
}

// RegisterRoutes 注册路由
func (api *RepoAPI) RegisterRoutes(r *gin.Engine) {
	repoGroup := r.Group("/api/v1/repos")
	{
		// 仓库 CRUD - 需要认证的路由
		repoGroup.GET("", api.ListRepos)
		repoGroup.POST("", middleware.AuthMiddleware(), api.CreateRepo)
		repoGroup.GET("/:owner/:repo", api.GetRepo)
		repoGroup.PUT("/:owner/:repo", middleware.AuthMiddleware(), api.UpdateRepo)
		repoGroup.DELETE("/:owner/:repo", middleware.AuthMiddleware(), api.DeleteRepo)

		// Fork & 导入 - 需要认证
		repoGroup.POST("/:owner/:repo/forks", middleware.AuthMiddleware(), api.ForkRepo)
		repoGroup.POST("/import", middleware.AuthMiddleware(), api.ImportRepo)

		// 分支操作 - 需要认证
		repoGroup.GET("/:owner/:repo/branches", api.ListBranches)
		repoGroup.POST("/:owner/:repo/branches", middleware.AuthMiddleware(), api.CreateBranch)
		repoGroup.DELETE("/:owner/:repo/branches/:branch", middleware.AuthMiddleware(), api.DeleteBranch)

		// 标签操作
		repoGroup.GET("/:owner/:repo/tags", api.ListTags)
		repoGroup.POST("/:owner/:repo/tags", middleware.AuthMiddleware(), api.CreateTag)
		repoGroup.DELETE("/:owner/:repo/tags/:tag", middleware.AuthMiddleware(), api.DeleteTag)

		// 代码浏览
		repoGroup.GET("/:owner/:repo/contents/*path", api.GetContent)
		repoGroup.GET("/:owner/:repo/blob/:ref/*path", api.GetBlob)
		repoGroup.GET("/:owner/:repo/blame/:ref/*path", api.GetBlame)

		// 代码搜索
		repoGroup.GET("/:owner/:repo/search", api.SearchCode)

		// 统计 - 需要认证
		repoGroup.POST("/:owner/:repo/star", middleware.AuthMiddleware(), api.StarRepo)
		repoGroup.DELETE("/:owner/:repo/star", middleware.AuthMiddleware(), api.UnstarRepo)
	}

	// LFS 路由
	lfsHandler := NewLFSHandler(api.gitSvc)
	lfsHandler.RegisterRoutes(r.Group("/api/v1/repos"))
}

// ListRepos 列出仓库
func (api *RepoAPI) ListRepos(c *gin.Context) {
	filter := &app.RepoFilter{
		Page:    1,
		PerPage: 30,
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		filter.Page = page
	}

	if perPage, err := strconv.Atoi(c.Query("per_page")); err == nil && perPage > 0 {
		filter.PerPage = perPage
	}

	filter.Search = c.Query("q")
	filter.Visibility = c.Query("visibility")

	repos, total, err := api.repoService.ListRepos(c, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"items": repos,
		"page":  filter.Page,
		"per_page": filter.PerPage,
	})
}

// CreateRepo 创建仓库
func (api *RepoAPI) CreateRepo(c *gin.Context) {
	var req app.CreateRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果没有指定所有者，默认使用当前用户
	if req.OwnerID == 0 {
		req.OwnerID = getCurrentUserID(c)
		req.OwnerType = domain.OwnerTypeUser
	}

	repo, err := api.repoService.CreateRepo(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, repo)
}

// GetRepo 获取仓库详情
func (api *RepoAPI) GetRepo(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	c.JSON(http.StatusOK, repo)
}

// UpdateRepo 更新仓库
func (api *RepoAPI) UpdateRepo(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	var req app.UpdateRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedRepo, err := api.repoService.UpdateRepo(c, repo.ID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedRepo)
}

// DeleteRepo 删除仓库
func (api *RepoAPI) DeleteRepo(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	if err := api.repoService.DeleteRepo(c, repo.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repository deleted successfully"})
}

// ForkRepo Fork 仓库
func (api *RepoAPI) ForkRepo(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	var req struct {
		Namespace string `json:"namespace" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	forkedRepo, err := api.repoService.ForkRepo(c, repo.ID, req.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, forkedRepo)
}

// ImportRepo 导入仓库
func (api *RepoAPI) ImportRepo(c *gin.Context) {
	var req app.ImportRepoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := api.repoService.ImportRepo(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, task)
}

// ListBranches 列出分支
func (api *RepoAPI) ListBranches(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	branches, err := api.repoService.ListBranches(c, repo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, branches)
}

// CreateBranch 创建分支
func (api *RepoAPI) CreateBranch(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	var req app.CreateBranchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	branch, err := api.repoService.CreateBranch(c, repo.ID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, branch)
}

// DeleteBranch 删除分支
func (api *RepoAPI) DeleteBranch(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	branch := c.Param("branch")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	if err := api.repoService.DeleteBranch(c, repo.ID, branch); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Branch deleted successfully"})
}

// ListTags 列出标签
func (api *RepoAPI) ListTags(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	tags, err := api.repoService.ListTags(c, repo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}

// CreateTag 创建标签
func (api *RepoAPI) CreateTag(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	var req app.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag, err := api.repoService.CreateTag(c, repo.ID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tag)
}

// DeleteTag 删除标签
func (api *RepoAPI) DeleteTag(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	tag := c.Param("tag")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	if err := api.repoService.DeleteTag(c, repo.ID, tag); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag deleted successfully"})
}

// GetContent 获取文件内容
func (api *RepoAPI) GetContent(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	path := c.Param("path")
	ref := c.DefaultQuery("ref", "main")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	tree, err := api.repoService.GetTree(c, repo.ID, ref, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tree)
}

// GetBlob 获取文件内容（原始）
func (api *RepoAPI) GetBlob(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	ref := c.Param("ref")
	path := c.Param("path")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	blob, err := api.repoService.GetBlob(c, repo.ID, ref, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blob)
}

// GetBlame 获取 Blame 信息
func (api *RepoAPI) GetBlame(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	ref := c.Param("ref")
	path := c.Param("path")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	blame, err := api.repoService.GetBlame(c, repo.ID, ref, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, blame)
}

// SearchCode 搜索代码
func (api *RepoAPI) SearchCode(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	query := &app.SearchQuery{
		Query:   c.Query("q"),
		RepoID:  repo.ID,
		Page:    1,
		PerPage: 30,
	}

	if page, err := strconv.Atoi(c.Query("page")); err == nil && page > 0 {
		query.Page = page
	}

	if perPage, err := strconv.Atoi(c.Query("per_page")); err == nil && perPage > 0 {
		query.PerPage = perPage
	}

	results, total, err := api.repoService.SearchCode(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    total,
		"items":    results,
		"page":     query.Page,
		"per_page": query.PerPage,
	})
}

// StarRepo 收藏仓库
func (api *RepoAPI) StarRepo(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	if err := api.repoService.StarRepo(c, repo.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repository starred successfully"})
}

// UnstarRepo 取消收藏仓库
func (api *RepoAPI) UnstarRepo(c *gin.Context) {
	owner := c.Param("owner")
	repoName := c.Param("repo")
	fullPath := owner + "/" + repoName

	repo, err := api.repoService.GetRepoByPath(c, fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	if err := api.repoService.UnstarRepo(c, repo.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repository unstarred successfully"})
}