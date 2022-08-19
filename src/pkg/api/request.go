package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/common"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/request"
)

type RequestHandler struct {
	service *request.Service
}

func NewRequestHandler(docService *service.DefaultDocumentService) *RequestHandler {
	return &RequestHandler{
		service: request.NewRequestService(docService),
	}
}

// GetAll
// @Summary      Get list of transactions requests by authorized user
// @Description  It is returning only transactions which in progress
// @Description
// @Tags     REQUEST
// @Accept   json
// @Produce  json
// @Param    AuthUserId  header    string  true  "UserId UUID"
// @Param    limit       query     string  true  "default: 10"
// @Param    offset      query     string  true  "id namespace. Example: examples"
// @Success  200         {object}  processing.RequestsResult
// @Failure  400         "Is returned when userId is empty"
// @Failure  404         "Is returned when requests not exist"
// @Failure  500         "Is returned when an unexpected error occurs while processing a request"
// @Router   /requests/ [get]
func (h RequestHandler) GetAll(c *gin.Context) {
	reqLimit := c.Query("limit")
	reqOffset := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(reqLimit)
	if err != nil || limit > common.PageLimit || limit <= 0 {
		limit = common.PageLimit
	}

	offset, err := strconv.Atoi(reqOffset)
	if err != nil || offset <= 0 {
		offset = 0
	}

	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	data, err := h.service.Doc.Proc.GetRequests(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request error"})
		return
	}

	if data == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}

// GetByID
// @Summary      Get list of transactions by certain request id for authorized user
// @Description  It's returning only transactions which in progress
// @Description
// @Tags     REQUEST
// @Accept   json
// @Produce  json
// @Param    AuthUserId  header    string  true  "UserId UUID"
// @Param    request_id  path      string  true  "Unique id of request"
// @Success  200         {object}  processing.RequestResult
// @Failure  400         "Is returned when userId or request_id is empty"
// @Failure  404         "Is returned when requests not exist"
// @Failure  500         "Is returned when an unexpected error occurs while processing a request"
// @Router   /requests/{request_id} [get]
func (h RequestHandler) GetByID(c *gin.Context) {
	userID := c.GetString("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	reqID := c.Param("reqId")
	if reqID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "requestId is empty"})
		return
	}

	data, err := h.service.Doc.Proc.GetRequest(reqID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request error"})
		return
	}

	if data == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}