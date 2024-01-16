package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules/dto"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/utils"
	"net/http"
)

func (h *Handler) initTreeRoutes(api *gin.RouterGroup) {
	crud := api.Group("/")
	{
		crud.POST("/", authorize(h.keycloak, []string{"admin", "manager", "student"}), h.createTree)
		crud.GET("/:treeID", authorize(h.keycloak, []string{"admin", "manager", "student"}), h.getTree)
		crud.GET("/:treeID/list", authorize(h.keycloak, []string{"admin", "manager", "student"}), h.listTree)
		crud.PUT("/:treeID", authorize(h.keycloak, []string{"admin", "manager", "student"}), h.updateTree)
		crud.DELETE("/:treeID", authorize(h.keycloak, []string{"admin", "manager", "student"}), h.deleteTree)
	}
}

type TreeInput struct {
	TreeID uint `uri:"treeID" binding:"required"`
}

func (h *Handler) createTree(ctx *gin.Context) {
	var tree dto.Tree
	if err := ctx.ShouldBind(&tree); err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	tree, err := h.services.TreeService.Create(ctx, tree)
	if err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tree)
	return
}

func (h *Handler) getTree(ctx *gin.Context) {
	id, err := utils.ParseUint(ctx.Param("treeID"))
	if err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	if id == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": "invalid tree id"})
		return
	}

	tree, err := h.services.TreeService.Get(ctx, dto.Tree{ID: id})
	if err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	docs, err := h.services.DocumentService.ListByTree(ctx, []uint{id})
	if err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	completeTree := h.services.TreeService.FormTree(ctx, []dto.Tree{tree}, docs)

	ctx.JSON(http.StatusOK, completeTree)
	return
}

func (h *Handler) listTree(ctx *gin.Context) {
	id, err := utils.ParseUint(ctx.Param("treeID"))
	if err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	tree := dto.Tree{ID: id}

	trees, err := h.services.TreeService.List(ctx, tree)
	if err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	ids := h.services.TreeService.GetTreeIDs(ctx, trees)

	docs, err := h.services.DocumentService.ListByTree(ctx, ids)
	if err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	completeTree := h.services.TreeService.FormTree(ctx, trees, docs)

	ctx.JSON(http.StatusOK, completeTree)
	return
}

func (h *Handler) updateTree(ctx *gin.Context) {
	var input TreeInput
	if err := ctx.ShouldBindUri(&input); err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	var tree dto.Tree
	if err := ctx.ShouldBind(&tree); err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	tree.ID = input.TreeID

	updated, err := h.services.TreeService.Update(ctx, tree)
	if err != nil {
		logrus.Errorf("[service error] - %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updated)
	return
}

func (h *Handler) deleteTree(ctx *gin.Context) {
	var input TreeInput
	if err := ctx.ShouldBindUri(&input); err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	tree := dto.Tree{ID: input.TreeID}

	deleted, err := h.services.TreeService.Delete(ctx, tree)
	if err != nil {
		logrus.Errorf("[service error] - %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, deleted)
	return
}
