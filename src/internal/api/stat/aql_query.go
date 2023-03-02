package stat

import (
	"context"
	"log"
	"net/http"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/gin-gonic/gin"
)

type AQLQuerier interface {
	ExecQuery(ctx context.Context, query *model.QueryRequest) (*model.QueryResponse, error)
}

type aqlQueryAPI struct {
	querier AQLQuerier
}

func newAQLQueryAPI(querier AQLQuerier) *aqlQueryAPI {
	return &aqlQueryAPI{
		querier: querier,
	}
}

func (api *aqlQueryAPI) QueryHandler(c *gin.Context) {
	req := model.QueryRequest{
		QueryParameters: map[string]interface{}{},
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request validation error"})
		return
	}

	resp, err := api.querier.ExecQuery(c.Request.Context(), &req)
	if err != nil {
		log.Printf("cannot exec query: %v", err)

		if errors.Is(err, errors.ErrTimeout) {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "timeout exceeded"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
