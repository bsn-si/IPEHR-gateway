package processing

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/ipfs/go-cid"
	"gorm.io/gorm"

	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/storage/filecoin"
	"hms/gateway/pkg/storage/ipfs"
)

type (
	Status            uint8
	TxKind            uint8 // TODO useless type because we already have RequestKind, require refactoring
	BlockChainService uint8

	Proc struct {
		db               *gorm.DB
		ethClient        *ethclient.Client
		filecoinClient   *filecoin.Client
		ipfsClient       *ipfs.Client
		httpClient       *http.Client
		lock             bool
		localStoragePath string
		done             chan bool
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

	Retrieve struct {
		CID     string `gorm:"primaryKey"`
		DealID  retrievalmarket.DealID
		Status  Status
		Comment string
	}
)

const (
	StatusFailed     Status = 0
	StatusSuccess    Status = 1
	StatusPending    Status = 2
	StatusProcessing Status = 3
	StatusUnknown    Status = 255

	TxSetEhrUser         TxKind = 0
	TxSetEhrBySubject    TxKind = 1
	TxSetEhrDocs         TxKind = 2
	TxSetDocAccess       TxKind = 3
	TxDeleteDoc          TxKind = 4
	TxFilecoinStartDeal  TxKind = 5
	TxEhrCreateWithID    TxKind = 6
	TxUpdateEhrStatus    TxKind = 7
	TxAddEhrDoc          TxKind = 8
	TxSetDocKeyEncrypted TxKind = 9
	TxUnknown            TxKind = 255

	BcEthereum BlockChainService = 0
	BcFileCoin BlockChainService = 1
)

var (
	statuses = map[Status]string{
		StatusFailed:     "Failed",
		StatusSuccess:    "Success",
		StatusPending:    "Pending",
		StatusProcessing: "Processing",
		StatusUnknown:    "Unknown",
	}

	txKinds = map[TxKind]string{
		TxSetEhrUser:         "SetEhrUser",
		TxSetEhrBySubject:    "SetEhrBySubject",
		TxSetEhrDocs:         "SetEhrDocs",
		TxSetDocAccess:       "SetDocAccess",
		TxDeleteDoc:          "DeleteDoc",
		TxFilecoinStartDeal:  "FilecoinStartDeal",
		TxEhrCreateWithID:    "EhrCreateWithID",
		TxUpdateEhrStatus:    "UpdateEhrStatus",
		TxAddEhrDoc:          "AddEhrDoc",
		TxSetDocKeyEncrypted: "SetDocKeyEncrypted",

		TxUnknown: "Unknown",
	}

	reqKinds = map[RequestKind]string{
		RequestEhrCreate:          "EhrCreate",
		RequestEhrGetBySubject:    "EhrGetBySubject",
		RequestEhrGetByID:         "EhrGetByID",
		RequestEhrStatusCreate:    "EhrStatusCreate",
		RequestEhrStatusUpdate:    "EhrStatusUpdate",
		RequestEhrStatusGetByID:   "EhrStatusGetByID",
		RequestEhrStatusGetByTime: "EhrStatusGetByTime",
		RequestCompositionCreate:  "CompositionCreate",
		RequestCompositionUpdate:  "CompositionUpdate",
		RequestCompositionGetByID: "CompositionGetByID",
		RequestCompositionDelete:  "CompositionDelete",
	}
)

func (s Status) String() string {
	if status, ok := statuses[s]; ok {
		return status
	}

	return statuses[StatusUnknown]
}

func (k TxKind) String() string {
	if tk, ok := txKinds[k]; ok {
		return tk
	}

	return txKinds[TxUnknown]
}

func (k RequestKind) String() string {
	if rk, ok := reqKinds[k]; ok {
		return rk
	}

	return reqKinds[RequestUnknown]
}

func logf(format string, a ...interface{}) {
	fmt.Printf("[PROC] %19s | %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		fmt.Sprintf(format, a...),
	)
}

func New(db *gorm.DB, ethClient *ethclient.Client, filecoinClient *filecoin.Client, ipfsClient *ipfs.Client, storagePath string) *Proc {
	return &Proc{
		db:               db,
		ethClient:        ethClient,
		filecoinClient:   filecoinClient,
		ipfsClient:       ipfsClient,
		httpClient:       http.DefaultClient,
		done:             make(chan bool),
		localStoragePath: storagePath,
	}
}

func (p *Proc) AddTx(reqID, txHash, comment string, kind TxKind) error {
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
		Status:  StatusPending,
		Comment: comment,
	})
	if result.Error != nil {
		return fmt.Errorf("db.Create error: %w", result.Error)
	}

	return nil
}

func (p *Proc) AddRetrieve(CID string) error {
	result := p.db.Create(&Retrieve{
		CID:     CID,
		Status:  StatusPending,
		Comment: "",
	})
	if result.Error != nil {
		return fmt.Errorf("db.Create error: %w", result.Error)
	}

	return nil
}

func (p *Proc) GetRetrieveStatus(CID *cid.Cid) (Status, error) {
	var ret Retrieve

	result := p.db.Model(&ret).Find(&ret, "c_id = ?", CID.String())
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return StatusUnknown, nil
	} else if result.RowsAffected == 0 {
		return StatusUnknown, nil
	} else if result.Error != nil {
		return 0, fmt.Errorf("Retrieve get error: %w CID: %s", result.Error, CID)
	}

	return ret.Status, nil
}

func (p *Proc) Start() {
	tickerBlockchain := time.NewTicker(5 * time.Second)
	tickerFilecoinStartDeal := time.NewTicker(5 * time.Minute)
	tickerFilecoinRetrieve := time.NewTicker(1 * time.Minute)

	go func() {
		logf("Started")

		for {
			select {
			case <-tickerBlockchain.C:
				if !p.lock {
					p.execBlockchain()
				}
			case <-tickerFilecoinStartDeal.C:
				p.execFilecoinStartDealStatus()
			case <-tickerFilecoinRetrieve.C:
				p.execFilecoinRetrieve()
			case <-p.done:
				logf("Stopped")
				return
			}
		}
	}()
}

func (p *Proc) Stop() {
	p.done <- true
}

func (p *Proc) BeginDbTx() *gorm.DB {
	return p.db.Begin()
}

func (p *Proc) CommitDbTx(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (p *Proc) RollbackDbTx(tx *gorm.DB) {
	tx.Rollback()
}

func (p *Proc) execBlockchain() {
	p.lock = true
	defer func() { p.lock = false }()

	// TODO we should separate different blockchain requests in different tables
	txKinds := []TxKind{
		TxSetEhrUser,
		TxSetEhrBySubject,
		TxSetEhrDocs,
		TxSetDocAccess,
		TxDeleteDoc,
		TxEhrCreateWithID,
		TxUpdateEhrStatus,
		TxAddEhrDoc,
		TxSetDocKeyEncrypted,
	}

	statuses := []Status{
		StatusPending,
		StatusProcessing,
	}

	for {
		tx := Tx{}

		result := p.db.Model(&tx).Select("req_id, hash, status").Where("kind IN ? AND status IN ?", txKinds, statuses).Group("hash").Find(&tx).Limit(1)
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
				log.Printf(tx.ReqID, "TransactionReceipt error: %v txHash %s", err, h.String())
				break
			}
		} else {
			status = Status(receipt.Status)
		}

		log.Println(tx.ReqID, "Tx", tx.Hash, "status", status)

		if status != tx.Status {
			if result = p.db.Model(Tx{}).Where("hash = ?", tx.Hash).Update("status", status); result.Error != nil {
				log.Println(tx.ReqID, "db.Update error:", result.Error)
				break
			}

			reqStatus, err := p.checkRequestStatus(tx.ReqID)
			if err != nil {
				log.Println(tx.ReqID, "checkRequestDone error:", err, "reqID", tx.ReqID)
				break
			}

			if reqStatus == StatusSuccess || reqStatus == StatusFailed {
				// update request status
				result = p.db.Exec("UPDATE requests SET status = ? WHERE req_id = ?", reqStatus, tx.ReqID)
				if result.Error != nil {
					break
				}
			}

			continue
		}

		// tx pending yet
		break
	}
}

// nolint
func (p *Proc) execFilecoinStartDealStatus() {
	logf("Filecoin StartDeal statuses started")
	defer logf("Filecoin StartDeal statuses finished")

	txKinds := []TxKind{
		TxFilecoinStartDeal,
	}

	statuses := []Status{
		StatusPending,
		StatusProcessing,
	}

	txs := []Tx{}

	for _, tx := range txs {
		result := p.db.Where("kind IN ? AND status IN ?", txKinds, statuses).Find(&txs)
		if result.Error != nil {
			if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				log.Println("DB get filecoin transactions error:", result.Error)
			}
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		dealCID, err := cid.Parse(tx.Hash)
		if err != nil {
			log.Println(tx.ReqID, "cid.Parse error:", err, "tx.Hash:", tx.Hash)
			break
		}

		dealStatus, err := p.filecoinClient.GetDealStatus(ctx, &dealCID)
		if err != nil {
			log.Println(tx.ReqID, "filecoinClient.GetDealStatus error:", err, "dealCID", dealCID.String())
			break
		}

		switch tx.Status {
		case StatusPending, StatusProcessing:
			if dealStatus == storagemarket.StorageDealActive {
				result = p.db.Exec("UPDATE requests SET status = ? WHERE req_id = ?", StatusSuccess, tx.ReqID)
				if result.Error != nil {
					log.Println(tx.ReqID, "db.Save error:", result.Error)
					return
				}
			}
		case StatusSuccess:
		}
	}
}

func (p *Proc) execFilecoinRetrieve() {
	logf("Filecoin retrieve started")
	defer logf("Filecoin retrieve finished")

	var rets []Retrieve

	statuses := []Status{
		StatusPending,
		StatusProcessing,
	}

	result := p.db.Model(&Retrieve{}).Find(&rets, "status IN ?", statuses)
	if result.Error != nil && result.RowsAffected == 0 {
		return
	} else if result.Error != nil {
		logf("Filecoin retrieve error: %v", result.Error)
	}

	for i, ret := range rets {
		var (
			nextStatus  Status
			comment     string
			ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
		)

		switch ret.Status {
		case StatusPending:
			nextStatus = StatusProcessing

			CID, err := cid.Decode(ret.CID)
			if err != nil {
				comment = fmt.Sprintf("Filecoin retrieve CID parse error: %v", err)
				logf(comment)

				nextStatus = StatusFailed
			}

			dealID, err := p.filecoinClient.StartRetrieve(ctx, &CID)
			if err != nil {
				comment = fmt.Sprintf("Filecoin retrieve StartRetrieve error: %v CID: %s", err, CID)
				logf(comment)

				nextStatus = StatusFailed
			}

			if err = p.db.Model(&rets[i]).Updates(Retrieve{
				Status:  nextStatus,
				Comment: comment,
				DealID:  dealID,
			}).Error; err != nil {
				logf("Filecoin retrieve DB update error: %v CID: %s", err, ret.CID)
			}
		case StatusProcessing:
			nextStatus = StatusSuccess

			dealStatus, err := p.filecoinClient.GetRetrieveStatus(ctx, ret.DealID)
			switch {
			case err != nil && !errors.Is(err, errors.ErrNotFound):
				logf("Filecoin retrieve error: %v, dealID: %d", err, ret.DealID)

				cancel()

				continue
			case err != nil && errors.Is(err, errors.ErrNotFound): // OK
			case dealStatus == retrievalmarket.DealStatusCompleted: // OK
			default:
				cancel()
				continue
			}

			CID, err := cid.Decode(ret.CID)
			if err != nil {
				comment = fmt.Sprintf("Filecoin retrieve CID parse error: %v CID: %s", err, ret.CID)
				logf(comment)
			}

			// Save retrieved file on lotus client host
			err = p.filecoinClient.SaveFile(ctx, &CID, ret.DealID)
			if err != nil {
				comment = fmt.Sprintf("Filecoin retrieve SaveFile error: %v", err)
				logf(comment)

				nextStatus = StatusFailed
			}

			// Download file
			file, err := p.downloadFile(&CID)
			if err != nil {
				comment = fmt.Sprintf("Filecoin retrieve download file error: %v", err)
				logf(comment)

				nextStatus = StatusFailed
			}

			// Add to IPFS
			_, err = p.ipfsClient.Add(ctx, file)
			if err != nil {
				comment = fmt.Sprintf("IpfsClient.Add error: %v", err)
				logf(comment)

				nextStatus = StatusFailed
			}

			if err = p.db.Model(&rets[i]).Updates(Retrieve{
				Status:  nextStatus,
				Comment: comment,
			}).Error; err != nil {
				logf("Filecoin retrieve DB update error: %v CID: %s", err, ret.CID)
			}

			logf("Filecoin file recovery complete CID %s", ret.CID)
		}

		cancel()
	}
}

func (p *Proc) downloadFile(CID *cid.Cid) ([]byte, error) {
	url := p.filecoinClient.BaseURL() + "/files/" + CID.String()

	resp, err := p.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request error: %w URL: %s", err, url)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w Download file response status error: %s", errors.ErrCustom, resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("downloadFile ReadAll error: %w", err)
	}

	return data, nil
}

func (p *Proc) checkRequestStatus(reqID string) (Status, error) {
	var txs []Tx
	result := p.db.Model(&Tx{}).Where("req_id = ? AND kind != ?", reqID, TxFilecoinStartDeal).Find(&txs)

	if result.Error != nil {
		return 0, fmt.Errorf("db Find tx error: %w", result.Error)
	}

	for _, tx := range txs {
		if tx.Status != StatusSuccess {
			return tx.Status, nil
		}
	}

	return StatusSuccess, nil
}

type TxResult struct {
	Kind   string `json:"kind"`
	Status string `json:"status"`
	Hash   string `json:"hash"`
}

type DocResult struct {
	Kind         string `json:"kind"`
	CID          string `json:"cid"`
	DealCID      string `json:"dealCid"`
	MinerAddress string `json:"minerAddress"`
}

type RequestResult struct {
	Status string       `json:"status"`
	Docs   []*DocResult `json:"docs"`
	Txs    []*TxResult  `json:"txs"`
}

type RequestsResult map[string]*RequestResult

func (p *Proc) GetRequests(userID string, limit, offset int) ([]byte, error) {
	result := make(map[string]*RequestResult)

	query := `SELECT r.req_id, r.status, r.kind, r.c_id, r.deal_c_id, r.miner_address, t.kind, t.hash, t.status
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
		reqID        string
		reqStatus    uint8
		reqKind      uint8
		reqCID       string
		reqDealCID   string
		reqMinerAddr string
		txStatus     uint8
		txKind       uint8
		txHash       string
	)

	for rows.Next() {
		err = rows.Scan(&reqID, &reqStatus, &reqKind, &reqCID, &reqDealCID, &reqMinerAddr, &txKind, &txHash, &txStatus)
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
				Kind:         RequestKind(reqKind).String(),
				CID:          reqCID,
				DealCID:      reqDealCID,
				MinerAddress: reqMinerAddr,
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

	query := `SELECT r.status, r.kind, r.c_id, r.deal_c_id, r.miner_address, t.kind, t.hash, t.status
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
		reqStatus    uint8
		reqKind      uint8
		reqCID       string
		reqDealCID   string
		reqMinerAddr string
		txStatus     uint8
		txKind       uint8
		txHash       string
	)

	for rows.Next() {
		err = rows.Scan(&reqStatus, &reqKind, &reqCID, &reqDealCID, &reqMinerAddr, &txKind, &txHash, &txStatus)
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
				Kind:         RequestKind(reqKind).String(),
				CID:          reqCID,
				DealCID:      reqDealCID,
				MinerAddress: reqMinerAddr,
			})
		}
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("GetRequestByID marshal error: %w", err)
	}

	return resultBytes, nil
}
