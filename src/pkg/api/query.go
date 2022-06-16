package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"hms/gateway/pkg/common/fake_data"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/query"

	"github.com/gin-gonic/gin"
)

type QueryHandler struct {
	*query.QueryService
}

func NewQueryHandler(docService *service.DefaultDocumentService, cfg *config.Config) *QueryHandler {
	return &QueryHandler{
		query.NewQueryService(docService, cfg),
	}
}

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
	c.JSON(http.StatusOK, fake_data.QueryExecResponse(queryRequest.Query))
}
