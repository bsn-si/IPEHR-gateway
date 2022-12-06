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
	"hms/gateway/pkg/errors"
)

type DocAccessHandler struct {
	service *docAccess.Service
}

func NewDocAccessHandler(docService *service.DefaultDocumentService) *DocAccessHandler {
	return &DocAccessHandler{
		service: docAccess.NewService(docService),
	}
}

// List
// @Summary      Get a document access list
// @Description  Returns the list of documents available to the user
// @Tags         ACCESS
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer AccessToken"
// @Param        AuthUserId     header  string  true  "UserId UUID"
// @Param    EhrSystemId    header    string                 true   "The identifier of the system, typically a reverse domain identifier"
// @Success      200            ""
// @Failure      400            "Is returned when the request has invalid content."
// @Failure      500            "Is returned when an unexpected error occurs while processing a request"
// @Router       /access/document/ [get]
func (h *DocAccessHandler) List(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

	acl, err := h.service.List(c, userID)
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	type item struct {
		CID         string `json:"CID"`
		Level       string `json:"level"`
		Description string `json:"description"`
	}

	resp := []item{}

	for _, a := range acl {
		CID, ok := a.Fields["CID"]
		if !ok {
			log.Printf("Error: Doc ACL field CID required")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		level, ok := a.Fields["level"]
		if !ok {
			log.Printf("Error: Doc ACL field 'level' required")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		resp = append(resp, item{
			CID:   string(CID),
			Level: string(level),
		})
	}

	c.JSON(http.StatusOK, resp)
}

// Set
// @Summary      Set user access to the document
// @Description  Sets access to the document with the specified CID for the user with the userID.
// @Description  Possible access levels: `owner`, `admin`, `read`, `noAccess`
// @Tags         ACCESS
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string                     true  "Bearer AccessToken"
// @Param        AuthUserId     header  string                     true  "UserId UUID"
// @Param    EhrSystemId    header    string                 true   "The identifier of the system, typically a reverse domain identifier"
// @Param        Request        body    model.DocAccessSetRequest  true  "DTO with data to create group access"
// @Success      200            "Indicates that the request to change the level of access to the document was successfully created"
// @Failure      400            "Is returned when the request has invalid content."
// @Failure      404            "Is returned when the userID for which access is set is not found "
// @Failure      500            "Is returned when an unexpected error occurs while processing a request"
// @Router       /access/document [post]
func (h *DocAccessHandler) Set(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}
	defer c.Request.Body.Close()

	var req model.DocAccessSetRequest
	if err = json.Unmarshal(data, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	level := access.LevelFromString(req.AccessLevel)
	if level == access.Unknown {
		c.JSON(http.StatusBadRequest, gin.H{"error": "AccessLevel is incorrect"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "systemID is empty"})
		return
	}

	CID, err := cid.Parse(req.CID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CID is incorrect"})
		return
	}

	reqID := c.GetString("reqID")
	if reqID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requestId is empty"})
		return
	}

	if err := h.service.Set(c, userID, systemID, req.UserID, reqID, &CID, level); err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Println(err)

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}
