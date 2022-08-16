package processing

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/ipfs/go-cid"
	"gorm.io/gorm"

	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/storage/filecoin"
)

type (
	Status      uint8
	TxKind      uint8
	RequestKind uint8

	Proc struct {
		db             *gorm.DB
		ethClient      *ethclient.Client
		filecoinClient *filecoin.Client
		lock           bool
		done           chan bool
	}

	Request struct {
		ReqID        string      `gorm:"index:idx_request,unique"`
		Kind         RequestKind `gorm:"index:idx_request,unique"`
		Status       Status
		UserID       string
		EhrUUID      string
		CID          string
		DealCID      string
		MinerAddress string
		BaseUIDHash  string
		Version      string
		//Txs          []Tx `gorm:"foreignKey:ReqID"`
	}

	Tx struct {
		gorm.Model
		Time    time.Time
		ReqID   string
		Kind    TxKind
		Hash    string
		Status  Status
		Comment string
	}
)

const (
	StatusFailed     Status = 0
	StatusSuccess    Status = 1
	StatusPending    Status = 2
	StatusProcessing Status = 3

	RequestEhrCreate RequestKind = iota
	RequestEhrGetBySubject
	RequestEhrGetByID
	RequestEhrStatusCreate
	RequestEhrStatusUpdate
	RequestEhrStatusGetByID
	RequestEhrStatusGetByTime
	RequestCompositionCreate
	RequestCompositionUpdate
	RequestCompositionGetByID
	RequestCompositionDelete

	TxSetEhrUser        TxKind = 0
	TxSetEhrBySubject   TxKind = 1
	TxSetEhrDocs        TxKind = 2
	TxSetDocAccess      TxKind = 3
	TxDeleteDoc         TxKind = 4
	TxFilecoinStartDeal TxKind = 5
)

func (s Status) String() string {
	switch s {
	case StatusFailed:
		return "Failed"
	case StatusSuccess:
		return "Success"
	case StatusPending:
		return "Pending"
	case StatusProcessing:
		return "Processing"
	default:
		// nolint
		return "Unknown"
	}
}

func (k TxKind) String() string {
	switch k {
	case TxSetEhrUser:
		return "SetEhrUser"
	case TxSetEhrBySubject:
		return "SetEhrBySubject"
	case TxSetEhrDocs:
		return "SetEhrDocs"
	case TxSetDocAccess:
		return "SetDocAccess"
	case TxDeleteDoc:
		return "DeleteDoc"
	case TxFilecoinStartDeal:
		return "FilecoinStartDeal"
	default:
		return "Unknown"
	}
}

func (k RequestKind) String() string {
	switch k {
	case RequestEhrCreate:
		return "EhrCreate"
	case RequestEhrGetBySubject:
		return "EhrGetBySubject"
	case RequestEhrGetByID:
		return "EhrGetByID"
	case RequestEhrStatusCreate:
		return "EhrStatusCreate"
	case RequestEhrStatusUpdate:
		return "EhrStatusUpdate"
	case RequestEhrStatusGetByID:
		return "EhrStatusGetByID"
	case RequestEhrStatusGetByTime:
		return "EhrStatusGetByTime"
	case RequestCompositionCreate:
		return "CompositionCreate"
	case RequestCompositionUpdate:
		return "CompositionUpdate"
	case RequestCompositionGetByID:
		return "CompositionGetByID"
	case RequestCompositionDelete:
		return "CompositionDelete"
	default:
		return "Unknown"
	}
}

func New(db *gorm.DB, ethClient *ethclient.Client, filecoinClient *filecoin.Client) *Proc {
	return &Proc{
		db:             db,
		ethClient:      ethClient,
		filecoinClient: filecoinClient,
		done:           make(chan bool),
	}
}

func (p *Proc) AddRequest(req *Request) error {
	if result := p.db.Create(req); result.Error != nil {
		return fmt.Errorf("db.Create error: %w", result.Error)
	}

	return nil
}

func (p *Proc) AddTx(reqID, txHash, comment string, kind TxKind, status Status) error {
	if reqID == "" {
		return fmt.Errorf("%w reqID %s", errors.ErrIncorrectRequest, reqID)
	}

	var req Request

	result := p.db.Model(&req).First(&req, "req_id = ?", reqID)
	if result.Error != nil {
		return fmt.Errorf("db request get error: %w reqID %s", result.Error, reqID)
	}

	result = p.db.Create(&Tx{
		Time:    time.Now(),
		ReqID:   reqID,
		Kind:    kind,
		Hash:    txHash,
		Status:  status,
		Comment: comment,
	})
	if result.Error != nil {
		return fmt.Errorf("db.Create error: %w", result.Error)
	}

	return nil
}

func (p *Proc) Start() {
	tickerBlockchain := time.NewTicker(5 * time.Second)
	tickerFilecoin := time.NewTicker(5 * time.Minute)

	go func() {
		for {
			select {
			case <-tickerBlockchain.C:
				if !p.lock {
					p.execBlockchain()
				}
			case <-tickerFilecoin.C:
				//p.execFilecoin()
			case <-p.done:
				return
			}
		}
	}()
}

func (p *Proc) Stop() {
	p.done <- true
}

func (p *Proc) execBlockchain() {
	p.lock = true
	defer func() { p.lock = false }()

	txKinds := []TxKind{
		TxSetEhrUser,
		TxSetEhrBySubject,
		TxSetEhrDocs,
		TxSetDocAccess,
		TxDeleteDoc,
	}

	statuses := []Status{
		StatusPending,
		StatusProcessing,
	}

	for {
		tx := Tx{}

		result := p.db.Model(&tx).Find(&tx, "kind IN ? AND status IN ?", txKinds, statuses).Limit(1)
		if result.Error != nil || result.RowsAffected == 0 {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		h := common.HexToHash(tx.Hash)

		var status Status

		receipt, err := p.ethClient.TransactionReceipt(ctx, h)
		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				status = StatusPending
			} else {
				log.Printf("TransactionReceipt error: %v txHash %s", err, h.String())
				break
			}
		} else {
			status = Status(receipt.Status)
		}

		log.Println("Tx", tx.Hash, "status", status)

		if status != tx.Status {
			if result = p.db.Model(&tx).Update("status", status); result.Error != nil {
				log.Println("db.Update error:", result.Error)
				break
			}

			done, err := p.isRequestDone(tx.ReqID)
			if err != nil {
				log.Println("checkRequestDone error:", err, "reqID", tx.ReqID)
				break
			}

			if done {
				// update tx statuses
				result = p.db.Exec("UPDATE txes SET status = ? WHERE req_id = ?", StatusSuccess, tx.ReqID)
				if result.Error != nil {
					//log.Printf("db txs update error: %v reqID %s", result.Error, tx.ReqID)
					break
				}

				// update request status
				result = p.db.Exec("UPDATE requests SET status = ? WHERE req_id = ?", StatusSuccess, tx.ReqID)
				if result.Error != nil {
					//log.Printf("db request delete error: %v reqID %s", result.Error, tx.ReqID)
					break
				}
			}

			continue
		}

		break
	}
}

// nolint
func (p *Proc) execFilecoin() {
	txKinds := []TxKind{
		TxFilecoinStartDeal,
	}

	statuses := []Status{
		StatusPending,
		StatusProcessing,
	}

	for {
		tx := Tx{}

		result := p.db.Where("kind IN ? AND status IN ?", txKinds, statuses).First(&tx)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		dealCID, err := cid.Parse(tx.Hash)
		if err != nil {
			log.Println("cid.Parse error:", err, "tx.Hash:", tx.Hash)
			break
		}

		dealStatus, err := p.filecoinClient.GetDealStatus(ctx, &dealCID)
		if err != nil {
			log.Println("filecoinClient.GetDealStatus error:", err, "dealCID", dealCID.String())
			break
		}

		switch tx.Status {
		case StatusProcessing:
			if dealStatus == storagemarket.StorageDealActive {
				tx.Status = StatusSuccess
				if result = p.db.Save(&tx); result.Error != nil {
					log.Println("db.Save error:", result.Error)
					return
				}
			}
		case StatusPending, StatusSuccess:
		}
	}
}

func (p *Proc) isRequestDone(reqID string) (bool, error) {
	request := Request{}

	result := p.db.Find(&request, "req_id = ?", reqID).Limit(1)
	if result.Error != nil {
		return false, fmt.Errorf("db.Find error: %w", result.Error)
	}

	var txsToCheck []TxKind

	switch request.Kind {
	case RequestEhrCreate:
		txsToCheck = []TxKind{
			TxSetEhrUser,
			TxSetEhrBySubject,
			TxSetEhrDocs,
			TxSetDocAccess,
			//TxFilecoinStartDeal,
		}
	case RequestEhrGetBySubject:
	case RequestEhrGetByID:
	case RequestEhrStatusCreate:
	case RequestEhrStatusUpdate:
	case RequestEhrStatusGetByID:
	case RequestEhrStatusGetByTime:
	case RequestCompositionCreate:
		txsToCheck = []TxKind{
			TxSetEhrDocs,
			TxSetDocAccess,
			//TxFilecoinStartDeal,
		}
	case RequestCompositionUpdate:
		txsToCheck = []TxKind{
			TxSetEhrDocs,
			TxSetDocAccess,
			//TxFilecoinStartDeal,
		}
	case RequestCompositionGetByID:
	case RequestCompositionDelete:
	default:
		return false, fmt.Errorf("%w Unknown request.Kind %d", errors.ErrCustom, request.Kind)
	}

	for _, txKind := range txsToCheck {
		var txs []Tx

		result := p.db.Model(&Tx{}).Where("req_id = ? AND kind = ?", reqID, txKind).Find(&txs)
		if result.Error != nil {
			return false, fmt.Errorf("db Find tx error: %w reqID %s txKind %d", result.Error, reqID, txKind)
		}

		// а он не может быть failed? что тогда реквест вообще не закроется?
		for _, tx := range txs {
			if tx.Status != StatusSuccess {
				return false, nil
			}
		}
	}

	return true, nil
}

func (p *Proc) GetRequestByID(reqID string, userID string) ([]byte, error) {
	var requestData Request

	result := p.db.Model(&Request{}).
		Where(Request{UserID: userID}).
		Where(Request{ReqID: reqID}).
		//Preload("BlockchainTxs").
		Find(&requestData)
	if result.Error != nil {
		return nil, fmt.Errorf("GetRequestByID error: %w reqID %s userID %s", result.Error, reqID, userID)
	}

	if result.RowsAffected == 0 {
		return nil, nil
	}

	resultBytes, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("GetRequestByID marshal error: %w", err)
	}

	return resultBytes, nil
}

type TxResult struct {
	Kind   string `json:"kind"`
	Status string `json:"status"`
	Hash   string `json:"hash"`
}

type DocResult struct {
	Kind string `json:"kind"`
	CID  string `json:"cid"`
}

type RequestResult struct {
	Status string       `json:"status"`
	Docs   []*DocResult `json:"docs"`
	Txs    []*TxResult  `json:"txs"`
}

func (p *Proc) GetRequests(userID string, limit, offset int) ([]byte, error) {
	result := make(map[string]*RequestResult)

	query := `SELECT r.req_id, r.status, r.kind, r.c_id, t.kind, t.hash, t.status
			FROM requests r, txes t
			WHERE 
				r.user_id = ? AND
				r.req_id = t.req_id`

	rows, err := p.db.Raw(query, userID).Rows()
	if err != nil {
		return nil, fmt.Errorf("GetRequests error: %w, userID %s", err, userID)
	}
	defer rows.Close()

	var (
		reqID     string
		reqStatus uint8
		reqKind   uint8
		reqCID    string
		txStatus  uint8
		txKind    uint8
		txHash    string
	)

	for rows.Next() {
		err = rows.Scan(&reqID, &reqStatus, &reqKind, &reqCID, &txKind, &txHash, &txStatus)
		if err != nil {
			return nil, fmt.Errorf("db.Scan error: %w", err)
		}

		req, ok := result[reqID]
		if !ok {
			req = &RequestResult{
				Status: Status(reqStatus).String(),
				Docs:   []*DocResult{},
				Txs:    []*TxResult{},
			}

			result[reqID] = req
		}

		exists := false

		for _, tx := range req.Txs {
			if tx.Hash == txHash {
				exists = true
				break
			}
		}

		if !exists {
			req.Txs = append(req.Txs, &TxResult{
				Kind:   TxKind(txKind).String(),
				Status: Status(txStatus).String(),
				Hash:   txHash,
			})
		}

		if reqCID == "" {
			continue
		}

		exists = false

		for _, doc := range req.Docs {
			if doc.CID == reqCID {
				exists = true
				break
			}
		}

		if !exists {
			req.Docs = append(req.Docs, &DocResult{
				Kind: RequestKind(reqKind).String(),
				CID:  reqCID,
			})
		}
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("GetRequestByID marshal error: %w", err)
	}

	return resultBytes, nil
}

func (p *Proc) GetRequest(reqID string) ([]byte, error) {
	var result RequestResult

	query := `SELECT r.status, r.kind, r.c_id, t.kind, t.hash, t.status
			FROM requests r, txes t
			WHERE 
				r.req_id = ? AND
				r.req_id = t.req_id`

	rows, err := p.db.Raw(query, reqID).Rows()
	if err != nil {
		return nil, fmt.Errorf("GetRequest error: %w, userID %s", err, reqID)
	}
	defer rows.Close()

	var (
		reqStatus uint8
		reqKind   uint8
		reqCID    string
		txStatus  uint8
		txKind    uint8
		txHash    string
	)

	for rows.Next() {
		err = rows.Scan(&reqStatus, &reqKind, &reqCID, &txKind, &txHash, &txStatus)
		if err != nil {
			return nil, fmt.Errorf("db.Scan error: %w", err)
		}

		result.Status = Status(reqStatus).String()

		exists := false

		for _, tx := range result.Txs {
			if tx.Hash == txHash {
				exists = true
				break
			}
		}

		if !exists {
			result.Txs = append(result.Txs, &TxResult{
				Kind:   TxKind(txKind).String(),
				Status: Status(txStatus).String(),
				Hash:   txHash,
			})
		}

		if reqCID == "" {
			continue
		}

		exists = false

		for _, doc := range result.Docs {
			if doc.CID == reqCID {
				exists = true
				break
			}
		}

		if !exists {
			result.Docs = append(result.Docs, &DocResult{
				Kind: RequestKind(reqKind).String(),
				CID:  reqCID,
			})
		}
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("GetRequestByID marshal error: %w", err)
	}

	return resultBytes, nil
}
