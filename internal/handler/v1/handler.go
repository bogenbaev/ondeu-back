package v1

import (
	"gitlab.com/a5805/ondeu/ondeu-back/internal/repository"
	keycloak2 "gitlab.com/a5805/ondeu/ondeu-back/pkg/gocloak"

	"github.com/gin-gonic/gin"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/service"
)

type Handler struct {
	services *service.Services
	keycloak keycloak2.IKeycloak
	repos    *repository.Repository
}

func NewHandler(services *service.Services, repos *repository.Repository, keycloak keycloak2.IKeycloak) *Handler {
	return &Handler{
		services: services,
		keycloak: keycloak,
		repos:    repos,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		tree := v1.Group("/tree")
		{
			h.initDocumentsRoutes(tree)
			h.initTreeRoutes(tree)
		}
		info := v1.Group("/info")
		{
			h.initInfoRoutes(info)
		}
	}
}
