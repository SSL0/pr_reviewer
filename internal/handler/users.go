package handler

import "github.com/gin-gonic/gin"

func RegisterUsers(r *gin.RouterGroup) {
	r.POST("/users/setIsActive", setIsActive)
	r.GET("/users/getReview", getReview)
}

func setIsActive(c *gin.Context) {

}

func getReview(c *gin.Context) {

}
