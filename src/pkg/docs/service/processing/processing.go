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
	TxStatus    uint8
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
		ReqID        string `gorm:"primaryKey"`
		UserID       string
		EhrUUID      string
		Kind         RequestKind
		CID          string
		DealCID      string
		MinerAddress string
		BaseUIDHash  string
		Version      string
	}

	Tx struct {
		gorm.Model
		Time    time.Time
		ReqID   string `gorm:"index:idx_tx,unique"`
		Kind    TxKind `gorm:"index:idx_tx,unique"`
		Hash    string
		Status  TxStatus
		Comment string
	}
)

const (
	StatusFailed     TxStatus = 0
	StatusSuccess    TxStatus = 1
	StatusPending    TxStatus = 2
	StatusProcessing TxStatus = 3

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

func (p *Proc) AddTx(reqID, txHash, comment string, kind TxKind, status TxStatus) error {
	result := p.db.Create(&Tx{
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

	for {
		tx := Tx{}

		result := p.db.Where("kind IN ?", txKinds).First(&tx)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		h := common.HexToHash(tx.Hash)

		receipt, err := p.ethClient.TransactionReceipt(ctx, h)
		if err != nil {
			if errors.Is(err, ethereum.NotFound) {
				receipt.Status = uint64(StatusPending)
			} else {
				log.Printf("TransactionReceipt error: %v txHash %s", err, h.String())
				break
			}
		}

		if receipt.Status != uint64(tx.Status) {
			tx.Status = TxStatus(receipt.Status)

			if result = p.db.Save(&tx); result.Error != nil {
				log.Println("db.Save error:", result.Error)
				break
			}

			done, err := p.isRequestDone(tx.ReqID)
			if err != nil {
				log.Println("checkRequestDone error:", err, "reqID", tx.ReqID)
				break
			}

			if done {
				// delete transactions connected with reqID
				result = p.db.Where("req_id = ?", tx.ReqID).Delete(&Tx{})
				if result.Error != nil {
					log.Printf("db txs delete error: %v reqID %s", result.Error, tx.ReqID)
					break
				}

				// delete request with reqID
				result = p.db.Delete(&Request{}, tx.ReqID)
				if result.Error != nil {
					log.Printf("db request delete error: %v reqID %s", result.Error, tx.ReqID)
					break
				}
			}
		}
	}
}

func (p *Proc) execFilecoin() {
	txKinds := []TxKind{
		TxFilecoinStartDeal,
	}

	for {
		tx := Tx{}

		result := p.db.Where("kind IN ?", txKinds).First(&tx)
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

	result := p.db.First(&request, reqID)
	if result.Error != nil {
		return false, fmt.Errorf("db.First error: %w", result.Error)
	}

	var txsToCheck []TxKind

	switch request.Kind {
	case RequestEhrCreate:
		txsToCheck = []TxKind{
			TxSetEhrUser,
			TxSetEhrBySubject,
			TxSetEhrDocs,
			TxSetDocAccess,
			TxFilecoinStartDeal,
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
			TxFilecoinStartDeal,
		}
	case RequestCompositionUpdate:
		txsToCheck = []TxKind{
			TxSetEhrDocs,
			TxSetDocAccess,
			TxFilecoinStartDeal,
		}
	case RequestCompositionGetByID:
	case RequestCompositionDelete:
	default:
		return false, fmt.Errorf("%w Unknown request.Kind %d", errors.ErrCustom, request.Kind)
	}

	for _, txKind := range txsToCheck {
		tx := &Tx{}

		result := p.db.Model(tx).Where("req_id = ? AND kind = ?", reqID, txKind).First(tx)
		if result.Error != nil {
			return false, fmt.Errorf("db Find tx error: %w reqID %s txKind %d", result.Error, reqID, txKind)
		}

		if tx.Status != StatusSuccess {
			return false, nil
		}
	}

	return true, nil
}
