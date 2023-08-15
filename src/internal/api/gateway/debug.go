package gateway

import (
	"context"
	"encoding/csv"
	"fmt"
	"net/http"
	"sort"
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

type req struct {
	reqID       string
	handlerName string
	method      string
	url         string
	start       time.Time
	end         time.Time
}

type DebugHandler struct {
	enabled bool
	storage Storage
	indexer EthIndexer

	requests map[string]req
	mx       sync.RWMutex
}

func NewDebugHandler(enabled bool, storage Storage, indexer EthIndexer) *DebugHandler {
	return &DebugHandler{
		enabled: enabled,
		storage: storage,
		indexer: indexer,

		requests: make(map[string]req),
		mx:       sync.RWMutex{},
	}
}

func (d *DebugHandler) addRequest(r req) {
	d.mx.Lock()
	d.requests[r.reqID] = r
	d.mx.Unlock()
}

func (d *DebugHandler) getRequests() []req {
	d.mx.RLock()

	rrs := make([]req, 0, len(d.requests))

	for _, r := range d.requests {
		rrs = append(rrs, r)
	}

	d.mx.RUnlock()

	sort.Slice(rrs, func(i, j int) bool {
		return rrs[i].start.Before(rrs[j].start)
	})

	return rrs
}

func (d *DebugHandler) getEthTransactions(c *gin.Context) {
	requests, err := d.storage.GetAllRequests(c.Request.Context())
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	result := [][]string{
		{
			"req_id",
			"time_start",
			"time_end",
			"duration",
			"handler_name",
			"method",
			"full_path",
			"req_kind",
			"req_status",
			"req_created_at",
			"req_updated_at",
			"req_duration",
			"tx_hash",
			"tx_kind",
			"tx_status",
			"tx_cumulativeGasUsed",
			"gas",
			"block_number",
			"transaction_index",
		},
	}

	for _, rr := range d.getRequests() {
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
				rr.handlerName,
				rr.method,
				rr.url,
				"",
				"",
				"",
				"",
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
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		for _, tx := range txs {
			r, err := d.indexer.GetTxReceipt(c.Request.Context(), tx.Hash)
			if err != nil {
				_ = c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			result = append(result, []string{
				rr.reqID,
				rr.start.Format("2006-01-02 15:04:05"),
				rr.end.Format("2006-01-02 15:04:05"),
				fmt.Sprintf("%v", rr.end.Sub(rr.start).Seconds()),
				rr.handlerName,
				rr.method,
				rr.url,
				req.Kind.String(),
				req.Status.String(),
				tx.CreatedAt.Format("2006-01-02 15:04:05"),
				tx.UpdatedAt.Format("2006-01-02 15:04:05"),
				fmt.Sprintf("%v", tx.UpdatedAt.Sub(tx.CreatedAt).Seconds()),
				tx.Hash,
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
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (d *DebugHandler) getAvgRequestsTime(c *gin.Context) {
	result := [][]string{
		{
			"hander",
			"avg_time",
		},
	}

	type t struct {
		count int
		sum   time.Duration
	}

	m := map[string]t{}

	for _, rr := range d.getRequests() {
		if _, ok := m[rr.handlerName]; !ok {
			m[rr.handlerName] = t{}
		}

		m[rr.handlerName] = t{
			count: m[rr.handlerName].count + 1,
			sum:   m[rr.handlerName].sum + rr.end.Sub(rr.start),
		}
	}

	for url, tt := range m {
		result = append(result, []string{
			url,
			fmt.Sprintf("%v", tt.sum.Seconds()/float64(tt.count)),
		})
	}

	csvWriter := csv.NewWriter(c.Writer)
	if err := csvWriter.WriteAll(result); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	csvWriter.Flush()

	if err := csvWriter.Error(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (d *DebugHandler) debugMiddleware(c *gin.Context) {
	if !d.enabled {
		c.Next()
		return
	}

	r := req{
		reqID:       c.GetString("reqID"),
		start:       time.Now(),
		handlerName: c.HandlerName(),
		method:      c.Request.Method,
		url:         c.FullPath(),
	}

	// serve the request to the next middleware
	c.Next()

	r.end = time.Now()

	d.addRequest(r)
}
