package directory

import (
	"context"
	"fmt"
	"hms/gateway/pkg/docs/service"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/types"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/infrastructure"
	userModel "github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
	"hms/gateway/pkg/compressor"
	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/docs/status"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/indexer/ehrIndexer"
	userModel "hms/gateway/pkg/user/model"
)

type Service struct {
	*service.DefaultDocumentService
}

func NewService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		docService,
	}
}
func (s *Service) NewProcRequest(reqID, userID, ehrUUID string, kind processing.RequestKind) (processing.RequestInterface, error) {
	return s.Proc.NewRequest(reqID, userID, ehrUUID, kind)
}

func (s *Service) Create(ctx context.Context, req processing.RequestInterface, ehrUUID, patientID, dirUID string, d *model.Directory) error {
	timestamp := time.Now()

	id := []byte(ehrUUID + patientID + dirUID)
	idHash := sha3.Sum256(id)
	key := chachaPoly.GenerateKey()

	content, err := msgpack.Marshal(d)
	if err != nil {
		return fmt.Errorf("msgpack.Marshal error: %w", err)
	}

	contentCompresed, err := compressor.New(compressor.BestCompression).Compress(content)
	if err != nil {
		return fmt.Errorf("Query Compress error: %w", err)
	}

	contentEncr, err := key.Encrypt(contentCompresed)
	if err != nil {
		return fmt.Errorf("key.Encrypt content error: %w", err)
	}

	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(patientID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, patientID)
	}

	keyEncr, err := keybox.SealAnonymous(key.Bytes(), userPubKey)
	if err != nil {
		return fmt.Errorf("keybox.SealAnonymous error: %w", err)
	}

	docMeta := &model.DocumentMeta{
		Status:    uint8(status.ACTIVE),
		Id:        idHash[:],
		Version:   []byte(d.UID.Value),
		Timestamp: uint32(timestamp.Unix()),
		IsLast:    true,
		Attrs: []ehrIndexer.AttributesAttribute{
			{Code: model.AttributeKeyEncr, Value: keyEncr},         // encrypted with key
			{Code: model.AttributeContentEncr, Value: contentEncr}, // encrypted with userPubKey
			{Code: model.AttributeDocUIDHash, Value: idHash[:]},
		},
	}

	packed, err := s.Infra.Index.AddEhrDoc(ctx, types.Directory, docMeta, userPrivKey, nil)
	if err != nil {
		return fmt.Errorf("Index.AddEhrDoc error: %w", err)
	}

	txHash, err := s.Infra.Index.SendSingle(ctx, packed, indexer.MulticallEhr)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return errors.ErrNotFound
		} else if strings.Contains(err.Error(), "AEX") {
			return errors.ErrAlreadyExist
		}

		return fmt.Errorf("Index.SendSingle error: %w", err)
	}

	req.AddEthereumTx(processing.TxAddEhrDoc, txHash)

	return nil
}

// TODO
func (s *Service) Update(ctx context.Context, req processing.RequestInterface, systemID string, ehrUUID *uuid.UUID, user *userModel.UserInfo, d *model.Directory) error {
	if err := s.increaseVersion(d); err != nil {
		return fmt.Errorf("Directory increaseVersion error: %w directory.UID %s", err, d.UID.Value)
	}

	// TODO need realization
	//err = s.save(ctx, multiCallTx, procRequest, userID, systemID, ehrUUID, groupAccess, d)

	return errors.ErrNotImplemented
}

func (s *Service) Delete(ctx context.Context, req processing.RequestInterface, systemID string, ehrUUID *uuid.UUID, versionUID, userID string) (string, error) {
	objectVersionID, err := base.NewObjectVersionID(versionUID, systemID)
	if err != nil {
		return "", fmt.Errorf("NewObjectVersionID error: %w versionUID %s ehrSystemID %s", err, versionUID, systemID)
	}

	_, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return "", fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	baseDocumentUID := []byte(objectVersionID.BasedID())
	baseDocumentUIDHash := sha3.Sum256(baseDocumentUID)

	txHash, err := s.Infra.Index.DeleteDoc(ctx, ehrUUID, types.Directory, &baseDocumentUIDHash, objectVersionID.VersionBytes(), userPrivKey, nil)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return "", err
		}
		return "", fmt.Errorf("Index.DeleteDoc error: %w", err)
	}

	req.AddEthereumTx(processing.TxDeleteDoc, txHash)

	if _, err = objectVersionID.IncreaseUIDVersion(); err != nil {
		return "", fmt.Errorf("IncreaseUIDVersion error: %w objectVersionID %s", err, objectVersionID.String())
	}

	return objectVersionID.String(), nil
}

// TODO
func (s *Service) GetByTime(ctx context.Context, systemID string, ehrUUID *uuid.UUID, userID string, versionTime time.Time) (*model.Directory, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Service) GetByID(ctx context.Context, patientID string, ehrUUID, versionUID string) (*model.Directory, error) {
	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(patientID)
	if err != nil {
		return nil, fmt.Errorf("Keystore.Get error: %w patientID %s", err, patientID)
	}

	id := []byte(ehrUUID + patientID + versionUID)

	idHash := sha3.Sum256(id)

	var vID [32]byte

	copy(vID[:], versionUID)

	e := uuid.MustParse(ehrUUID)

	docMeta, err := s.Infra.Index.GetDocByVersion(ctx, &e, types.Directory, &idHash, &vID)
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

	var d model.Directory

	err = msgpack.Unmarshal(content, &d)
	if err != nil {
		return nil, fmt.Errorf("DIRECTORY content unmarshal error: %w", err)
	}

	return &d, nil
}

func (s *Service) increaseVersion(d *model.Directory) error {
	if _, err := d.IncreaseUIDVersion(); err != nil {
		return fmt.Errorf("Directory %s IncreaseUIDVersion error: %w", d.UID.Value, err)
	}

	return nil
}
