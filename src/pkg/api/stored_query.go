package api

import (
	"io"
	"log"
	"net/http"

	"hms/gateway/pkg/docs/model"

	"github.com/gin-gonic/gin"
)

// Get
// @Summary      Get list stored queries
// @Description  Retrieves list of all stored queries on the system matched by qualified_query_name as pattern.
// @Description  https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/Query/operation/definition_query_list
// @Tags         QUERY
// @Accept       json
// @Produce      json
// @Param        qualified_query_name    path      string  false  "If pattern should given be in the format of [{namespace}::]{query-name}, and when is empty, it will be treated as "wildcard" in the search."
// @Param        Authorization           header    string  true  "Bearer AccessToken"
// @Param        AuthUserId              header    string  true  "UserId UUID"
// @Success      200            {object}  []model.StoredQuery
// @Failure      500            "Is returned when an unexpected error occurs while processing a request"
// @Router       /definition/query/{qualifiedQueryName} [get]
func (h *QueryHandler) ListStored(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	// Here we do not check the existence of the argument.
	// This does not satisfy the specification. https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/Query/operation/definition_query_list
	// But otherwise it is not clear how the client can get the full list of stored queries
	qName := c.Param("qualifiedQueryName")

	queryList, err := h.service.List(c, userID, qName)
	if err != nil {
		log.Printf("StoredQuery service error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if len(queryList) == 0 {
		queryList = make([]*model.StoredQuery, 0)
	}

	c.JSON(http.StatusOK, queryList)
}

// Store a query
// @Summary      Stores a new query, or updates an existing query on the system
// @Description
// @Description  https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/Query/operation/definition_query_store.yaml
// @Tags         QUERY
// @Accept       json
// @Produce      json
// @Param        qualified_query_name    path      string  true  "If pattern should given be in the format of [{namespace}::]{query-name}, and when is empty, it will be treated as "wildcard" in the search."
// @Param        query_type              query     string  true  "Parameter indicating the query language/type"
// @Param        Authorization           header    string  true  "Bearer AccessToken"
// @Param        AuthUserId              header    string  true  "UserId"
// @Param        EhrSystemId			 header    string  true  "The identifier of the system, typically a reverse domain identifier"
// @Header       200  {string}  Location "{baseUrl}/definition/query/org.openehr::compositions/1.0.1"
// @Success      200            "Is returned when the query was successfully stored."
// @Failure      400            "Is returned when the server was unable to store the query. This could be due to incorrect request body (could not be parsed, etc), unknown query type, etc."
// @Failure      500            "Is returned when an unexpected error occurs while processing a request"
// @Router       /definition/query/{qualifiedQueryName} [put]
func (h *QueryHandler) Store(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "header AuthUserId is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")
	if systemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "header EhrSystemId is empty"})
		return
	}

	qName := c.Param("qualifiedQueryName")
	if qName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "qualified_query_name is empty"})
		return
	}

	qType := c.GetString("query_type")
	if qType == "" {
		qType = model.QueryTypeAQL
	}

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}

	defer c.Request.Body.Close()

	if len(data) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
		return
	}

	if !h.service.Validate(data) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	reqID := c.GetString("reqID")

	sQ, err := h.service.Store(c, userID, systemID, reqID, qType, qName, string(data))
	if err != nil {
		log.Printf("StoredQuery service error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("Location", h.baseURL+"/v1/definition/query/"+sQ.Name+"/"+sQ.Version)

	c.Status(http.StatusOK)
}
