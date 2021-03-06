package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/groupAccess"
)

type GroupAccessHandler struct {
	cfg     *config.Config
	service *groupAccess.Service
}

func NewGroupAccessHandler(docService *service.DefaultDocumentService, cfg *config.Config) *GroupAccessHandler {
	return &GroupAccessHandler{
		cfg:     cfg,
		service: groupAccess.NewGroupAccessService(docService, cfg),
	}
}

// Create
// @Summary      Create new access group
// @Description  Creates new access group for use with part of users data and return this group
// @Tags         GROUP_ACCESS
// @Accept       json
// @Produce      json
// @Param        AuthUserId  header    string                          true  "UserId UUID"
// @Param        Request     body      model.GroupAccessCreateRequest  true  "DTO with data to create group access"
// @Success      200         {object}  model.GroupAccess
// @Failure      400         "Is returned when the request has invalid content."
// @Failure      500         "Is returned when an unexpected error occurs while processing a request"
// @Router       /v1/access/group [post]
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

	newGroupAccess, err := h.service.Create(userID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Group Access creating error"})
		return
	}

	h.respondWithDoc(newGroupAccess, c)
}

// Get
// @Summary      Get access group
// @Description  Return access group object
// @Tags         GROUP_ACCESS
// @Accept       json
// @Produce      json
// @Param        group_id    path      string  true  "access group id (UUID). Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param        AuthUserId  header    string  true  "UserId UUID"
// @Success      200         {object}  model.GroupAccess
// @Failure      400         "Is returned when the request has invalid content."
// @Failure      500         "Is returned when an unexpected error occurs while processing a request"
// @Router       /v1/access/group/{group_id} [get]
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

	accessGroup, err := h.service.Get(userID, &groupUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group access not found"})
		return
	}

	c.JSON(http.StatusOK, accessGroup)
}

func (h *GroupAccessHandler) respondWithDoc(doc *model.GroupAccess, c *gin.Context) {
	c.Header("Location", h.cfg.BaseURL+"/v1/access/group/"+doc.GroupUUID.String())
	c.Header("ETag", doc.GroupUUID.String())

	c.JSON(http.StatusCreated, doc)
}
