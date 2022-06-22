package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"hms/gateway/pkg/common/fake_data"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/query"
)

type QueryHandler struct {
	cfg     *config.Config
	service *query.QueryService
}

func NewQueryHandler(docService *service.DefaultDocumentService, cfg *config.Config) *QueryHandler {
	return &QueryHandler{
		cfg:     cfg,
		service: query.NewQueryService(docService),
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
// @Param    AuthUserId  header    string              true  "UserId UUID"
// @Param    Request     body      model.QueryRequest  true  "Query Request"
// @Success  200         {object}  model.QueryResponse
// @Header   201         {string}  ETag  "A unique identifier of the resultSet. Example: cdbb5db1-e466-4429-a9e5-bf80a54e120b"
// @Failure  400         "Is returned when the server was unable to execute the query due to invalid input, e.g. a request with missing `q` parameter or an invalid query syntax."
// @Failure  408         "Is returned when there is a query execution timeout (i.e. maximum query execution time reached, therefore the server aborted the execution of the query)."
// @Failure  500         "Is returned when an unexpected error occurs while processing a request"
// @Router   /query/aql [post]
func (h QueryHandler) ExecPost(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
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

	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is empty"})
		return
	}

	//TODO make real job
	c.Data(http.StatusOK, "application/json", fake_data.QueryExecResponse(queryRequest.Query))
}
