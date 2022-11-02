package api

import (
	"encoding/json"
	"io"
	"log"
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

// Create
// @Summary      Create new access group
// @Description  Creates new access group for use with part of users data and return this group
// @Tags         GROUP_ACCESS
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                          true  "Bearer <JWT>"
// @Param        AuthUserId     header    string                          true  "UserId UUID"
// @Param        Request        body      model.GroupAccessCreateRequest  true  "DTO with data to create group access"
// @Success      200            {object}  model.GroupAccess
// @Failure      400            "Is returned when the request has invalid content."
// @Failure      500            "Is returned when an unexpected error occurs while processing a request"
// @Router       /access/group [post]
func (h *GroupAccessHandler) Create(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
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

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

	newGroupAccess, err := h.service.Create(c, userID, &request)
	if err != nil {
		log.Println(err)
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
// @Param        group_id       path      string  true  "access group id (UUID). Example: 7d44b88c-4199-4bad-97dc-d78268e01398"
// @Param        Authorization  header    string  true  "Bearer <JWT>"
// @Param        AuthUserId     header    string  true  "UserId UUID"
// @Success      200            {object}  model.GroupAccess
// @Failure      400            "Is returned when the request has invalid content."
// @Failure      500            "Is returned when an unexpected error occurs while processing a request"
// @Router       /access/group/{group_id} [get]
func (h *GroupAccessHandler) Get(c *gin.Context) {
	groupID := c.Param("group_id")

	groupUUID, err := uuid.Parse(groupID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "groupId is incorrect"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
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
