package handler

import (
	"errors"
	"log"
	"net/http"
	"pr_reviewer/internal/domain"
	"pr_reviewer/internal/dto"
	"pr_reviewer/internal/service"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterTeams(r *gin.RouterGroup) {
	r.POST("/team/add", h.add)
	r.GET("/team/get", h.get)
}

func (h *Handler) add(c *gin.Context) {
	var req dto.Team

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "invalid body data")
		return
	}

	if req.TeamName == "" {
		c.JSON(http.StatusBadRequest, "invalid body data")
		return
	}

	var domainTeamMembers []domain.TeamMember

	for _, m := range req.Members {
		domainTeamMembers = append(domainTeamMembers,
			domain.TeamMember{
				UserID:   m.UserID,
				Username: m.Username,
				IsActive: m.IsActive,
			})

	}

	team, err := h.services.AddTeam(
		domain.Team{TeamName: req.TeamName, Members: domainTeamMembers},
	)

	if err != nil {
		log.Println(err)
		if errors.Is(err, service.ErrTeamExists) {
			c.JSON(http.StatusBadRequest, h.jsonError(ErrCodeTeamExists, "team_name already exists"))
			return
		}
		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusCreated, team)
}

func (h *Handler) get(c *gin.Context) {
	teamName := c.Query("team_name")

	team, err := h.services.GetTeam(teamName)

	if err != nil {
		log.Println(err)

		if errors.Is(err, service.ErrResourceNotFound) {
			c.JSON(http.StatusNotFound, h.jsonError(ErrCodeNotFound, "resource not found"))
			return
		}

		c.JSON(http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(http.StatusOK, team)
}
