package processing

import (
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type (
	RequestKind uint8

	SuperRequest struct {
		dbTx    *gorm.DB
		request *Request
	}

	Request struct {
		gorm.Model
		ReqID   string `gorm:"index:idx_request,unique"`
		Kind    RequestKind
		Status  Status
		UserID  string
		EhrUUID string
	}

	RequestDataEtherium struct {
		gorm.Model
		RequestID   uint
		Request     Request
		BaseUIDHash string
		Version     string
	}

	RequestDataFileCoin struct {
		gorm.Model
		RequestID    uint
		Request      Request
		CID          string
		DealCID      string
		MinerAddress string
	}
)

const (
	RequestUnknown            RequestKind = 0
	RequestEhrCreate          RequestKind = 1
	RequestEhrGetBySubject    RequestKind = 2
	RequestEhrGetByID         RequestKind = 3
	RequestEhrStatusCreate    RequestKind = 4
	RequestEhrStatusUpdate    RequestKind = 5
	RequestEhrStatusGetByID   RequestKind = 6
	RequestEhrStatusGetByTime RequestKind = 7
	RequestCompositionCreate  RequestKind = 8
	RequestCompositionUpdate  RequestKind = 9
	RequestCompositionGetByID RequestKind = 10
	RequestCompositionDelete  RequestKind = 11
)

func (p *Proc) AddRequest(dbTx *gorm.DB, req *Request) (*SuperRequest, error) {
	if result := dbTx.Create(req); result.Error != nil {
		return nil, fmt.Errorf("db.Create error: %w", result.Error)
	}

	return &SuperRequest{dbTx: dbTx, request: req}, nil
}

func (s *SuperRequest) SetRequestKind(kind RequestKind) error {
	return s.dbTx.Model(&s.request).Update("kind", kind).Error
}

func (s *SuperRequest) UpdateFileCoinData(cid, deaclCid, minerAddress string) error {
	return s.dbTx.Model(&s.request).Updates(RequestDataFileCoin{
		CID:          cid,
		DealCID:      deaclCid,
		MinerAddress: minerAddress,
	}).Error
}

func (s *SuperRequest) AddEthData(baseUIDHash, version string) (*RequestDataEtherium, error) {
	var dataEtherium = RequestDataEtherium{
		BaseUIDHash: baseUIDHash,
		Version:     version,
		Request:     *s.request,
		RequestID:   s.request.ID,
	}

	db := s.dbTx.Create(&dataEtherium)

	return &dataEtherium, db.Error
}

func (s *SuperRequest) AddFileCoinData(cid, deaclCid, minerAddress string) (*RequestDataFileCoin, error) {
	var dataFileCoin = RequestDataFileCoin{
		CID:          cid,
		DealCID:      deaclCid,
		MinerAddress: minerAddress,
		Request:      *s.request,
		RequestID:    s.request.ID,
	}

	db := s.dbTx.Create(&dataFileCoin)

	return &dataFileCoin, db.Error
}

func (s *SuperRequest) AddTx(tx *Tx) error {
	db := s.dbTx.Create(&tx)

	return db.Error
}

func (s *SuperRequest) ReqID() string {
	return s.request.ReqID
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
	Status string       `json:"status"`
	Docs   []*DocResult `json:"docs"`
	Txs    []*TxResult  `json:"txs"`
}

type RequestsResult map[string]*RequestResult

func (p *Proc) requests(userID, reqId string, limit, offset int) (RequestsResult, error) {

	criteria := ""
	if reqId != "" {
		criteria = `AND requests.req_id = "@req_id"`
		criteria = strings.Replace(criteria, "@req_id", reqId, 1)
	}

	query := `SELECT
				r.req_id, r.status, r.kind,
				fc.c_id, fc.deal_c_id, fc.miner_address,
				t.kind, t.hash, coalesce(t_parent.hash, '') AS t_parent_hash, t.status
			FROM (select * from requests 
			               WHERE requests.user_id = @user_id @criteria
			               ORDER BY requests.id desc LIMIT @limit OFFSET @offset 
			     ) r
					 left join request_data_file_coins fc on fc.request_id = r.id
					 left join txes t on t.request_id = r.id
			         left join txes t_parent on t_parent.id = t.parent_tx_id
			ORDER BY r.id desc, t.id;`

	query = strings.Replace(query, "@criteria", criteria, 1)
	rows, err := p.db.Raw(query, map[string]interface{}{"user_id": userID, "limit": limit, "offset": offset, "criteria": criteria}).Rows()

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
		txKind       uint8
		txHash       string
		txParentHash string
		txStatus     uint8
	)

	result := make(RequestsResult)
	for rows.Next() {
		err = rows.Scan(&reqID, &reqStatus, &reqKind, &reqCID, &reqDealCID, &reqMinerAddr, &txKind, &txHash, &txParentHash, &txStatus)
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

		req.Txs = append(req.Txs, &TxResult{
			Kind:       TxKind(txKind).String(),
			Status:     Status(txStatus).String(),
			Hash:       txHash,
			ParentHash: txParentHash,
		})

		if reqCID != "" {
			req.Docs = append(req.Docs, &DocResult{
				Kind:         RequestKind(reqKind).String(),
				CID:          reqCID,
				DealCID:      reqDealCID,
				MinerAddress: reqMinerAddr,
			})
		}
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

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("GetRequestById marshal error: %w", err)
	}

	return resultBytes, nil
}
