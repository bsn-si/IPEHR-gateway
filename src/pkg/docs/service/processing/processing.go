package processing

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/ipfs/go-cid"
	"gorm.io/gorm"

	project_common "hms/gateway/pkg/common"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/storage/filecoin"
	"hms/gateway/pkg/storage/ipfs"
)

type (
	Status            uint8
	TxKind            uint8
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
		ParentTxID uint
		RequestID  uint
		Request    Request
		ServiceID  uint
		Service    BlockChainService
		Kind       TxKind
		Hash       string // TODO should be normalized for space saving
		Status     Status
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

	TxUnknown TxKind = iota
	TxMultiCall
	TxSetEhrUser
	TxSetEhrBySubject
	TxSetEhrDocs
	TxSetDocAccess
	TxDeleteDoc
	TxFilecoinStartDeal
	TxEhrCreateWithID
	TxUpdateEhrStatus
	TxAddEhrDoc
	TxSetDocKeyEncrypted

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
		TxMultiCall:          "MultiCall",
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

func (p *Proc) AddTx(request *SuperRequest, txHash string, kind TxKind, service BlockChainService, serviceID uint, parentTxID uint) (*Tx, error) {
	var dbTx = Tx{
		ParentTxID: parentTxID,
		Hash:       txHash,
		Kind:       kind,
		Status:     StatusPending,
		Service:    service,
		ServiceID:  serviceID,
		Request:    *request.request,
		RequestID:  request.request.ID,
	}

	if err := request.AddTx(&dbTx); err != nil {
		return nil, fmt.Errorf("db.Create transaction error: %w", err)
	}

	return &dbTx, nil
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
	tickerFilecoinStartDeal := time.NewTicker(5 * time.Second)
	tickerDealFinisher := time.NewTicker(5 * time.Second)
	tickerFilecoinRetrieve := time.NewTicker(1 * time.Minute)

	var mFinisher sync.Mutex

	go func() {
		logf("Started")

		for {
			select {
			case <-tickerBlockchain.C:
				if !p.lock {
					p.ExecBlockchain()
				}
			case <-tickerFilecoinStartDeal.C:
				p.execFilecoinStartDealStatus()
			case <-tickerFilecoinRetrieve.C:
				p.execFilecoinRetrieve()
			case <-tickerDealFinisher.C:
				p.execDealFinisher(&mFinisher)
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

func (p *Proc) ExecBlockchain() {
	p.lock = true
	defer func() { p.lock = false }()

	statuses := []Status{
		StatusPending,
		StatusProcessing,
	}

	var txs []Tx

	result := p.db.Model(&Tx{}).Select("txes.id, txes.hash, txes.status, txes.request_id").Joins("LEFT JOIN requests ON requests.id = txes.request_id").Where("txes.parent_tx_id = 0 AND txes.service = ?", BcEthereum).Where("requests.status IN ? AND requests.deleted_at IS NULL", statuses).Find(&txs)

	if result.Error != nil || result.RowsAffected == 0 {
		return
	}

	for _, tx := range txs {
		ctx, cancel := context.WithTimeout(context.Background(), project_common.BlockchainTxProcAwaitTime*time.Second)
		defer cancel()

		h := common.HexToHash(tx.Hash)
		txID := strconv.Itoa(int(tx.ID))

		var status Status

		receipt, err := p.ethClient.TransactionReceipt(ctx, h)
		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				status = StatusPending
			} else {
				log.Printf(txID, "TransactionReceipt error: %v txHash %s", err, h.String())
				continue
			}
		} else {
			status = Status(receipt.Status)
		}

		if status == tx.Status {
			continue
		}

		log.Println(txID, "Tx", tx.Hash, "status", status)

		if result = p.db.Model(Tx{}).Where("id = ? OR parent_tx_id = ?", tx.ID, tx.ID).Update("status", status); result.Error != nil {
			log.Println(txID, "db.Update error:", result.Error)
		}
	}
}

func (p *Proc) execFilecoinStartDealStatus() {
	logf("Filecoin StartDeal statuses started")
	defer logf("Filecoin StartDeal statuses finished")

	statuses := []Status{
		StatusPending,
		StatusProcessing,
	}

	txs := []Tx{}

	result := p.db.Model(&Tx{}).Select("txes.id, txes.hash, txes.status, txes.request_id").Joins("LEFT JOIN requests ON requests.id = txes.request_id").Where("txes.parent_tx_id = 0 AND txes.service = ?", BcFileCoin).Where("requests.status IN ? AND requests.deleted_at IS NULL", statuses).Find(&txs)

	if result.Error != nil {
		log.Println("DB get filecoin transactions error:", result.Error)
	}

	if result.RowsAffected == 0 {
		return
	}

	for _, tx := range txs {
		ctx, cancel := context.WithTimeout(context.Background(), project_common.FilecoinTxProcAwaitTime*time.Second)
		defer cancel()

		txID := strconv.Itoa(int(tx.ID))

		dealCID, err := cid.Parse(tx.Hash)
		if err != nil {
			log.Println(txID, "cid.Parse error:", err, "tx.Hash:", tx.Hash)
			break
		}

		dealStatus, err := p.filecoinClient.GetDealStatus(ctx, &dealCID)
		if err != nil {
			log.Println(txID, "filecoinClient.GetDealStatus error:", err, "dealCID", dealCID.String())
			break
		}

		var status Status

		switch dealStatus {
		case storagemarket.StorageDealActive:
			status = StatusSuccess
		case storagemarket.StorageDealError: // TODO need to do research
			fallthrough
		case storagemarket.StorageDealProposalRejected: // TODO need to do research
			fallthrough
		case storagemarket.StorageDealExpired: // TODO need to do research
			fallthrough
		case storagemarket.StorageDealSlashed: // TODO need to do research
			status = StatusFailed
		default:
			status = StatusProcessing
		}

		if status == tx.Status {
			continue
		}

		if result = p.db.Model(Tx{}).Where("id = ? OR parent_tx_id = ?", tx.ID, tx.ID).Update("status", status); result.Error != nil {
			log.Println(txID, "db.Update error:", result.Error)
			break
		}
	}
}

func (p *Proc) execDealFinisher(m *sync.Mutex) {
	m.Lock()

	logf("execDealFinisher started")

	defer func() {
		logf("execDealFinisher finished")
		m.Unlock()
	}()

	query := `select id
		from (SELECT requests.id,
					 (count(txes.status)) - (count(txes_succeeds.status)) as txs_left
			  FROM requests
					   LEFT JOIN txes ON requests.id = txes.request_id
					   LEFT JOIN txes txes_succeeds ON txes.id = txes_succeeds.id and txes_succeeds.status = @status_success
			  WHERE (requests.status IN (@status_pending, @status_processing) AND requests.deleted_at IS NULL)
				AND txes.deleted_at IS NULL
				AND txes.parent_tx_id = 0
			 )
		where id is not null and txs_left = 0`

	requests := []int{}
	result := p.db.Raw(query, map[string]interface{}{"status_success": StatusSuccess, "status_pending": StatusPending, "status_processing": StatusProcessing}).Find(&requests)

	if result.Error != nil {
		log.Println("DB get list of success transactions error:", result.Error)
	}

	if result.RowsAffected == 0 {
		return
	}

	p.db.Model(Request{}).Where("id IN ?", requests).Updates(Retrieve{Status: StatusSuccess})

	if p.db.Error != nil {
		log.Println("DB update error:", result.Error)
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

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("downloadFile ReadAll error: %w", err)
	}

	return data, nil
}
