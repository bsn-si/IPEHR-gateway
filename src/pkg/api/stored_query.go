package api

import (
	"context"
	"log"
	"net/http"

	"hms/gateway/pkg/docs/model"

	"github.com/gin-gonic/gin"
)

type StoredQueryService interface {
	Get(ctx context.Context, userID string, qualifiedQueryName string) ([]model.StoredQuery, error)
}

type StoredQueryHandler struct {
	service StoredQueryService
}

func NewStoredQueryHandler(storedQueryService StoredQueryService) *StoredQueryHandler {
	return &StoredQueryHandler{
		service: storedQueryService,
	}
}

// Get TODO
// @Summary      Get list stored queries
// @Description  Retrieves list of all stored queries on the system matched by qualified_query_name as pattern.
// @Description  https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/Query/operation/definition_query_list
// @Tags         STORED_QUERY
// @Accept       json
// @Produce      json
// @Param        qualified_query_name    path      string  true  "If pattern should given be in the format of [{namespace}::]{query-name}, and when is empty, it will be treated as "wildcard" in the search."
// @Param        Authorization           header    string  true  "Bearer <JWT>"
// @Param        AuthUserId              header    string  true  "UserId UUID"
// @Success      200            {object}  []model.StoredQuery
// @Failure      500            "Is returned when an unexpected error occurs while processing a request"
// @Router       /access/group/{group_id} [get]
func (h *StoredQueryHandler) Get(c *gin.Context) {
	qName := c.Param("qualifiedQueryName")

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	if qName == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	queryList, err := h.service.Get(c, userID, qName)
	if err != nil {
		log.Printf("StoredQuery service error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(queryList) == 0 {
		queryList = make([]model.StoredQuery, 0)
	}

	c.JSON(http.StatusOK, queryList)
}
