package handler

import "github.com/gin-gonic/gin"

func RegisterTeams(r *gin.RouterGroup) {
	r.POST("/team/add", add)
	r.GET("/team/get", get)
}

func add(c *gin.Context) {

}

func get(c *gin.Context) {

}
