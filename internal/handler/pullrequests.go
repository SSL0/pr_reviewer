package handler

import (
	"errors"
	"net/http"
	"pr_reviewer/internal/dto"
	"pr_reviewer/internal/service"

	"github.com/gin-gonic/gin"
)

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
	var req dto.CreatePullRequestRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "invalid body data")
		return
	}

	pr, err := h.services.Create(req.PullRequestID, req.PullRequestName, req.AuthorID)

	if err != nil {
		if errors.Is(err, service.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, h.jsonError(ErrorCodeNotFound, "resource not found"))
			return
		}

		if errors.Is(err, service.ErrPRExists) {
			c.JSON(http.StatusNotFound, h.jsonError(ErrorCodePRExists, "PR id already exists"))
			return
		}

		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, pr)

}

func (h *Handler) merge(c *gin.Context) {

}

func (h *Handler) reassign(c *gin.Context) {

}
