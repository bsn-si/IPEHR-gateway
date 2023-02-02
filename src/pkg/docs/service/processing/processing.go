package processing

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/filecoin-project/go-fil-markets/retrievalmarket"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/ipfs/go-cid"
	"gorm.io/gorm"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage/filecoin"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage/ipfs"
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
		lockEthereum     bool
		lockFilecoin     bool
		localStoragePath string
		done             chan bool
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
	TxSetDocGroupAccess
	TxDeleteDoc
	TxFilecoinStartDeal
	TxEhrCreateWithID
	TxUpdateEhrStatus
	TxAddEhrDoc
	TxSetDocKeyEncrypted
	TxSaveEhr
	TxSaveEhrStatus
	TxSaveComposition
	TxSaveTemplate
	TxUserRegister
	TxUserNew
	TxUserGroupCreate
	TxUserGroupAddUser
	TxUserGroupRemoveUser
	TxDocGroupCreate
	TxDocGroupAddDoc
	TxIndexDataUpdate
	TxCreateDirectory
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
		TxMultiCall:           "MultiCall",
		TxSetEhrUser:          "SetEhrUser",
		TxSetEhrBySubject:     "SetEhrBySubject",
		TxSetEhrDocs:          "SetEhrDocs",
		TxSetDocAccess:        "SetDocAccess",
		TxSetDocGroupAccess:   "SetDocGroupAccess",
		TxDeleteDoc:           "DeleteDoc",
		TxFilecoinStartDeal:   "FilecoinStartDeal",
		TxEhrCreateWithID:     "EhrCreateWithID",
		TxUpdateEhrStatus:     "UpdateEhrStatus",
		TxAddEhrDoc:           "AddEhrDoc",
		TxSetDocKeyEncrypted:  "SetDocKeyEncrypted",
		TxSaveEhr:             "SaveEhr",
		TxSaveEhrStatus:       "SaveEhrStatus",
		TxSaveComposition:     "SaveComposition",
		TxSaveTemplate:        "SaveTemplate",
		TxCreateDirectory:     "CreateDirectory",
		TxUserRegister:        "UserRegister",
		TxUserNew:             "UserNew",
		TxUserGroupCreate:     "UserGroupCreate",
		TxUserGroupAddUser:    "UserGroupAddUser",
		TxUserGroupRemoveUser: "UserGroupRemoveUser",
		TxDocGroupCreate:      "DocGroupCreate",
		TxDocGroupAddDoc:      "DocGroupAddDoc",
		TxIndexDataUpdate:     "IndexDataUpdate",
		TxUnknown:             "Unknown",
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
		RequestTemplateCreate:     "TemplateStore",
		RequestDirectoryCreate:    "DirectoryCreate",
		RequestDirectoryUpdate:    "DirectoryUpdate",
		RequestDirectoryDelete:    "DirectoryDelete",
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
	tickerEthereum := time.NewTicker(5 * time.Second)
	tickerFilecoin := time.NewTicker(5 * time.Minute)
	tickerDealFinisher := time.NewTicker(1 * time.Minute)
	tickerFilecoinRetrieve := time.NewTicker(1 * time.Minute)

	go func() {
		logf("Started")

		for {
			select {
			case <-tickerEthereum.C:
				if !p.lockEthereum {
					p.execEthereum()
				}
			case <-tickerFilecoin.C:
				if !p.lockFilecoin {
					p.execFilecoin()
				}
			case <-tickerFilecoinRetrieve.C:
				p.execFilecoinRetrieve()
			case <-tickerDealFinisher.C:
				if !p.lockEthereum && !p.lockFilecoin {
					p.execDealFinisher()
				}
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

func (p *Proc) execEthereum() {
	p.lockEthereum = true
	//logf("Ethereum transactions started")

	defer func() {
		//logf("Ethereum transactions finished")
		p.lockEthereum = false
	}()

	statuses := []Status{
		StatusPending,
		StatusProcessing,
	}

	var txs []struct {
		EthereumTx
	}

	result := p.db.Model(&EthereumTx{}).
		Select("ethereum_txes.req_id, ethereum_txes.hash, ethereum_txes.status").
		Where("ethereum_txes.status IN ?", statuses).
		Group("ethereum_txes.hash").
		Find(&txs)
	if result.Error != nil {
		logf("execEthereum get transactions error: %v", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		return
	}

	for _, tx := range txs {
		ctx, cancel := context.WithTimeout(context.Background(), common.BlockchainTxProcAwaitTime)

		h := eth_common.HexToHash(tx.Hash)

		var status Status

		receipt, err := p.ethClient.TransactionReceipt(ctx, h)
		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				status = StatusPending
			} else {
				cancel()
				logf("TransactionReceipt error: %v txHash %s", err, h.String())
				return
			}
		} else {
			status = Status(receipt.Status)
		}

		cancel()

		logf("Tx: %s status: %s", tx.Hash, status)

		if status == tx.Status {
			continue
		}

		if result = p.db.Model(&EthereumTx{}).Where("hash = ?", tx.Hash).Update("status", status); result.Error != nil {
			logf("db.Update error: %v", result.Error)
		}
	}
}

func (p *Proc) execFilecoin() {
	p.lockFilecoin = true

	logf("Filecoin deals started")

	defer func() {
		p.lockFilecoin = false

		logf("Filecoin deals finished")
	}()

	statuses := []Status{
		StatusPending,
		StatusProcessing,
	}

	txs := []FileCoinTx{}

	result := p.db.Model(&FileCoinTx{}).
		Select("req_id, c_id, deal_c_id, miner_address, status").
		Where("status IN ?", statuses).
		Find(&txs)
	if result.Error != nil {
		logf("DB get filecoin transactions error: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return
	}

	for _, tx := range txs {
		dealCID, err := cid.Parse(tx.DealCID)
		if err != nil {
			logf("cid.Parse error: %v dealCID: %s", err, tx.DealCID)
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), common.FilecoinTxProcAwaitTime)

		dealStatus, dealID, err := p.filecoinClient.GetDealStatus(ctx, &dealCID)
		if err != nil {
			logf("filecoinClient.GetDealStatus error: %v dealCID: %s", err, tx.DealCID)
			cancel()

			continue
		}

		cancel()

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

		err = p.db.Model(&FileCoinTx{}).Where("deal_c_id", tx.DealCID).Updates(map[string]interface{}{"status": status, "deal_id": dealID}).Error
		if err != nil {
			logf("db.Update error: %v", err)
		}
	}
}

func (p *Proc) execDealFinisher() {
	p.lockFilecoin = true
	p.lockEthereum = true

	logf("DealFinisher started")

	defer func() {
		p.lockEthereum = false
		p.lockFilecoin = false

		logf("DealFinisher finished")
	}()

	query := `UPDATE requests SET status = @success 
				WHERE req_id NOT IN (
					SELECT req_id 
						FROM ethereum_txes 
						WHERE status IN (@failed,@pending,@processing) 
						GROUP BY req_id
					UNION
					SELECT req_id 
						FROM file_coin_txes
						WHERE status IN (@failed,@pending,@processing)
						GROUP BY req_id
				)`

	if err := p.db.Exec(query, map[string]interface{}{
		"failed":     StatusFailed,
		"success":    StatusSuccess,
		"pending":    StatusPending,
		"processing": StatusProcessing,
	}).Error; err != nil {
		logf("DB get list of success transactions error: %v", err)
	}
}

func (p *Proc) execFilecoinRetrieve() {
	p.lockFilecoin = true

	logf("Filecoin retrieve started")

	defer func() {
		p.lockFilecoin = false

		logf("Filecoin retrieve finished")
	}()

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
