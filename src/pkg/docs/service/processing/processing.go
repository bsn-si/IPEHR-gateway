package processing

import (
	"context"
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

	RequestEhrCreate          RequestKind = 0
	RequestEhrGetBySubject    RequestKind = 1
	RequestEhrGetByID         RequestKind = 2
	RequestEhrStatusUpdate    RequestKind = 3
	RequestEhrStatusGetByID   RequestKind = 4
	RequestEhrStatusGetByTime RequestKind = 5
	RequestCompositionCreate  RequestKind = 6
	RequestCompositionUpdate  RequestKind = 7
	RequestCompositionGetByID RequestKind = 8
	RequestCompositionDelete  RequestKind = 9

	TxSetEhrUser        TxKind = 0
	TxSetEhrBySubject   TxKind = 1
	TxSetEhrDocs        TxKind = 2
	TxSetDocAccess      TxKind = 3
	TxDeleteDoc         TxKind = 4
	TxFilecoinStartDeal TxKind = 5
)

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
				p.execFilecoin()
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

		result := p.db.Where("kind IN ? AND statuses IN ?", txKinds, statuses).First(&tx)
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
			// do nothing
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

		for _, tx := range txs {
			if tx.Status != StatusSuccess {
				return false, nil
			}
		}
	}

	return true, nil
}

func (p *Proc) RequestStatus(reqID string) (Status, error) {
	var req Request
	result := p.db.Model(&req).Where("req_id = ?", reqID).First(&req)
	if result.Error != nil {
		return 0, result.Error
	}

	return req.Status, nil
}
