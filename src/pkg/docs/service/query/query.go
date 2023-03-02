package query

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"

	aqlprocessor "github.com/bsn-si/IPEHR-gateway/src/pkg/aql/processor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/compressor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/status"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/types"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/ehrIndexer"
)

const defaultVersion = "1.0.1"

type QueryExecuter interface { //nolint
	//Validate(q string) error
	ExecQuery(ctx context.Context, query *model.QueryRequest) (*model.QueryResponse, error)
}

type Service struct {
	*service.DefaultDocumentService

	qExec QueryExecuter
}

func NewService(docService *service.DefaultDocumentService, qExec QueryExecuter) *Service {
	return &Service{
		DefaultDocumentService: docService,
		qExec:                  qExec,
	}
}

func (s *Service) List(ctx context.Context, userID, systemID, qualifiedQueryName string) ([]*model.StoredQuery, error) {
	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	list, err := s.Infra.Index.ListDocByType(ctx, userID, systemID, types.Query)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("ListDocByType error: %w", err)
	}

	var result []*model.StoredQuery

	for i, dm := range list {
		dm := dm

		key, err := s.KeyFromAttribures(&dm, userPubKey, userPrivKey)
		if err != nil {
			return nil, fmt.Errorf("index %d KeyFromAttribures error: %w", i, err)
		}

		content, err := s.ContentFromAttributes(&dm, key)
		if err != nil {
			return nil, fmt.Errorf("index %d ContentFromAttributes error: %w", i, err)
		}

		var storedQuery model.StoredQuery

		err = msgpack.Unmarshal(content, &storedQuery)
		if err != nil {
			return nil, fmt.Errorf("index %d StoredQuery content unmarshal error: %w", i, err)
		}

		if qualifiedQueryName != "" && storedQuery.Name != qualifiedQueryName {
			continue
		}

		result = append(result, &storedQuery)
	}

	return result, nil
}

func (s *Service) GetByVersion(ctx context.Context, userID, systemID, name string, version *base.VersionTreeID) (*model.StoredQuery, error) {
	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	ehrUUID, err := s.Infra.Index.GetEhrUUIDByUserID(ctx, userID, systemID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("Index.GetEhrIDByUserId error: %w", err)
	}

	id := []byte(userID + systemID + name + version.String())
	idHash := sha3.Sum256(id)

	var vID [32]byte

	copy(vID[:], version.String())

	docMeta, err := s.Infra.Index.GetDocByVersion(ctx, ehrUUID, types.Query, &idHash, &vID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("Index.GetDocByVersion error: %w", err)
	}

	key, err := s.KeyFromAttribures(docMeta, userPubKey, userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("KeyFromAttribures error: %w", err)
	}

	content, err := s.ContentFromAttributes(docMeta, key)
	if err != nil {
		return nil, fmt.Errorf("ContentFromAttributes error: %w", err)
	}

	var storedQuery model.StoredQuery

	err = msgpack.Unmarshal(content, &storedQuery)
	if err != nil {
		return nil, fmt.Errorf("StoredQuery content unmarshal error: %w", err)
	}

	return &storedQuery, nil
}

func (s *Service) Validate(data []byte) bool {
	_, err := aqlprocessor.NewAqlProcessor(string(data)).Process()

	return err == nil
}

func (s *Service) Store(ctx context.Context, userID, systemID, reqID, qType, name, q string) (*model.StoredQuery, error) {
	v, _ := base.NewVersionTreeID(defaultVersion)

	return s.StoreVersion(ctx, userID, systemID, reqID, qType, name, v, q)
}

func (s *Service) StoreVersion(ctx context.Context, userID, systemID, reqID, qType, name string, version *base.VersionTreeID, q string) (*model.StoredQuery, error) {
	timestamp := time.Now()

	storedQuery := &model.StoredQuery{
		Name:        name,
		Type:        qType,
		Version:     version.String(),
		TimeCreated: timestamp.Format(common.OpenEhrTimeFormat),
		Query:       q,
	}

	id := []byte(userID + systemID + storedQuery.Name + storedQuery.Version)
	idHash := sha3.Sum256(id)
	key := chachaPoly.GenerateKey()

	content, err := msgpack.Marshal(storedQuery)
	if err != nil {
		return nil, fmt.Errorf("msgpack.Marshal error: %w", err)
	}

	contentCompresed, err := compressor.New(compressor.BestCompression).Compress(content)
	if err != nil {
		return nil, fmt.Errorf("Query Compress error: %w", err)
	}

	contentEncr, err := key.Encrypt(contentCompresed)
	if err != nil {
		return nil, fmt.Errorf("key.Encrypt content error: %w", err)
	}

	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	keyEncr, err := keybox.SealAnonymous(key.Bytes(), userPubKey)
	if err != nil {
		return nil, fmt.Errorf("keybox.SealAnonymous error: %w", err)
	}

	docMeta := &model.DocumentMeta{
		Status:    uint8(status.ACTIVE),
		Id:        idHash[:],
		Version:   []byte(version.String()),
		Timestamp: uint32(timestamp.Unix()),
		IsLast:    true,
		Attrs: []ehrIndexer.AttributesAttribute{
			{Code: model.AttributeKeyEncr, Value: keyEncr},         // encrypted with key
			{Code: model.AttributeContentEncr, Value: contentEncr}, // encrypted with userPubKey
			{Code: model.AttributeDocUIDHash, Value: idHash[:]},
		},
	}

	packed, err := s.Infra.Index.AddEhrDoc(ctx, types.Query, docMeta, userPrivKey, nil)
	if err != nil {
		return nil, fmt.Errorf("Index.AddEhrDoc error: %w", err)
	}

	procRequest, err := s.Proc.NewRequest(reqID, userID, "", processing.RequestQueryStore)
	if err != nil {
		return nil, fmt.Errorf("Proc.NewRequest error: %w", err)
	}

	txHash, err := s.Infra.Index.SendSingle(ctx, packed, indexer.MulticallEhr)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, errors.ErrNotFound
		} else if strings.Contains(err.Error(), "AEX") {
			return nil, errors.ErrAlreadyExist
		}

		return nil, fmt.Errorf("Index.SendSingle error: %w", err)
	}

	procRequest.AddEthereumTx(processing.TxAddEhrDoc, txHash)

	if err := procRequest.Commit(); err != nil {
		return nil, fmt.Errorf("EHR create procRequest commit error: %w", err)
	}

	return storedQuery, nil
}

func (s *Service) ExecStoredQuery(ctx context.Context, userID, systemID, qualifiedQueryName string, query *model.QueryRequest) (*model.QueryResponse, error) {
	v, _ := base.NewVersionTreeID(defaultVersion)

	storedQuery, err := s.GetByVersion(ctx, userID, systemID, qualifiedQueryName, v)
	if err != nil {
		return nil, errors.Wrap(err, "cannot find stored query")
	}

	query.Query = storedQuery.Query

	resp, err := s.qExec.ExecQuery(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "cannot exec query")
	}

	resp.Name = qualifiedQueryName

	return resp, nil
}

func (s *Service) ExecQuery(ctx context.Context, query *model.QueryRequest) (*model.QueryResponse, error) {
	resp, err := s.qExec.ExecQuery(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "cannot exec query")
	}

	return resp, nil
}
