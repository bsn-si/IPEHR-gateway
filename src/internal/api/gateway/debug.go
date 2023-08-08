package gateway

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
)

type Storage interface {
	GetAllRequests(ctx context.Context) ([]*processing.Request, error)
	GetEthTransactionsForRequest(ctx context.Context, reqID string) ([]*processing.EthereumTx, error)
	GetAllEthTransactions(ctx context.Context) ([]*processing.EthereumTx, error)
}

type EthIndexer interface {
	GetTxReceipt(ctx context.Context, hash string) (*types.Receipt, error)
}

type DebugHandler struct {
	storage Storage
	indexer EthIndexer
}

func NewDebugHandler(storage Storage, indexer EthIndexer) *DebugHandler {
	return &DebugHandler{
		storage: storage,
		indexer: indexer,
	}
}

func (d *DebugHandler) GetEthTransactions(c *gin.Context) {
	requests, err := d.storage.GetAllRequests(c.Request.Context())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	result := [][]string{
		{
			"req_id",
			"time_start",
			"time_end",
			"duration",
			"full_path",
			"req_kind",
			"req_status",
			"tx_kind",
			"tx_status",
			"tx_cumulativeGasUsed",
			"gas",
			"block_number",
			"transaction_index",
		},
	}

	rrs := []req{}

	mx.RLock()
	for _, rr := range rMap {
		rrs = append(rrs, rr)
	}
	mx.RUnlock()

	for _, rr := range rrs {
		var req *processing.Request

		for _, r := range requests {
			if r.ReqID == rr.reqID || r.ReqID == rr.reqID+"_register" {
				req = r
				break
			}
		}

		if req == nil {
			result = append(result, []string{
				rr.reqID,
				rr.start.Format("2006-01-02 15:04:05"),
				rr.end.Format("2006-01-02 15:04:05"),
				fmt.Sprintf("%v", rr.end.Sub(rr.start).Seconds()),
				rr.url,
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
			})

			continue
		}

		txs, err := d.storage.GetEthTransactionsForRequest(c.Request.Context(), req.ReqID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		for _, tx := range txs {
			r, err := d.indexer.GetTxReceipt(c.Request.Context(), tx.Hash)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			result = append(result, []string{
				rr.reqID,
				rr.start.Format("2006-01-02 15:04:05"),
				rr.end.Format("2006-01-02 15:04:05"),
				fmt.Sprintf("%v", rr.end.Sub(rr.start).Seconds()),
				rr.url,
				req.Kind.String(),
				req.Status.String(),
				tx.Kind.String(),
				fmt.Sprintf("%d", r.Status),
				fmt.Sprintf("%d", r.CumulativeGasUsed),
				fmt.Sprintf("%d", r.GasUsed),
				r.BlockNumber.Text(10),
				fmt.Sprintf("%d", r.TransactionIndex),
			})
		}
	}

	// convert result to csv and send it
	csvWriter := csv.NewWriter(c.Writer)
	if err := csvWriter.WriteAll(result); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

type req struct {
	reqID string
	url   string
	start time.Time
	end   time.Time
}

var (
	rMap = map[string]req{}

	mx = sync.RWMutex{}
)

func middleware(c *gin.Context) {
	r := req{
		reqID: c.GetString("reqID"),
		start: time.Now(),
		url:   c.FullPath(),
	}

	// serve the request to the next middleware
	c.Next()

	r.end = time.Now()

	mx.Lock()
	rMap[r.reqID] = r
	mx.Unlock()
}
