package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/groupAccess"
)

type GroupAccessHandler struct {
	service *groupAccess.Service
	baseURL string
}

func NewGroupAccessHandler(docService *service.DefaultDocumentService, groupAccessService *groupAccess.Service, baseURL string) *GroupAccessHandler {
	return &GroupAccessHandler{
		service: groupAccessService,
		baseURL: baseURL,
	}
}

func (h *GroupAccessHandler) Create(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}
	defer c.Request.Body.Close()

	var request model.GroupAccessCreateRequest

	if err = json.Unmarshal(data, &request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request decoding error"})
		return
	}

	if !request.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	newGroupAccess, err := h.service.Create(c, userID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Group Access creating error"})
		return
	}

	h.respondWithDoc(newGroupAccess, c)
}

func (h *GroupAccessHandler) Get(c *gin.Context) {
	groupID := c.Param("group_id")

	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "groupId is incorrect"})
		return
	}

	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	accessGroup, err := h.service.Get(c, userID, &groupUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group access not found"})
		return
	}

	c.JSON(http.StatusOK, accessGroup)
}

func (h *GroupAccessHandler) respondWithDoc(doc *model.GroupAccess, c *gin.Context) {
	c.Header("Location", h.baseURL+"/v1/access/group/"+doc.GroupUUID.String())
	c.Header("ETag", doc.GroupUUID.String())

	c.JSON(http.StatusCreated, doc)
}
