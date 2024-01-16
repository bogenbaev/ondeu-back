package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) initInfoRoutes(api *gin.RouterGroup) {
	info := api.Group("/")
	{
		info.GET("/roles", authorize(h.keycloak, []string{"admin", "manager", "student"}), h.getRoles)
	}
}

func (h *Handler) getRoles(ctx *gin.Context) {
	roles, err := h.services.InformationService.GetRoles(ctx)
	if err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, roles)
	return
}
