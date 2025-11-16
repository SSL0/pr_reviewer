package handlers

import (
	"pr_reviewer/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	ErrCodeTeamExists  = "TEAM_EXISTS"
	ErrCodePRExists    = "PR_EXISTS"
	ErrCodePRMerged    = "PR_MERGED"
	ErrCodeNotAssigned = "NOT_ASSIGNED"
	ErrCodeNoCandidate = "NO_CANDIDATE"
	ErrCodeNotFound    = "NOT_FOUND"
)

type ErrorResponse struct {
	Code    string
	Message string
}

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) RegisterRoutes() *gin.Engine {
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group("/api/v1")
	{
		h.RegisterTeams(api)
		h.RegisterUsers(api)
		h.RegisterPullRequsets(api)
	}

	return router
}

func (h *Handler) jsonError(code string, message string) gin.H {
	return gin.H{
		"error": ErrorResponse{
			Code:    code,
			Message: message,
		},
	}
}
