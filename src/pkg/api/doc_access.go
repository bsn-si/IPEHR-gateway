package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-cid"

	"hms/gateway/pkg/access"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/docAccess"
)

type DocAccessHandler struct {
	service *docAccess.Service
}

func NewDocAccessHandler(docService *service.DefaultDocumentService) *DocAccessHandler {
	return &DocAccessHandler{
		service: docAccess.NewService(docService),
	}
}

// Create
// @Summary      Create new access group
// @Description  Creates new access group for use with part of users data and return this group
// @Tags         ACCESS
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                          true  "Bearer <JWT>"
// @Param        AuthUserId     header    string                          true  "UserId UUID"
// @Param        Request        body      model.GroupAccessCreateRequest  true  "DTO with data to create group access"
// @Success      200            {object}  model.GroupAccess
// @Failure      400            "Is returned when the request has invalid content."
// @Failure      500            "Is returned when an unexpected error occurs while processing a request"
// @Router       /access/document/manage [post]
func (h *DocAccessHandler) Manage(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}
	defer c.Request.Body.Close()

	var req model.DocAccessManageRequest
	if err = json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	switch access.Level(req.AccessLevel) {
	case access.NoAccess, access.Owner, access.Admin, access.Read:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "AccessLevel is incorrect"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

	CID, err := cid.Parse(req.CID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CID is incorrect"})
		return
	}

	reqID := c.GetString("reqId")
	if reqID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requestId is empty"})
		return
	}

	if err := h.service.Set(c, userID, req.ToUserID, reqID, &CID, req.AccessLevel); err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
