package directory

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/types"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/infrastructure"
	userModel "github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
	"hms/gateway/pkg/common"
	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/service"
	docGroupService "hms/gateway/pkg/docs/service/docGroup"
	"hms/gateway/pkg/docs/service/processing"
	proc "hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/docs/status"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/indexer/ehrIndexer"
	userModel "hms/gateway/pkg/user/model"
)

type Service struct {
	*service.DefaultDocumentService
	DocGroup *docGroupService.Service
}

func NewService(docService *service.DefaultDocumentService, docGroupSvc *docGroupService.Service) *Service {
	return &Service{
		docService,
		docGroupSvc,
	}
}
func (s *Service) NewProcRequest(reqID, userID, ehrUUID string, kind processing.RequestKind) (processing.RequestInterface, error) {
	return s.Proc.NewRequest(reqID, userID, ehrUUID, kind)
}

func (s *Service) Create(ctx context.Context, req processing.RequestInterface, ehrUUID, patientID, systemID, dirUID string, d *model.Directory) error {
	key := chachaPoly.GenerateKey()

	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(patientID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, patientID)
	}

	objectVersionID, err := base.NewObjectVersionID(dirUID, systemID)
	if err != nil {
		return fmt.Errorf("saving error: %w versionUID %s ehrSystemID %s", err, objectVersionID, systemID)
	}

	baseDocumentUID := []byte(objectVersionID.BasedID())
	baseDocumentUIDHash := sha3.Sum256(baseDocumentUID)

	// Searching 'all documents' group
	var allDocGroup *model.DocumentGroup
	{
		docGroups, err := s.DocGroup.GroupGetList(ctx, patientID, systemID)
		if err != nil {
			return fmt.Errorf("DocGroup.GroupGetList error: %w", err)
		}

		for _, dg := range docGroups {
			if dg.Name == common.DefaultGroupAllDocuments {
				allDocGroup = dg
				break
			}
		}

		if allDocGroup == nil {
			return fmt.Errorf("user 'all documents' group not found: %w", errors.ErrNotFound)
		}
	}

	docBytes, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("DIRECTORY marshal error: %w", err)
	}

	if s.Infra.CompressionEnabled {
		docBytes, err = s.Infra.Compressor.Compress(docBytes)
		if err != nil {
			return fmt.Errorf("DIRECTORY compress error: %w", err)
		}
	}

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(objectVersionID.String()))
	if err != nil {
		return fmt.Errorf("DIRECTORY encryption error: %w", err)
	}

	// IPFS saving
	CID, err := s.Infra.IpfsClient.Add(ctx, docEncrypted)
	if err != nil {
		return fmt.Errorf("IpfsClient.Add error: %w", err)
	}

	// Filecoin saving
	dealCID, minerAddr, err := s.Infra.FilecoinClient.StartDeal(ctx, CID, uint64(len(docEncrypted)))
	if err != nil {
		return fmt.Errorf("FilecoinClient.StartDeal error: %w", err)
	}

	req.AddFilecoinTx(proc.TxCreateDirectory, CID.String(), dealCID.String(), minerAddr)

	{
		docIDEncrypted, err := key.Encrypt([]byte(objectVersionID.String()))
		if err != nil {
			return fmt.Errorf("EncryptWithAuthData error: %w", err)
		}

		CIDEncr, err := key.Encrypt(CID.Bytes())
		if err != nil {
			return fmt.Errorf("CID encryption error: %w", err)
		}

		keyEncr, err := keybox.SealAnonymous(key.Bytes(), userPubKey)
		if err != nil {
			return fmt.Errorf("keybox.SealAnonymous error: %w", err)
		}

		nameEncr, err := key.Encrypt([]byte(d.Name.Value))
		if err != nil {
			return fmt.Errorf("Encrypt name error: %w", err)
		}

		docMeta := &model.DocumentMeta{
			Status:    uint8(status.ACTIVE),
			Id:        CID.Bytes(),
			Version:   objectVersionID.VersionBytes()[:],
			Timestamp: uint32(time.Now().Unix()),
			IsLast:    true,
			Attrs: []ehrIndexer.AttributesAttribute{
				{Code: model.AttributeIDEncr, Value: CIDEncr},
				{Code: model.AttributeKeyEncr, Value: keyEncr},
				{Code: model.AttributeDocUIDHash, Value: baseDocumentUIDHash[:]},
				{Code: model.AttributeDocUIDEncr, Value: docIDEncrypted},
				{Code: model.AttributeDealCid, Value: dealCID.Bytes()},
				{Code: model.AttributeMinerAddress, Value: []byte(minerAddr)},
				{Code: model.AttributeNameEncr, Value: nameEncr},
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
	}

	{
		docCIDHash := indexer.Keccak256(CID.Bytes())

		docCIDEncr, err := allDocGroup.GroupKey.Encrypt(CID.Bytes())
		if err != nil {
			return fmt.Errorf("CID encryption error: %w", err)
		}

		packed, err := s.Infra.Index.DocGroupAddDoc(ctx, &allDocGroup.GroupID, docCIDHash, docCIDEncr, userPrivKey, nil)
		if err != nil {
			return fmt.Errorf("Index.DocGroupAddDoc error: %w", err)
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

		req.AddEthereumTx(processing.TxDocGroupAddDoc, txHash)
	}

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

func (s *Service) GetByTime(ctx context.Context, systemID string, ehrUUID *uuid.UUID, userID string, versionTime time.Time) (*model.Directory, error) {
	docMeta, err := s.Infra.Index.GetDocByTime(ctx, ehrUUID, types.Directory, uint32(versionTime.Unix()))
	if err != nil && errors.Is(err, errors.ErrNotFound) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("Index.GetDocByTime error: %w ehrID %s nearestTime %s docType %s", err, ehrUUID.String(), versionTime.String(), types.EhrStatus)
	}

	CID, err := cid.Parse(docMeta.Id)
	if err != nil {
		return nil, fmt.Errorf("cid.Parse error: %w", err)
	}

	docUIDEncrypted := docMeta.GetAttr(model.AttributeDocUIDEncr)
	if docUIDEncrypted == nil {
		return nil, errors.ErrFieldIsEmpty("DocUIDEncrypted")
	}

	docDecrypted, err := s.DocGroup.GetDocFromStorageByID(ctx, userID, systemID, &CID, docMeta.Version, docUIDEncrypted)
	//docDecrypted, err := s.DocGroup.GetDocFromStorageByID(ctx, userID, systemID, &CID, ehrUUID[:], docUIDEncrypted)
	if err != nil && errors.Is(err, errors.ErrIsInProcessing) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("GetDocFromStorageByID error: %w", err)
	}

	var d model.Directory

	err = json.Unmarshal(docDecrypted, &d)
	if err != nil {
		return nil, fmt.Errorf("DIRECTORY content unmarshal error: %w", err)
	}

	return &d, nil
}

func (s *Service) GetByID(ctx context.Context, patientID string, systemID string, ehrUUID *uuid.UUID, versionID *base.ObjectVersionID) (*model.Directory, error) {
	baseDocumentUID := versionID.BasedID()
	baseDocumentUIDHash := sha3.Sum256([]byte(baseDocumentUID))

	docMeta, err := s.Infra.Index.GetDocByVersion(ctx, ehrUUID, types.Directory, &baseDocumentUIDHash, versionID.VersionBytes())
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("Index.GetDocByVersion error: %w", err)
	}

	if docMeta.Status == uint8(status.DELETED) {
		return nil, errors.ErrAlreadyDeleted
	}

	CID, err := cid.Parse(docMeta.Id)
	if err != nil {
		return nil, fmt.Errorf("cid.Parse error: %w", err)
	}

	docUIDEncrypted := docMeta.GetAttr(model.AttributeDocUIDEncr)
	if docUIDEncrypted == nil {
		return nil, errors.ErrFieldIsEmpty("DocUIDEncrypted")
	}

	docDecrypted, err := s.DocGroup.GetDocFromStorageByID(ctx, patientID, systemID, &CID, []byte(versionID.String()), docUIDEncrypted)
	if err != nil && errors.Is(err, errors.ErrIsInProcessing) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("GetDocFromStorageByID error: %w", err)
	}

	var d model.Directory

	err = json.Unmarshal(docDecrypted, &d)
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
