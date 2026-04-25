package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"laima/internal/cicd/app"
	"laima/internal/cicd/domain"
)

type RunnerHandler struct {
	service app.RunnerService
}

func NewRunnerHandler(service app.RunnerService) *RunnerHandler {
	return &RunnerHandler{
		service: service,
	}
}

func (h *RunnerHandler) RegisterRoutes(r *gin.RouterGroup) {
	runners := r.Group("/runners")
	{
		runners.POST("/register", h.RegisterRunner)
		runners.POST("/heartbeat", h.Heartbeat)
		runners.POST("/job/request", h.RequestJob)
		runners.POST("/job/update", h.UpdateJobStatus)
		runners.GET("", h.ListRunners)
		runners.GET("/:id", h.GetRunner)
		runners.DELETE("/:id", h.DeleteRunner)
	}
}

func (h *RunnerHandler) RegisterRunner(c *gin.Context) {
	var req domain.RunnerRegistration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.service.RegisterRunner(&req)
	if err != nil {
		if err == app.ErrRunnerAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *RunnerHandler) Heartbeat(c *gin.Context) {
	var req domain.RunnerHeartbeat
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Heartbeat(&req); err != nil {
		if err == app.ErrRunnerNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == app.ErrInvalidToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "heartbeat received"})
}

func (h *RunnerHandler) RequestJob(c *gin.Context) {
	var req domain.RunnerJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.service.RequestJob(&req)
	if err != nil {
		if err == app.ErrRunnerNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == app.ErrInvalidToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if err == app.ErrRunnerBusy {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if job == nil {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.JSON(http.StatusOK, job)
}

func (h *RunnerHandler) UpdateJobStatus(c *gin.Context) {
	var req domain.RunnerJobUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateJobStatus(&req); err != nil {
		if err == app.ErrRunnerNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == app.ErrInvalidToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "job status updated"})
}

func (h *RunnerHandler) ListRunners(c *gin.Context) {
	runners, err := h.service.ListRunners()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, runners)
}

func (h *RunnerHandler) GetRunner(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid runner id"})
		return
	}

	runner, err := h.service.GetRunner(id)
	if err != nil {
		if err == app.ErrRunnerNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, runner)
}

func (h *RunnerHandler) DeleteRunner(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid runner id"})
		return
	}

	if err := h.service.DeleteRunner(id); err != nil {
		if err == app.ErrRunnerNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
