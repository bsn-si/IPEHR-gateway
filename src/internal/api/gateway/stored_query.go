package gateway

import (
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

// Get
//
//	@Summary		Get list stored queries
//	@Description	Retrieves list of all stored queries on the system matched by qualified_query_name as pattern.
//	@Description	https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/Query/operation/definition_query_list
//	@Tags			DEFINITION
//	@Accept			json
//	@Produce		json
//	@Param			qualified_query_name	path		string	false	"If pattern should given be in the format of [{namespace}::]{query-name},  and  when  is  empty,  it  will  be  treated  as  "wildcard"  in  the  search."
//	@Param			Authorization			header		string	true	"Bearer AccessToken"
//	@Param			AuthUserId				header		string	true	"UserId"
//	@Param			EhrSystemId				header		string	true	"The identifier of the system, typically a reverse domain identifier"
//	@Success		200						{object}	[]model.StoredQuery
//	@Failure		500						"Is returned when an unexpected error occurs while processing a request"
//	@Router			/definition/query/{qualified_query_name} [get]
func (h *QueryHandler) ListStored(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")

	// Here we do not check the existence of the argument.
	// This does not satisfy the specification. https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/Query/operation/definition_query_list
	// But otherwise it is not clear how the client can get the full list of stored queries
	qName := c.Param("qualified_query_name")

	queryList, err := h.service.List(c, userID, systemID, qName)
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

// Get
//
//	@Summary		Get stored query by version
//	@Description	Retrieves the definition of a particular stored query (at specified version) and its associated metadata.
//	@Description	https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/Query/operation/definition_query_list
//	@Tags			DEFINITION
//	@Produce		json
//	@Param			qualified_query_name	path		string	false	"If pattern should given be in the format of [{namespace}::]{query-name},  and  when  is       empty,  it       will     be  treated  as    "wildcard"  in       the  search."
//	@Param			version					path		string	false	"A SEMVER version number. This can be a an exact version (e.g. 1.7.1),     or   a     pattern  as      partial  prefix,  in  a        form  of          {major}  or   {major}.{minor}  (e.g. 1 or 1.0),  in  which  case  the  highest  (latest)  version  matching  the  prefix  will  be  considered."
//	@Param			Authorization			header		string	true	"Bearer AccessToken"
//	@Param			AuthUserId				header		string	true	"UserId"
//	@Param			EhrSystemId				header		string	true	"The identifier of the system, typically a reverse domain identifier"
//	@Success		200						{object}	model.StoredQuery
//	@Failure		400						"Is returned when the request has invalid content."
//	@Failure		404						"Is returned when a stored query with {qualified_query_name}  and  {version}  does  not  exist."
//	@Failure		500						"Is returned when an unexpected error occurs while processing a request"
//	@Router			/definition/query/{qualified_query_name}/{version} [get]
func (h *QueryHandler) GetStoredByVersion(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")

	version := c.Param("version")

	v, err := base.NewVersionTreeID(version)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	qName := c.Param("qualified_query_name")
	if qName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "qualified_query_name is empty"})
		return
	}

	sq, err := h.service.GetByVersion(c, userID, systemID, qName, v)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Printf("StoredQuery service error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, sq)
}

// Store a query
//
//	@Summary	Stores a new query, or updates an existing query on the system
//	@Description
//	@Description	https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/Query/operation/definition_query_store.yaml
//	@Tags			DEFINITION
//	@Accept			json
//	@Produce		json
//	@Param			qualified_query_name	path		string		true	"If pattern should given be in the format of [{namespace}::]{query-name}, and when is empty, it will be treated as "wildcard" in the search."
//	@Param			query_type				query		string		true	"Parameter indicating the query language/type"
//	@Param			Authorization			header		string		true	"Bearer AccessToken"
//	@Param			AuthUserId				header		string		true	"UserId"
//	@Param			EhrSystemId				header		string		true	"The identifier of the system, typically a reverse domain identifier"
//	@Header			200						{string}	Location	"{baseUrl}/definition/query/org.openehr::compositions/1.0.1"
//	@Success		200						"Is returned when the query was successfully stored."
//	@Failure		400						"Is returned when the server was unable to store the query. This could be due to incorrect request body (could not be parsed, etc), unknown query type, etc."
//	@Failure		500						"Is returned when an unexpected error occurs while processing a request"
//	@Router			/definition/query/{qualified_query_name} [put]
func (h *QueryHandler) Store(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "header AuthUserId is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")

	qName := c.Param("qualified_query_name")
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

// Store a query version
//
//	@Summary	Stores a query, at a specified version, on the system.
//	@Description
//	@Description	https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/Query/operation/definition_query_store.yaml
//	@Tags			DEFINITION
//	@Accept			json
//	@Produce		json
//	@Param			qualified_query_name	path		string		true	"If pattern should given be in the format of [{namespace}::]{query-name},  and  when  is       empty,  it       will     be  treated  as    "wildcard"  in       the  search"
//	@Param			version					path		string		true	"A SEMVER version number. This can be a an exact version (e.g. 1.7.1),     or   a     pattern  as      partial  prefix,  in  a        form  of          {major}  or   {major}.{minor}  (e.g. 1 or 1.0),  in  which  case  the  highest  (latest)  version  matching  the  prefix  will  be  considered"
//	@Param			query_type				query		string		true	"Parameter indicating the query language/type"
//	@Param			Authorization			header		string		true	"Bearer AccessToken"
//	@Param			AuthUserId				header		string		true	"UserId"
//	@Header			200						{string}	Location	"{baseUrl}/definition/query/org.openehr::compositions/1.0.1"
//	@Success		200						"Is returned when the query was successfully stored"
//	@Failure		400						"Is returned when the server was unable to store the query. This could be due to incorrect request body (could not be parsed, etc),  unknown  query  type,  etc"
//	@Failure		409						"Is returned when a query with the given 'qualified_query_name' and 'version' already exists on the server"
//	@Failure		500						"Is returned when an unexpected error occurs while processing a request"
//	@Router			/definition/query/{qualified_query_name}/{version} [put]
func (h *QueryHandler) StoreVersion(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	systemID := c.GetString("ehrSystemID")

	qName := c.Param("qualified_query_name")
	if qName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "qualified_query_name is empty"})
		return
	}

	qType := c.GetString("query_type")
	if qType == "" {
		qType = model.QueryTypeAQL
	}

	version := c.Param("version")

	v, err := base.NewVersionTreeID(version)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	query, err := h.service.GetByVersion(c, userID, systemID, qName, v)
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		log.Printf("GetByVersion error: %v", err) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if query != nil {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}
	defer c.Request.Body.Close()

	if string(data) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
		return
	}

	if !h.service.Validate(data) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	reqID := c.GetString("reqID")

	sQ, err := h.service.StoreVersion(c, userID, systemID, reqID, qType, qName, v, string(data))
	if err != nil {
		log.Printf("StoredQuery service error: %s", err.Error()) // TODO replace to ErrorF after merge IPEHR-32

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("Location", h.baseURL+"/v1/definition/query/"+sQ.Name+"/"+v.String())

	c.Status(http.StatusOK)
}
