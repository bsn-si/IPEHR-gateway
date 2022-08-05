package processing

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"

	"hms/gateway/pkg/errors"
)

type (
	TxStatus    uint8
	TxKind      uint8
	RequestKind uint8

	Proc struct {
		db        *gorm.DB
		ethClient *ethclient.Client
		lock      bool
		done      chan bool
	}

	Request struct {
		ReqID       string `gorm:"primaryKey"`
		UserID      string
		EhrUUID     string
		Kind        RequestKind
		CID         string
		BaseUIDHash string
		Version     string
	}

	BlockchainTx struct {
		gorm.Model
		Time    time.Time
		ReqID   string `gorm:"index:idx_bcTx,unique"`
		Kind    TxKind `gorm:"index:idx_bcTx,unique"`
		Hash    string
		Status  TxStatus
		Comment string
	}
)

const (
	StatusFailed  TxStatus = 0
	StatusSuccess TxStatus = 1
	StatusPending TxStatus = 2

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

	TxSetEhrUser      TxKind = 0
	TxSetEhrBySubject TxKind = 1
	TxSetEhrDocs      TxKind = 2
	TxSetDocAccess    TxKind = 3
	TxDeleteDoc       TxKind = 4
)

func New(db *gorm.DB, ethClient *ethclient.Client) *Proc {
	return &Proc{
		db:        db,
		ethClient: ethClient,
		done:      make(chan bool),
	}
}

func (p *Proc) AddRequest(req *Request) error {
	if result := p.db.Create(req); result.Error != nil {
		return fmt.Errorf("db.Create error: %w", result.Error)
	}

	return nil
}

func (p *Proc) AddBlockchainTx(reqID, txHash, comment string, kind TxKind, status TxStatus) error {
	result := p.db.Create(&BlockchainTx{
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
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				if !p.lock {
					p.exec()
				}
			case <-p.done:
				return
			}
		}
	}()
}

func (p *Proc) Stop() {
	p.done <- true
}

func (p *Proc) exec() {
	p.lock = true
	defer func() { p.lock = false }()

	for {
		tx := BlockchainTx{}

		result := p.db.First(&tx)
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
				result = p.db.Where("req_id = ?", tx.ReqID).Delete(&BlockchainTx{})
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
		}
	case RequestEhrGetBySubject:
	case RequestEhrGetByID:
	case RequestEhrStatusUpdate:
	case RequestEhrStatusGetByID:
	case RequestEhrStatusGetByTime:
	case RequestCompositionCreate:
	case RequestCompositionUpdate:
	case RequestCompositionGetByID:
	case RequestCompositionDelete:
	default:
		return false, fmt.Errorf("%w Unknown request.Kind %d", errors.ErrCustom, request.Kind)
	}

	for _, txKind := range txsToCheck {
		tx := &BlockchainTx{}

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
