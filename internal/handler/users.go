package handler

import (
	"errors"
	"log"
	"net/http"
	"pr_service/internal/service"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterUsers(r *gin.RouterGroup) {
	r.POST("/users/setIsActive", h.setIsActive)
	r.GET("/users/getReview", h.getReview)
}

type SetIsActiveBody struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (h *Handler) setIsActive(c *gin.Context) {
	var body SetIsActiveBody

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, "invalid body data")
		return
	}

	user, err := h.services.SetIsActive(body.UserID, body.IsActive)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {

			c.JSON(http.StatusBadRequest, gin.H{
				"error": ErrorResponse{
					Code:    ErrorCodeNotFound,
					Message: "resource not found",
				},
			})

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
