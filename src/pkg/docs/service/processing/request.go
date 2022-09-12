package processing

import (
	"fmt"

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
		ReqID string `gorm:"index:idx_request,unique"`
		//BcService BlockChainService
		Kind RequestKind
		//Kind    RequestKind `gorm:"index:idx_request,unique"` // TODO move to another table?
		Status  Status
		UserID  string
		EhrUUID string
		RequestDataFileCoin
		RequestDataEtherium
	}

	RequestDataEtherium struct {
		gorm.Model
		BaseUIDHash string
		Version     string
	}

	RequestDataFileCoin struct {
		gorm.Model
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

func (s *SuperRequest) UpdateEthData(baseUIDHash, version string) error {
	return s.dbTx.Model(&s.request).Updates(RequestDataEtherium{
		BaseUIDHash: baseUIDHash,
		Version:     version,
	}).Error
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
	}

	db := s.dbTx.Model(&s.request).Create(&dataEtherium)

	return &dataEtherium, db.Error
}

func (s *SuperRequest) AddFileCoinData(cid, deaclCid, minerAddress string) (*RequestDataFileCoin, error) {
	var dataFileCoin = RequestDataFileCoin{
		CID:          cid,
		DealCID:      deaclCid,
		MinerAddress: minerAddress,
	}

	db := s.dbTx.Model(&s.request).Create(&dataFileCoin)

	return &dataFileCoin, db.Error
}

func (s *SuperRequest) AddTx(tx *Tx) error {
	db := s.dbTx.Model(&s.request).Create(&tx)

	return db.Error
}

func (s *SuperRequest) ReqID() string {
	return s.request.ReqID
}
