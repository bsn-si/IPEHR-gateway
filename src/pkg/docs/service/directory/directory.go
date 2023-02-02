package directory

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
	docGroupService "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/docGroup"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	proc "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/status"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/types"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/ehrIndexer"
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

func (s *Service) GetActiveProcRequest(userID string, kind processing.RequestKind) (string, error) {
	r, err := s.Proc.GetRequestsByKindInProgress(userID, kind)
	if err != nil {
		return "", err
	}

	keys := reflect.ValueOf(r).MapKeys()
	return keys[0].String(), nil
}

func (s *Service) Create(ctx context.Context, req processing.RequestInterface, patientID, systemID, dirUID string, d *model.Directory) error {
	docBytes, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("DIRECTORY marshal error: %w", err)
	}

	allDocGroup, err := s.getGroupAllDocs(ctx, patientID, systemID)
	if err != nil {
		return err
	}

	err = s.save(ctx, req, docBytes, patientID, systemID, dirUID, d.Name.Value, allDocGroup)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) getGroupAllDocs(ctx context.Context, patientID string, systemID string) (*model.DocumentGroup, error) {
	return s.DocGroup.GroupGetByName(ctx, common.DefaultGroupAllDocuments, patientID, systemID)
}

func (s *Service) save(ctx context.Context, req proc.RequestInterface, docBytes []byte, patientID, systemID, dirUID, encName string, allDocGroup *model.DocumentGroup) error {
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

	CID, dealCID, minerAddr, err := s.fileCoinSaving(ctx, req, docBytes, key, objectVersionID)
	if err != nil {
		return err
	}

	if err := s.addMetaData(ctx, req, key, objectVersionID, CID, userPubKey, encName, baseDocumentUIDHash, dealCID, minerAddr, userPrivKey); err != nil {
		return err
	}

	if err := s.addDocGroupData(ctx, req, CID, allDocGroup, userPrivKey); err != nil {
		return err
	}

	return nil
}

func (s *Service) fileCoinSaving(ctx context.Context, req proc.RequestInterface, docBytes []byte, key *chachaPoly.Key, objectVersionID *base.ObjectVersionID) (*cid.Cid, *cid.Cid, string, error) {
	var err error

	if s.Infra.CompressionEnabled {
		docBytes, err = s.Infra.Compressor.Compress(docBytes)
		if err != nil {
			return nil, nil, "", fmt.Errorf("DIRECTORY compress error: %w", err)
		}
	}

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(objectVersionID.String()))
	if err != nil {
		return nil, nil, "", fmt.Errorf("DIRECTORY encryption error: %w", err)
	}

	// IPFS saving
	CID, err := s.Infra.IpfsClient.Add(ctx, docEncrypted)
	if err != nil {
		return nil, nil, "", fmt.Errorf("IpfsClient.Add error: %w", err)
	}

	// Filecoin saving
	dealCID, minerAddr, err := s.Infra.FilecoinClient.StartDeal(ctx, CID, uint64(len(docEncrypted)))
	if err != nil {
		return nil, nil, "", fmt.Errorf("FilecoinClient.StartDeal error: %w", err)
	}

	req.AddFilecoinTx(proc.TxCreateDirectory, CID.String(), dealCID.String(), minerAddr)

	return CID, dealCID, minerAddr, nil
}

func (s *Service) addMetaData(ctx context.Context, req proc.RequestInterface, key *chachaPoly.Key, objectVersionID *base.ObjectVersionID, CID *cid.Cid, userPubKey *[32]byte, encName string, baseDocumentUIDHash [32]byte, dealCID *cid.Cid, minerAddr string, userPrivKey *[32]byte) error {
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

		nameEncr, err := key.Encrypt([]byte(encName))
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
	return nil
}

func (s *Service) addDocGroupData(ctx context.Context, req proc.RequestInterface, CID *cid.Cid, allDocGroup *model.DocumentGroup, userPrivKey *[32]byte) error {
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

func (s *Service) Update(ctx context.Context, req processing.RequestInterface, systemID string, userID string, d *model.Directory) error {
	if _, err := s.IncreaseVersion(d, systemID); err != nil {
		return fmt.Errorf("Directory IncreaseVersion error: %w directory.UID %s", err, d.UID.Value)
	}

	docBytes, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("DIRECTORY marshal error: %w", err)
	}

	allDocGroup, err := s.getGroupAllDocs(ctx, userID, systemID)
	if err != nil {
		return err
	}

	err = s.save(ctx, req, docBytes, userID, systemID, d.UID.Value, d.Name.Value, allDocGroup)
	if err != nil {
		return err
	}

	return nil
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

func (s *Service) GetByTimeOrLast(ctx context.Context, systemID string, ehrUUID *uuid.UUID, patientID string, versionTime time.Time) (*model.Directory, error) {
	docMeta, err := s.Infra.Index.GetDocByTime(ctx, ehrUUID, types.Directory, uint32(versionTime.Unix()))

	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			docMeta, err = s.Infra.Index.GetDocLastByType(ctx, ehrUUID, types.Directory)
			if err != nil {
				return nil, fmt.Errorf("Index.GetDocLastByType error: %w", err)
			}
		} else {
			return nil, fmt.Errorf("Index.GetDocByTime error: %w ehrID %s nearestTime %s docType %s", err, ehrUUID.String(), versionTime.String(), types.Directory)
		}
	}

	CID, err := cid.Parse(docMeta.Id)
	if err != nil {
		return nil, fmt.Errorf("cid.Parse error: %w", err)
	}

	docUIDEncrypted := docMeta.GetAttr(model.AttributeDocUIDEncr)
	if docUIDEncrypted == nil {
		return nil, errors.ErrFieldIsEmpty("DocUIDEncrypted")
	}

	docDecrypted, err := s.DocGroup.GetDocFromStorageByID(ctx, patientID, systemID, &CID, docMeta.Version, docUIDEncrypted)
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

	if docMeta.Status == uint8(status.DELETED) {
		return &d, errors.ErrAlreadyDeleted
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

func (s *Service) IncreaseVersion(d *model.Directory, systemID string) (string, error) {
	dVersionUID, err := base.NewObjectVersionID(d.UID.Value, systemID)
	if err != nil {
		return "", fmt.Errorf("Directory %s NewObjectVersionID error: %w", d.UID.Value, err)
	}

	if _, err := dVersionUID.IncreaseUIDVersion(); err != nil {
		return "", fmt.Errorf("Directory %s IncreaseUIDVersion error: %w", d.UID.Value, err)
	}

	d.UID.Value = dVersionUID.String()

	return d.UID.Value, nil
}
