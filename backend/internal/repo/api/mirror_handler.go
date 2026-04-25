package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"laima/internal/repo/app"
	"laima/internal/repo/domain"
)

type MirrorHandler struct {
	service app.MirrorService
}

func NewMirrorHandler(service app.MirrorService) *MirrorHandler {
	return &MirrorHandler{
		service: service,
	}
}

func (h *MirrorHandler) RegisterRoutes(r *gin.RouterGroup) {
	mirrors := r.Group("/repos/:owner/:repo/mirror")
	{
		mirrors.POST("", h.CreateMirror)
		mirrors.GET("", h.GetMirror)
		mirrors.PUT("", h.UpdateMirror)
		mirrors.DELETE("", h.DeleteMirror)
		mirrors.POST("/sync", h.SyncMirror)
		mirrors.GET("/status", h.GetMirrorStatus)
		mirrors.GET("/logs", h.GetSyncLogs)
	}
}

func (h *MirrorHandler) CreateMirror(c *gin.Context) {
	var req domain.MirrorCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.CreateMirror(c.Request.Context(), &req)
	if err != nil {
		if err == app.ErrMirrorAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *MirrorHandler) GetMirror(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	repoID, err := h.getRepoID(owner, repo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repository"})
		return
	}

	resp, err := h.service.GetMirrorByRepo(c.Request.Context(), repoID)
	if err != nil {
		if err == app.ErrMirrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MirrorHandler) UpdateMirror(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	repoID, err := h.getRepoID(owner, repo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repository"})
		return
	}

	mirror, err := h.service.GetMirrorByRepo(c.Request.Context(), repoID)
	if err != nil {
		if err == app.ErrMirrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var req domain.MirrorUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.UpdateMirror(c.Request.Context(), mirror.ID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MirrorHandler) DeleteMirror(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	repoID, err := h.getRepoID(owner, repo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repository"})
		return
	}

	mirror, err := h.service.GetMirrorByRepo(c.Request.Context(), repoID)
	if err != nil {
		if err == app.ErrMirrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.DeleteMirror(c.Request.Context(), mirror.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *MirrorHandler) SyncMirror(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	repoID, err := h.getRepoID(owner, repo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repository"})
		return
	}

	mirror, err := h.service.GetMirrorByRepo(c.Request.Context(), repoID)
	if err != nil {
		if err == app.ErrMirrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	syncResp, err := h.service.SyncMirror(c.Request.Context(), mirror.ID)
	if err != nil {
		if err == app.ErrMirrorDisabled {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, syncResp)
}

func (h *MirrorHandler) GetMirrorStatus(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	repoID, err := h.getRepoID(owner, repo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repository"})
		return
	}

	resp, err := h.service.GetMirrorByRepo(c.Request.Context(), repoID)
	if err != nil {
		if err == app.ErrMirrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *MirrorHandler) GetSyncLogs(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	repoID, err := h.getRepoID(owner, repo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid repository"})
		return
	}

	mirror, err := h.service.GetMirrorByRepo(c.Request.Context(), repoID)
	if err != nil {
		if err == app.ErrMirrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	logs, err := h.service.GetSyncLogs(c.Request.Context(), mirror.ID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func (h *MirrorHandler) getRepoID(owner, repo string) (int64, error) {
	return int64(1), nil
}
