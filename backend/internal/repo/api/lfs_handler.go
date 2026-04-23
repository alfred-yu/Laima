package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"laima/internal/git"

	"github.com/gin-gonic/gin"
)

// LFSHandler 处理 Git LFS 相关请求
type LFSHandler struct {
	gitSvc *git.Service
}

// NewLFSHandler 创建 LFS 处理器
func NewLFSHandler(gitSvc *git.Service) *LFSHandler {
	return &LFSHandler{gitSvc: gitSvc}
}

// RegisterRoutes 注册 LFS 相关路由
func (h *LFSHandler) RegisterRoutes(router *gin.RouterGroup) {
	lfs := router.Group("/lfs")
	{
		lfs.POST("/objects", h.UploadObject)
		lfs.GET("/objects/:oid", h.DownloadObject)
	}
}

// LFSBatchRequest LFS 批量请求
type LFSBatchRequest struct {
	Operation string `json:"operation"`
	Objects   []struct {
		Oid  string `json:"oid"`
		Size int64  `json:"size"`
	} `json:"objects"`
	Transfers []string `json:"transfers,omitempty"`
	Ref       struct {
		Name string `json:"name"`
	} `json:"ref,omitempty"`
}

// LFSObjectAction LFS 对象动作
type LFSObjectAction struct {
	Href      string            `json:"href"`
	Header    map[string]string `json:"header,omitempty"`
	ExpiresIn int              `json:"expires_in,omitempty"`
}

// LFSObjectError LFS 对象错误
type LFSObjectError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// LFSResponseObject LFS 响应对象
type LFSResponseObject struct {
	Oid     string                     `json:"oid"`
	Size    int64                      `json:"size"`
	Actions map[string]LFSObjectAction `json:"actions,omitempty"`
	Error   *LFSObjectError            `json:"error,omitempty"`
}

// LFSBatchResponse LFS 批量响应
type LFSBatchResponse struct {
	Transfer string               `json:"transfer"`
	Objects  []LFSResponseObject `json:"objects"`
}

// Batch 处理 LFS 批量请求
func (h *LFSHandler) Batch(c *gin.Context) {
	var req LFSBatchRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	resp := LFSBatchResponse{
		Transfer: "basic",
		Objects:  make([]LFSResponseObject, len(req.Objects)),
	}

	for i, obj := range req.Objects {
		resp.Objects[i].Oid = obj.Oid
		resp.Objects[i].Size = obj.Size
		resp.Objects[i].Actions = make(map[string]LFSObjectAction)

		if req.Operation == "download" {
			resp.Objects[i].Actions["download"] = LFSObjectAction{
				Href:      "/api/v1/repos/lfs/objects/" + obj.Oid,
				ExpiresIn: 3600,
			}
		} else if req.Operation == "upload" {
			resp.Objects[i].Actions["upload"] = LFSObjectAction{
				Href:      "/api/v1/repos/lfs/objects",
				ExpiresIn: 3600,
			}
		}
	}

	c.JSON(http.StatusOK, resp)
}

// UploadObject 上传 LFS 对象
func (h *LFSHandler) UploadObject(c *gin.Context) {
	oid := c.Query("oid")
	sizeStr := c.Query("size")

	if oid == "" || sizeStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing oid or size"})
		return
	}

	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid size"})
		return
	}

	contentLength := c.Request.ContentLength
	if contentLength != size {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content length mismatch"})
		return
	}

	content := make([]byte, size)
	if _, err := c.Request.Body.Read(content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read content"})
		return
	}

	if err := h.gitSvc.UploadLFSObject(oid, size, content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

// DownloadObject 下载 LFS 对象
func (h *LFSHandler) DownloadObject(c *gin.Context) {
	oid := c.Param("oid")
	if oid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing oid"})
		return
	}

	content, err := h.gitSvc.DownloadLFSObject(oid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Object not found"})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", content)
}
