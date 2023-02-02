package processing

import (
	"encoding/json"
	"fmt"

	"gorm.io/gorm"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type RequestInterface interface {
	Commit() error
	AddEthereumTx(TxKind, string)
	AddFilecoinTx(TxKind, string, string, string)
}

type (
	RequestKind uint8

	Request struct {
		gorm.Model
		ReqID     string      `gorm:"index:idx_request,unique"`
		Kind      RequestKind `gorm:"kind" json:"-"`
		KindStr   string      `gorm:"-" json:"Kind"`
		Status    Status      `gorm:"status" json:"-"`
		StatusStr string      `gorm:"-" json:"Status"`
		UserID    string
		EhrUUID   string

		db     *gorm.DB      `gorm:"-"`
		ethTxs []*EthereumTx `gorm:"-"`
		fcTxs  []*FileCoinTx `gorm:"-"`
	}

	Tx struct {
		ReqID     string `gorm:"req_id" json:"-"`
		Kind      TxKind `gorm:"kind" json:"-"`
		KindStr   string `gorm:"-" json:"Kind"`
		Status    Status `gorm:"status" json:"-"`
		StatusStr string `gorm:"-" json:"Status"`
		Comment   string
	}

	EthereumTx struct {
		Tx
		Hash string
	}

	FileCoinTx struct {
		Tx
		CID          string
		DealCID      string
		MinerAddress string
		DealID       uint
	}
)

const (
	RequestUnknown RequestKind = iota
	RequestEhrCreate
	RequestEhrCreateWithID
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
	RequestUserRegister
	RequestDocAccessSet
	RequestQueryStore
	RequestUserGroupCreate
	RequestUserGroupAddUser
	RequestUserGroupRemoveUser
	RequestContributionCreate
	RequestTemplateCreate
	RequestDirectoryCreate
	RequestDirectoryUpdate
	RequestDirectoryDelete
)

func (p *Proc) NewRequest(reqID, userID, ehrUUID string, kind RequestKind) (*Request, error) {
	if userID == "" {
		return nil, fmt.Errorf("%w userID is empty", errors.ErrIncorrectRequest)
	}

	req := &Request{
		ReqID:   reqID,
		Kind:    kind,
		Status:  StatusProcessing,
		UserID:  userID,
		EhrUUID: ehrUUID,
		db:      p.db,
	}

	return req, nil
}

func (r *Request) Commit() error {
	dbTx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			dbTx.Rollback()
		}
	}()

	if result := dbTx.Create(r); result.Error != nil {
		return fmt.Errorf("NewRequest db.Create error: %w", result.Error)
	}

	for _, tx := range r.ethTxs {
		if result := dbTx.Create(tx); result.Error != nil {
			dbTx.Rollback()
			return fmt.Errorf("db.Create ethereum transaction error: %w", result.Error)
		}
	}

	for _, tx := range r.fcTxs {
		if result := dbTx.Create(tx); result.Error != nil {
			dbTx.Rollback()
			return fmt.Errorf("db.Create filecoin transaction error: %w", result.Error)
		}
	}

	if err := dbTx.Commit().Error; err != nil {
		dbTx.Rollback()
		return fmt.Errorf("Request commit error: %w reqID %s", err, r.ReqID)
	}

	return nil
}

func (r *Request) AddEthereumTx(kind TxKind, hash string) {
	tx := &EthereumTx{
		Tx: Tx{
			ReqID:  r.ReqID,
			Kind:   kind,
			Status: StatusPending,
		},
		Hash: hash,
	}

	r.ethTxs = append(r.ethTxs, tx)
}

func (r *Request) AddFilecoinTx(kind TxKind, CID, dealCID, minerAddress string) {
	tx := &FileCoinTx{
		Tx: Tx{
			ReqID:  r.ReqID,
			Kind:   kind,
			Status: StatusProcessing,
		},
		CID:          CID,
		DealCID:      dealCID,
		MinerAddress: minerAddress,
	}

	r.fcTxs = append(r.fcTxs, tx)
}

type TxResult struct {
	Kind       string `json:"kind"`
	Status     string `json:"status"`
	Hash       string `json:"hash"`
	ParentHash string `json:"parentHash"`
}

type DocResult struct {
	Kind         string `json:"kind"`
	CID          string `json:"cid"`
	DealCID      string `json:"dealCid"`
	MinerAddress string `json:"minerAddress"`
}

type RequestResult struct {
	Status   string        `json:"status"`
	Kind     string        `json:"kind"`
	Ethereum []*EthereumTx `json:"ethereum"`
	Filecoin []*FileCoinTx `json:"filecoin"`
}

type RequestsResult map[string]*RequestResult

type RequestsCriteria struct {
	kind   []RequestKind
	status []Status
}

func (rc *RequestsCriteria) ByStatus(status []Status) {
	rc.status = status
}

func (rc *RequestsCriteria) ByKind(kind []RequestKind) {
	rc.kind = kind
}

func (rc *RequestsCriteria) addInQuery(tx *gorm.DB) *gorm.DB {
	if len(rc.status) != 0 {
		tx.Where("status IN ?", rc.status)
	}

	if len(rc.kind) != 0 {
		tx.Where("kind IN ?", rc.kind)
	}

	return tx
}

func (p *Proc) requests(userID, reqID string, limit, offset int) (RequestsResult, error) {
	return p.requestsWithCriteria(userID, reqID, limit, offset, RequestsCriteria{})
}

func (p *Proc) requestsWithCriteria(userID, reqID string, limit, offset int, criteria RequestsCriteria) (RequestsResult, error) {
	var (
		requests    []*Request
		reqIDs      []string
		ethTxs      []*EthereumTx
		filecoinTxs []*FileCoinTx
		result      = make(RequestsResult)
	)

	// requests
	queryReq := p.db.Model(&Request{}).Where("user_id = ?", userID)

	if reqID != "" {
		queryReq = queryReq.Where("req_id = ?", reqID)
	}

	queryReq = criteria.addInQuery(queryReq)

	err := queryReq.Limit(limit).Offset(offset).Find(&requests).Error
	if err != nil {
		return nil, fmt.Errorf("Requests select error: %w userID: %s reqID: %s limit: %d offset: %d", err, userID, reqID, limit, offset)
	}

	if len(requests) == 0 {
		return nil, errors.ErrNotFound
	}

	for _, r := range requests {
		reqIDs = append(reqIDs, r.ReqID)
	}

	// ethereum transactions
	err = p.db.Model(&EthereumTx{}).Where("req_id IN ?", reqIDs).Find(&ethTxs).Error
	if err != nil {
		return nil, fmt.Errorf("Ethereum transactions select error: %w userID: %s reqID: %s limit: %d offset: %d", err, userID, reqID, limit, offset)
	}

	// filecoin transactions
	err = p.db.Model(&FileCoinTx{}).Where("req_id IN ?", reqIDs).Find(&filecoinTxs).Error
	if err != nil {
		return nil, fmt.Errorf("Ethereum transactions select error: %w userID: %s reqID: %s limit: %d offset: %d", err, userID, reqID, limit, offset)
	}

	for _, r := range requests {
		reqResult := &RequestResult{
			Status: r.Status.String(),
			Kind:   r.Kind.String(),
		}

		for _, ethTx := range ethTxs {
			if ethTx.ReqID == r.ReqID {
				ethTx.KindStr = ethTx.Kind.String()
				ethTx.StatusStr = ethTx.Status.String()
				reqResult.Ethereum = append(reqResult.Ethereum, ethTx)
			}
		}

		for _, fcTx := range filecoinTxs {
			if fcTx.ReqID == r.ReqID {
				fcTx.KindStr = fcTx.Kind.String()
				fcTx.StatusStr = fcTx.Status.String()
				reqResult.Filecoin = append(reqResult.Filecoin, fcTx)
			}
		}

		result[r.ReqID] = reqResult
	}

	return result, nil
}

func (p *Proc) GetRequests(userID string, limit, offset int) ([]byte, error) {
	result, err := p.requests(userID, "", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("GetRequests error: %w", err)
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("GetRequests marshal error: %w", err)
	}

	return resultBytes, nil
}

func (p *Proc) GetRequest(userID, reqID string) ([]byte, error) {
	result, err := p.requests(userID, reqID, 1, 0)
	if err != nil {
		return nil, fmt.Errorf("GetRequestById error: %w", err)
	}

	resultBytes, err := json.Marshal(result[reqID])
	if err != nil {
		return nil, fmt.Errorf("GetRequestById marshal error: %w", err)
	}

	return resultBytes, nil
}

func (p *Proc) GetRequestsByKindInProgress(userID string, kind RequestKind) (RequestsResult, error) {
	c := RequestsCriteria{}
	c.ByStatus([]Status{StatusPending, StatusProcessing})
	c.ByKind([]RequestKind{kind})

	return p.requestsWithCriteria(userID, "", 0, 0, c)
}
