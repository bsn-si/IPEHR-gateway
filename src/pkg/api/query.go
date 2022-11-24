package api

import (
	"context"
	"encoding/json"
	"hms/gateway/pkg/docs/model/base"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/docs/model"
)

type QueryService interface {
	List(ctx context.Context, userID, qualifiedQueryName string) ([]*model.StoredQuery, error)
	GetByVersion(ctx context.Context, userID string, qualifiedQueryName string, version *base.VersionTreeID) (*model.StoredQuery, error)
	Validate(data []byte) bool
	Store(ctx context.Context, userID, systemID, reqID, qType, name, q string) (*model.StoredQuery, error)
	StoreVersion(ctx context.Context, userID, systemID, reqID, qType, name, version, q string) (*model.StoredQuery, error)
	StoreVersion(ctx context.Context, userID string, qType string, qualifiedQueryName string, version *base.VersionTreeID, q []byte) (model.StoredQuery, error)
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
