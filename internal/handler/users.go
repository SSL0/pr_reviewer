package handler

import (
	"errors"
	"log"
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

	user, err := h.services.SetIsActive(req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, service.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, h.jsonError(ErrorCodeNotFound, "resource not found"))
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, user)

}

func (h *Handler) getReview(c *gin.Context) {
	userID := c.Query("user_id")

	res, err := h.services.GetReview(userID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}
	c.JSON(http.StatusOK, res)
}
