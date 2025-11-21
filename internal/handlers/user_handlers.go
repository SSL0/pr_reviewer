package handlers

import (
	"errors"
	"net/http"
	"pr_reviewer/internal/dto"
	"pr_reviewer/internal/service"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterUsers(r *gin.RouterGroup) {
	r.POST("/users/setIsActive", h.setIsActive)
	r.GET("/users/getReview", h.getReview)
}

func (h *Handler) setIsActive(c *gin.Context) {
	var req dto.SetIsActiveRequest

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "invalid body data")
		return
	}

	user, err := h.Services.SetUserIsActive(req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, service.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, h.jsonError(ErrCodeNotFound, "resource not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, "internal server error")
		h.logger.Error("failed to create pull request", "error", err.Error())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) getReview(c *gin.Context) {
	userID := c.Query("user_id")

	res, err := h.Services.GetUserReviews(userID)
	if err != nil {
		if errors.Is(err, service.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, h.jsonError(ErrCodeNotFound, "resource not found"))
			return
		}
		c.JSON(http.StatusInternalServerError, "internal server error")
		h.logger.Error("failed to create pull request", "error", err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}
