package handler

import "github.com/gin-gonic/gin"

type StatusType string

const (
	StatusOpen  = "OPEN"
	StatusMerge = "MERGED"
)

func RegisterPullRequsets(r *gin.RouterGroup) {
	r.POST("/pullRequest/create", create)
	r.POST("/pullRequest/merge", merge)
	r.POST("/pullRequest/reassign", reassign)
}

func create(c *gin.Context) {

}

func merge(c *gin.Context) {

}

func reassign(c *gin.Context) {

}
