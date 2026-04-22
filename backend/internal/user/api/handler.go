package api

import (
	"laima/internal/middleware"
	"laima/internal/user/app"
	"laima/internal/user/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// UserAPI 用户API处理器
type UserAPI struct {
	db          *gorm.DB
	redisClient *redis.Client
	userService app.UserService
}

// NewUserAPI 创建用户API实例
func NewUserAPI(db *gorm.DB, redisClient *redis.Client) *UserAPI {
	return &UserAPI{
		db:          db,
		redisClient: redisClient,
		userService: app.NewUserService(db),
	}
}

// RegisterRoutes 注册路由
func (api *UserAPI) RegisterRoutes(r *gin.Engine) {
	// 认证路由
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", api.Login)
		auth.POST("/register", api.Register)
	}

	// 用户路由
	user := r.Group("/api/users")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/me", api.GetCurrentUser)
		user.PUT("/me", api.UpdateCurrentUser)
		user.GET("/:id", api.GetUserByID)
		user.GET("/username/:username", api.GetUserByUsername)
	}

	// 组织路由
	org := r.Group("/api/orgs")
	org.Use(middleware.AuthMiddleware())
	{
		org.POST("/", api.CreateOrganization)
		org.GET("/:id", api.GetOrganization)
		org.PUT("/:id", api.UpdateOrganization)
		org.DELETE("/:id", api.DeleteOrganization)
		org.GET("/:id/members", api.GetOrganizationMembers)
		org.POST("/:id/members", api.AddOrganizationMember)
		org.DELETE("/:id/members/:user_id", api.RemoveOrganizationMember)
		org.PUT("/:id/members/:user_id/role", api.UpdateOrganizationMemberRole)
	}

	// 仓库成员路由
	repoMembers := r.Group("/api/repos/:repo_id/members")
	repoMembers.Use(middleware.AuthMiddleware())
	{
		repoMembers.GET("/", api.GetRepositoryMembers)
		repoMembers.POST("/", api.AddRepositoryMember)
		repoMembers.DELETE("/:user_id", api.RemoveRepositoryMember)
		repoMembers.PUT("/:user_id/role", api.UpdateRepositoryMemberRole)
	}
}

// getCurrentUserID 从上下文中获取当前用户ID
func getCurrentUserID(c *gin.Context) int {
	userID, _ := c.Get("user_id")
	return userID.(int)
}

// Login 用户登录
func (api *UserAPI) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := api.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Register 用户注册
func (api *UserAPI) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := api.userService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetCurrentUser 获取当前用户
func (api *UserAPI) GetCurrentUser(c *gin.Context) {
	userID := getCurrentUserID(c)

	user, err := api.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateCurrentUser 更新当前用户
func (api *UserAPI) UpdateCurrentUser(c *gin.Context) {
	userID := getCurrentUserID(c)

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 禁止修改密码和ID
	delete(updates, "password")
	delete(updates, "id")

	user, err := api.userService.UpdateUser(userID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserByID 根据ID获取用户
func (api *UserAPI) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := api.userService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUserByUsername 根据用户名获取用户
func (api *UserAPI) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := api.userService.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// CreateOrganization 创建组织
func (api *UserAPI) CreateOrganization(c *gin.Context) {
	userID := getCurrentUserID(c)

	var req struct {
		Name        string `json:"name" binding:"required"`
		DisplayName string `json:"display_name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := api.userService.CreateOrganization(userID, req.Name, req.DisplayName, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, org)
}

// GetOrganization 获取组织
func (api *UserAPI) GetOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	org, err := api.userService.GetOrganizationByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		return
	}

	c.JSON(http.StatusOK, org)
}

// UpdateOrganization 更新组织
func (api *UserAPI) UpdateOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org, err := api.userService.UpdateOrganization(id, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, org)
}

// DeleteOrganization 删除组织
func (api *UserAPI) DeleteOrganization(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	if err := api.userService.DeleteOrganization(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}

// GetOrganizationMembers 获取组织成员
func (api *UserAPI) GetOrganizationMembers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	members, err := api.userService.GetOrganizationMembers(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

// AddOrganizationMember 添加组织成员
func (api *UserAPI) AddOrganizationMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	var req struct {
		UserID int    `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := api.userService.AddOrganizationMember(id, req.UserID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// RemoveOrganizationMember 移除组织成员
func (api *UserAPI) RemoveOrganizationMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := api.userService.RemoveOrganizationMember(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}

// UpdateOrganizationMemberRole 更新组织成员角色
func (api *UserAPI) UpdateOrganizationMemberRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := api.userService.UpdateOrganizationMemberRole(id, userID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, member)
}

// GetRepositoryMembers 获取仓库成员
func (api *UserAPI) GetRepositoryMembers(c *gin.Context) {
	repoIDStr := c.Param("repo_id")
	repoID, err := strconv.Atoi(repoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	members, err := api.userService.GetRepositoryMembers(repoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

// AddRepositoryMember 添加仓库成员
func (api *UserAPI) AddRepositoryMember(c *gin.Context) {
	repoIDStr := c.Param("repo_id")
	repoID, err := strconv.Atoi(repoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	var req struct {
		UserID int    `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := api.userService.AddRepositoryMember(repoID, req.UserID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// RemoveRepositoryMember 移除仓库成员
func (api *UserAPI) RemoveRepositoryMember(c *gin.Context) {
	repoIDStr := c.Param("repo_id")
	repoID, err := strconv.Atoi(repoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := api.userService.RemoveRepositoryMember(repoID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}

// UpdateRepositoryMemberRole 更新仓库成员角色
func (api *UserAPI) UpdateRepositoryMemberRole(c *gin.Context) {
	repoIDStr := c.Param("repo_id")
	repoID, err := strconv.Atoi(repoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
		return
	}

	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member, err := api.userService.UpdateRepositoryMemberRole(repoID, userID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, member)
}
