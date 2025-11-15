package handler

import "github.com/gin-gonic/gin"

func (h *Handler) RegisterTeams(r *gin.RouterGroup) {
	r.POST("/team/add", h.add)
	r.GET("/team/get", h.get)
}

func (h *Handler) add(c *gin.Context) {

}

func (h *Handler) get(c *gin.Context) {

}
