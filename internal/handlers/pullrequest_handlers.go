package handlers

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

	res, err := h.services.CreatePullRequest(req.PullRequestID, req.PullRequestName, req.AuthorID)

	if err != nil {
		if errors.Is(err, service.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, h.jsonError(ErrCodeNotFound, "resource not found"))
			return
		}

		if errors.Is(err, service.ErrPRExists) {
			c.JSON(http.StatusNotFound, h.jsonError(ErrCodePRExists, "PR id already exists"))
			return
		}

		c.JSON(http.StatusInternalServerError, "internal server error")
		h.logger.Error("failed to create pull request", "error", err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) merge(c *gin.Context) {
	var req dto.MergePullRequestRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "invalid body data")
		return
	}

	res, err := h.services.MergePullRequest(req.PullRequestID)

	if err != nil {
		if errors.Is(err, service.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, h.jsonError(ErrCodeNotFound, "resource not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, "internal server error")
		h.logger.Error("failed to merge pull request", "error", err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) reassign(c *gin.Context) {
	var req dto.ReassignPullRequestRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "invalid body data")
		return
	}
	res, err := h.services.ReassignPullRequestReviewer(req.PullRequestID, req.OldReviewerID)
	if err != nil {
		if errors.Is(err, service.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, h.jsonError(ErrCodeNotFound, "resource not found"))
			return
		}

		if errors.Is(err, service.ErrPRMerged) {
			c.JSON(http.StatusConflict, h.jsonError(ErrCodePRMerged, "cannot reassign on merged PR"))
			return
		}

		if errors.Is(err, service.ErrNotAssigned) {
			c.JSON(http.StatusConflict, h.jsonError(ErrCodeNotAssigned, "reviewer is not assigned to this PR"))
			return
		}

		if errors.Is(err, service.ErrNoCanditate) {
			c.JSON(http.StatusConflict, h.jsonError(ErrCodeNoCandidate, "no active replacement candidate in team"))
			return
		}

		c.JSON(http.StatusInternalServerError, "internal server error")
		h.logger.Error("failed to reassign pull request reviewer", "error", err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}
