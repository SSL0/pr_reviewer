package handler

import "github.com/gin-gonic/gin"

type StatusType string

const (
	StatusOpen  = "OPEN"
	StatusMerge = "MERGED"
)

func (h *Handler) RegisterPullRequsets(r *gin.RouterGroup) {
	r.POST("/pullRequest/create", h.create)
	r.POST("/pullRequest/merge", h.merge)
	r.POST("/pullRequest/reassign", h.reassign)
}

func (h *Handler) create(c *gin.Context) {

}

func (h *Handler) merge(c *gin.Context) {

}

func (h *Handler) reassign(c *gin.Context) {

}
