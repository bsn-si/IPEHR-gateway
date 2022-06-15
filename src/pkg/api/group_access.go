package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/group_access"
	"io"
	"io/ioutil"
	"net/http"
)

type GroupAccessHandler struct {
	*group_access.GroupAccessService
}

func NewGroupAccessHandler(docService *service.DefaultDocumentService, cfg *config.Config) *GroupAccessHandler {
	return &GroupAccessHandler{
		group_access.NewGroupAccessService(docService, cfg),
	}
}

func (h *GroupAccessHandler) Create(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(c.Request.Body)

	var request model.GroupAccessCreateRequest

	if err = json.Unmarshal(data, &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request decoding error"})
		return
	}

	if !request.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	newGroupAccess, err := h.GroupAccessService.Create(userId, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Group Access creating error"})
		return
	}

	h.respondWithDoc(newGroupAccess, c)
}

func (h *GroupAccessHandler) Get(c *gin.Context) {
	groupId := c.Param("group_id")

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	accessGroup, err := h.GroupAccessService.Get(userId, groupId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group access not found"})
	}

	c.JSON(http.StatusOK, accessGroup)
}

func (h *GroupAccessHandler) respondWithDoc(doc *model.GroupAccess, c *gin.Context) {
	c.Header("Location", h.Cfg.BaseUrl+"/v1/access/group/"+doc.GroupId)
	c.Header("ETag", doc.GroupId)

	c.JSON(http.StatusCreated, doc)
}
