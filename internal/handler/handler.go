package handler

import (
	"gitlab.com/a5805/ondeu/ondeu-back/internal/repository"
	keycloak2 "gitlab.com/a5805/ondeu/ondeu-back/pkg/gocloak"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	v1 "gitlab.com/a5805/ondeu/ondeu-back/internal/handler/v1"
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

func (h *Handler) Init() *gin.Engine {
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))

	router.Use(v1.CORSMiddleware())
	router.Use(v1.ReadRequestBody())

	router.Use(gin.Recovery())

	// third party handlers
	router.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "Up"}) })

	h.InitRoutes(router)

	return router
}

func (h *Handler) InitRoutes(router *gin.Engine) {
	handler := v1.NewHandler(h.services, h.repos, h.keycloak)

	api := router.Group("/api")
	{
		handler.Init(api)
	}
}
