package api

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
)

type QueryService interface {
	List(ctx context.Context, userID, systemID, qualifiedQueryName string) ([]*model.StoredQuery, error)
	GetByVersion(ctx context.Context, userID, systemID, name string, version *base.VersionTreeID) (*model.StoredQuery, error)
	Validate(data []byte) bool
	Store(ctx context.Context, userID, systemID, reqID, qType, name, q string) (*model.StoredQuery, error)
	StoreVersion(ctx context.Context, userID, systemID, reqID, qType, name string, version *base.VersionTreeID, q string) (*model.StoredQuery, error)

	ExecStoredQuery(ctx context.Context, userID, systemID, qualifiedQueryName string, query *model.QueryRequest) (*model.QueryResponse, error)
}

type QueryHandler struct {
	service QueryService
	baseURL string
}

func NewQueryHandler(queryService QueryService, baseURL string) *QueryHandler {
	return &QueryHandler{
		service: queryService,
		baseURL: baseURL,
	}
}

// ExecPost
// @Summary      Execute ad-hoc (non-stored) AQL query
// @Description  Work in progress...
// @Description  Execute ad-hoc query, supplied by q attribute, fetching fetch numbers of rows from offset and passing query_parameters to the underlying query engine.
// @Description  See also details on usage of [query parameters](https://specifications.openehr.org/releases/ITS-REST/Release-1.0.2/query.html#requirements-common-headers-and-query-parameters).
// @Description
// @Tags     QUERY
// @Accept   json
// @Produce  json
// @Param    Authorization  header    string              true  "Bearer AccessToken"
// @Param    AuthUserId     header    string              true  "UserId UUID"
// @Param    Request        body      model.QueryRequest  true  "Query Request"
// @Success  200            {object}  model.QueryResponse
// @Header   201            {string}  ETag  "A unique identifier of the resultSet. Example: cdbb5db1-e466-4429-a9e5-bf80a54e120b"
// @Failure  400            "Is returned when the server was unable to execute the query due to invalid input, e.g. a request with missing `q` parameter or an invalid query syntax."
// @Failure  408            "Is returned when there is a query execution timeout (i.e. maximum query execution time reached, therefore the server aborted the execution of the query)."
// @Failure  500            "Is returned when an unexpected error occurs while processing a request"
// @Router   /query/aql [post]
func (h QueryHandler) ExecPost(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}
	defer c.Request.Body.Close()

	var queryRequest model.QueryRequest
	if err = json.Unmarshal(data, &queryRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	if !queryRequest.Validate() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request validation error"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is empty"})
		return
	}

	//TODO make real job
	c.Data(http.StatusOK, "application/json", fakeData.QueryExecResponse(queryRequest.Query))
}

// Get
// @Summary      Execute stored AQL
// @Description  Execute a stored query, identified by the supplied qualified_query_name (at latest version), fetching fetch numbers of rows from offset and passing query_parameters to the underlying query engine.
// @Description  See also details on usage of [query parameters](https://specifications.openehr.org/releases/ITS-REST/latest/query.html#tag/Request/Common-Headers-and-Query-Parameters).
// @Description  Queries can be stored or, once stored, their definition can be retrieved using the [definition endpoint](https://specifications.openehr.org/releases/ITS-REST/latest/definition.html#tag/Query).
// @Description  https://specifications.openehr.org/releases/ITS-REST/latest/query.html#tag/Query/operation/query_execute_stored_query
// @Tags         QUERY
// @Accept       json
// @Produce      json
// @Param        Authorization         header    string              true  "Bearer AccessToken"
// @Param        AuthUserId            header    string              true  "UserId UUID"
// @Param        qualified_query_name  path      string  true   "If pattern should given be in the format of [{namespace}::]{query-name},  and  when  is       empty,  it       will     be  treated  as    "wildcard"  in       the  search."
// @Param        ehr_id            	   query     string  false  "An optional parameter to execute the query within an EHR context."
// @Param        offset            	   query     string  false  "The row number in result-set to start result-set from (0-based), default is 0."
// @Param        fetch            	   query     string  false  "Number of rows to fetch (the default depends on the implementation)."
// @Success      200                   {object}  model.QueryResponse
// @Header       200                   {string}  ETag  "A unique identifier of the resultSet. Example: cdbb5db1-e466-4429-a9e5-bf80a54e120b"
// @Failure      400                   "Is returned when the server was unable to execute the query due to invalid input, e.g. a required parameter is missing, or at least one of the parameters has invalid syntax"
// @Failure      404                   "Is returned when a stored query with qualified_query_name does not exists."
// @Failure      408                   "Is returned when there is a query execution timeout"
// @Router       /query/{qualified_query_name} [get]
func (h QueryHandler) ExecStoredQuery(c *gin.Context) {
	userID := c.GetString("userID")
	systemID := c.GetString("ehrSystemID")

	qualifiedQueryName := c.Param("qualified_query_name")

	m := map[string]string{}

	if err := c.BindQuery(&m); err != nil {
		log.Printf("cannot bind query params to map: %f", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body error"})
		return
	}

	req := &model.QueryRequest{
		QueryParameters: map[string]interface{}{},
	}

	for key, val := range m {
		if key == "ehr_id" {
			ehrID, err := uuid.Parse(val)
			if err != nil {
				log.Printf("cannot parse ehr_id uuid: %f", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "ehr_id bad format"})
				return
			}

			req.QueryParameters["ehr_id"] = ehrID

			continue
		}

		if key == "offset" {
			offset, err := strconv.Atoi(val)
			if err != nil {
				log.Printf("cannot parse 'offset' from string: %f", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "offset bad format"})
				return
			}

			if offset < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "offset cannot be less than 0"})
				return
			}

			req.Offset = offset

			continue
		}

		if key == "fetch" {
			fetch, err := strconv.Atoi(val)
			if err != nil {
				log.Printf("cannot parse 'fetch' from string: %f", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "fetch bad format"})
				return
			}

			if fetch < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "fetch cannot be less than 0"})
				return
			}

			req.Fetch = fetch

			continue
		}

		req.QueryParameters[key] = val
	}

	resp, err := h.service.ExecStoredQuery(c, userID, systemID, qualifiedQueryName, req)
	if err != nil {
		log.Printf("cannot exec stored query: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
