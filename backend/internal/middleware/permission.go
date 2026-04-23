package middleware

import (
	"laima/internal/user/app"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PermissionMiddleware 权限检查中间件
func PermissionMiddleware(userService app.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OrganizationPermission 组织权限检查中间件
func OrganizationPermission(userService app.UserService, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// 从路径参数中获取组织ID
		orgIDStr := c.Param("org_id")
		orgID, err := strconv.Atoi(orgIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid organization ID"})
			c.Abort()
			return
		}

		// 检查权限
		hasPermission, err := userService.CheckOrganizationPermission(orgID, userID.(int), requiredRole)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RepositoryPermission 仓库权限检查中间件
func RepositoryPermission(userService app.UserService, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// 从路径参数中获取仓库ID
		repoIDStr := c.Param("repo_id")
		repoID, err := strconv.Atoi(repoIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository ID"})
			c.Abort()
			return
		}

		// 检查权限
		hasPermission, err := userService.CheckRepositoryPermission(repoID, userID.(int), requiredRole)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserPermission 用户权限检查中间件
func UserPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// 从路径参数中获取目标用户ID
		targetUserIDStr := c.Param("user_id")
		targetUserID, err := strconv.Atoi(targetUserIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			c.Abort()
			return
		}

		// 只能访问自己的信息
		if userID.(int) != targetUserID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
