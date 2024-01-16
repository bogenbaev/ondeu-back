package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gitlab.com/a5805/ondeu/ondeu-back/internal/repository"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/modules/dto"
	"gitlab.com/a5805/ondeu/ondeu-back/pkg/utils"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) initDocumentsRoutes(api *gin.RouterGroup) {
	crud := api.Group("/:treeID/document")
	{
		crud.POST("/", authorize(h.keycloak, []string{"admin", "manager", "student"}), h.createDocument)
		crud.GET("/:docID", authorize(h.keycloak, []string{"admin", "manager", "student"}), h.readDocument)
		crud.GET("/filter", authorize(h.keycloak, []string{"admin", "manager", "student"}), h.filterDocument)
		crud.GET("/:docID/share", authorize(h.keycloak, []string{"admin", "manager", "student"}), h.shareDocument)
		crud.PUT("/:docID", authorize(h.keycloak, []string{"admin", "manager", "student"}), h.updateDocument)
		crud.DELETE("/:docID", authorize(h.keycloak, []string{"admin", "manager", "student"}), h.deleteDocument)
	}
}

type DocumentInput struct {
	DocumentID uint `uri:"docID" binding:"required"`
	TreeID     uint `uri:"treeID" binding:"required"`
}

type DocumentRepos struct {
	repos repository.DocumentRepository
}

func (h *Handler) createDocument(ctx *gin.Context) {
	var document dto.Document
	if err := ctx.ShouldBind(&document); err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	treeID, err := utils.ParseUint(ctx.Param("treeID"))
	if err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	document.TreeID = treeID

	newDoc, err := h.services.DocumentService.Create(ctx, document, file)
	if err != nil {
		logrus.Errorf("[service error] - %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, newDoc)
	return
}

func (h *Handler) readDocument(ctx *gin.Context) {
	var input DocumentInput
	if err := ctx.ShouldBindUri(&input); err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	_, download := ctx.GetQuery("download")

	document := dto.Document{ID: input.DocumentID, TreeID: input.TreeID}

	stored, err := h.services.DocumentService.Get(ctx, document, download)
	if err != nil {
		logrus.Errorf("[service error] - %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	if download {
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", stored.Name+stored.Extension))
		ctx.Data(http.StatusOK, "application/octet-stream", stored.ResponseContent)
		return
	}

	ctx.JSON(http.StatusOK, stored)
	return
}

func (h *Handler) updateDocument(ctx *gin.Context) {
	var input DocumentInput
	if err := ctx.ShouldBindUri(&input); err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	var doc dto.Document
	if err := ctx.ShouldBind(&doc); err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	doc.ID = input.DocumentID
	doc.TreeID = input.TreeID

	updated, err := h.services.DocumentService.Update(ctx, doc)
	if err != nil {
		logrus.Errorf("[service error] - %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updated)
	return
}

func (h *Handler) deleteDocument(ctx *gin.Context) {
	var input DocumentInput
	if err := ctx.ShouldBindUri(&input); err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	document := dto.Document{ID: input.DocumentID, TreeID: input.TreeID}

	deleted, err := h.services.DocumentService.Delete(ctx, document)
	if err != nil {
		logrus.Errorf("[service error] - %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, deleted)
	return
}

func (h *Handler) shareDocument(ctx *gin.Context) {
	var input DocumentInput
	if err := ctx.ShouldBindUri(&input); err != nil {
		logrus.Errorf("[validaton error] - %+v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
		return
	}

	expire := 3600 * time.Second
	rawExpire, ok := ctx.GetQuery("expire")
	if ok {
		if try, err := strconv.Atoi(rawExpire); err == nil {
			expire = time.Duration(try) * time.Second
		}
	}

	document := dto.Document{ID: input.DocumentID, TreeID: input.TreeID}

	stored, err := h.services.DocumentService.Share(ctx, document, expire)
	if err != nil {
		logrus.Errorf("[service error] - %+v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stored)
	return
}

func (h *Handler) filterDocument(ctx *gin.Context) {
	field := ctx.Query("field")
	param := ctx.Query("param")

	tx, err := h.repos.DocumentRepository.FindByCondition(ctx, field, param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err})
		return
	}

	ctx.JSON(http.StatusOK, tx)
}
