package composition

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"hms/gateway/pkg/common"
	"hms/gateway/pkg/common/fakeData"
	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/groupAccess"
	"hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/docs/status"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer/service/dataSearch"
)

type Service struct {
	*service.DefaultDocumentService
	DataSearchIndex    *dataSearch.Index
	groupAccessService *groupAccess.Service
}

func NewCompositionService(docService *service.DefaultDocumentService, groupAccessService *groupAccess.Service) *Service {
	return &Service{
		DefaultDocumentService: docService,
		DataSearchIndex:        dataSearch.New(),
		groupAccessService:     groupAccessService,
	}
}

func (s *Service) Create(ctx context.Context, userID string, ehrUUID, groupAccessUUID *uuid.UUID, ehrSystemID base.EhrSystemID, composition *model.Composition) (*model.Composition, error) {
	groupAccessModel, err := s.groupAccessService.Get(ctx, userID, groupAccessUUID)
	if err != nil {
		return nil, fmt.Errorf("groupAccessService.Get error: %w userID %s groupAccessUUID %s", err, userID, groupAccessUUID.String())
	}

	var (
		reqID        = ctx.(*gin.Context).GetString("reqId")
		transactions = s.MultiCallTx.New(s.Infra.Index, s.Proc, processing.TxSetEhrUser, "", reqID)
	)

	err = s.save(ctx, &transactions, userID, ehrUUID, groupAccessModel, ehrSystemID, composition)
	if err != nil {
		return nil, fmt.Errorf("Composition %s save error: %w", composition.UID.Value, err)
	}

	if err := transactions.Commit(); err != nil {
		return nil, fmt.Errorf("Create composition commit error: %w", err)
	}

	return composition, nil
}

func (s *Service) Update(ctx context.Context, userID string, ehrUUID, groupAccessUUID *uuid.UUID, ehrSystemID base.EhrSystemID, composition *model.Composition) (*model.Composition, error) {
	groupAccessModel, err := s.groupAccessService.Get(ctx, userID, groupAccessUUID)
	if err != nil {
		return nil, fmt.Errorf("GroupAccessIndex.Get error: %w userID %s groupAccessUUID %s", err, userID, groupAccessUUID.String())
	}

	if err = s.increaseVersion(composition, ehrSystemID); err != nil {
		return nil, fmt.Errorf("Composition increaseVersion error: %w composition.UID %s", err, composition.UID.Value)
	}

	var (
		reqID        = ctx.(*gin.Context).GetString("reqId")
		transactions = s.MultiCallTx.New(s.Infra.Index, s.Proc, processing.TxSetEhrUser, "", reqID)
	)

	err = s.save(ctx, &transactions, userID, ehrUUID, groupAccessModel, ehrSystemID, composition)
	if err != nil {
		return nil, fmt.Errorf("Composition save error: %w userID %s ehrUUID %s composition.UID %s", err, userID, ehrUUID.String(), composition.UID.Value)
	}

	if err := transactions.Commit(); err != nil {
		return nil, fmt.Errorf("Update composition commit error: %w", err)
	}

	// TODO what we should do with prev composition?
	return composition, nil
}

func (s *Service) increaseVersion(c *model.Composition, ehrSystemID base.EhrSystemID) error {
	if c == nil || c.UID == nil || c.UID.Value == "" {
		return fmt.Errorf("%w Incorrect composition UID", errors.ErrIncorrectFormat)
	}

	objectVersionID, err := base.NewObjectVersionID(c.UID.Value, ehrSystemID)
	if err != nil {
		return fmt.Errorf("increaseVersion error: %w versionUID %s ehrSystemID %s", err, objectVersionID.String(), ehrSystemID.String())
	}

	if _, err := objectVersionID.IncreaseUIDVersion(); err != nil {
		return fmt.Errorf("Composition %s IncreaseUIDVersion error: %w", c.UID.Value, err)
	}

	c.UID.Value = objectVersionID.String()

	return nil
}

func (s *Service) save(ctx context.Context, multiCallTx *processing.MultiCallTx, userID string, ehrUUID *uuid.UUID, groupAccess *model.GroupAccess, ehrSystemID base.EhrSystemID, doc *model.Composition) error {
	objectVersionID, err := base.NewObjectVersionID(doc.UID.Value, ehrSystemID)
	if err != nil {
		return fmt.Errorf("saving error: %w versionUID %s ehrSystemID %s", err, objectVersionID.String(), ehrSystemID.String())
	}

	baseDocumentUID := []byte(objectVersionID.BasedID())
	baseDocumentUIDHash := sha3.Sum256(baseDocumentUID)

	// Checking the existence of the Composition
	docMeta, err := s.Infra.Index.GetDocByVersion(ctx, ehrUUID, types.Composition, &baseDocumentUIDHash, objectVersionID.VersionBytes())
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		return fmt.Errorf("Index.GetDocByVersion error: %w", err)
	} else if docMeta != nil {
		return fmt.Errorf("%w objectVersionID %s", errors.ErrAlreadyExist, objectVersionID.String())
	}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("Composition marshal error: %w", err)
	}

	if s.Infra.CompressionEnabled {
		docBytes, err = s.Infra.Compressor.Compress(docBytes)
		if err != nil {
			return fmt.Errorf("Compress error: %w", err)
		}
	}

	// Document encryption key generation
	key := chachaPoly.GenerateKey()

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(objectVersionID.String()))
	if err != nil {
		return fmt.Errorf("EncryptWithAuthData error: %w", err)
	}

	// IPFS saving
	CID, err := s.Infra.IpfsClient.Add(ctx, docEncrypted)
	if err != nil {
		return fmt.Errorf("IpfsClient.Add error: %w", err)
	}

	// Filecoin saving
	//dealCID, minerAddr, err := s.Infra.FilecoinClient.StartDeal(ctx, CID, uint64(len(docEncrypted)))
	//if err != nil {
	//	return fmt.Errorf("FilecoinClient.StartDeal error: %w", err)
	//}
	dealCID := fakeData.Cid()
	minerAddr := "123"

	docIDEncrypted, err := key.EncryptWithAuthData([]byte(objectVersionID.String()), ehrUUID[:])
	if err != nil {
		return fmt.Errorf("EncryptWithAuthData error: %w", err)
	}

	// Start processing request
	reqID := ctx.(*gin.Context).GetString("reqId")
	{
		procReq := &processing.Request{
			ReqID:        reqID,
			UserID:       userID,
			EhrUUID:      ehrUUID.String(),
			Status:       processing.StatusProcessing,
			Kind:         processing.RequestCompositionCreate,
			CID:          CID.String(),
			DealCID:      dealCID.String(),
			MinerAddress: minerAddr,
		}
		if err = s.Proc.AddRequest(procReq); err != nil {
			return fmt.Errorf("Proc.AddRequest error: %w", err)
		}

		/*
			err = s.Proc.AddTx(reqID, dealCID.String(), "", processing.TxFilecoinStartDeal, processing.StatusPending)
			if err != nil {
				return fmt.Errorf("Proc.AddTx error: %w", err)
			}
		*/
	}

	// Index Docs ehr_id -> doc_meta
	{
		docMeta := &model.DocumentMeta{
			DocType:         uint8(types.Composition),
			Status:          uint8(status.ACTIVE),
			CID:             CID.Bytes(),
			DealCID:         dealCID.Bytes(),
			MinerAddress:    []byte(minerAddr),
			DocUIDEncrypted: docIDEncrypted,
			DocBaseUIDHash:  baseDocumentUIDHash,
			Version:         *objectVersionID.VersionBytes(),
			IsLast:          true,
			Timestamp:       uint32(time.Now().Unix()),
		}

		packed, err := s.Infra.Index.AddEhrDoc(ehrUUID, docMeta)
		if err != nil {
			return fmt.Errorf("Index.AddEhrDoc error: %w", err)
		}
		multiCallTx.Add(packed)
	}

	// Index DataSearch
	_ = groupAccess
	/* TODO
	docStorageIDEncrypted, err := groupAccess.Key.EncryptWithAuthData(cidBytes[:], groupAccess.GroupUUID[:])
	if err != nil {
		return fmt.Errorf("EncryptWithAuthData error: %w", err)
	}

	if err = s.DataSearchIndex.UpdateIndexWithNewContent(doc.Content, groupAccess, docStorageIDEncrypted); err != nil {
		return fmt.Errorf("UpdateIndexWithNewContent error: %w", err)
	}
	*/

	// Index Access
	{
		userPubKey, _, err := s.Infra.Keystore.Get(userID)
		if err != nil {
			return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
		}

		docAccessValue, err := keybox.SealAnonymous(key.Bytes(), userPubKey)
		if err != nil {
			return fmt.Errorf("keybox.SealAnonymous error: %w", err)
		}

		docAccessKey := sha3.Sum256(append(CID.Bytes()[:], []byte(userID)...))

		packed, err := s.Infra.Index.SetDocKeyEncrypted(&docAccessKey, docAccessValue)
		if err != nil {
			return fmt.Errorf("Index.SetDocAccess error: %w", err)
		}

		multiCallTx.Add(packed)
	}

	return nil
}

func (s *Service) GetLastByBaseID(ctx context.Context, userID string, ehrUUID *uuid.UUID, versionUID string, ehrSystemID base.EhrSystemID) (*model.Composition, error) {
	objectVersionID, err := base.NewObjectVersionID(versionUID, ehrSystemID)
	if err != nil {
		return nil, fmt.Errorf("GetLastByBaseID error: %w versionUID %s ehrSystemID %s", err, objectVersionID.String(), ehrSystemID.String())
	}

	baseDocumentUID := []byte(objectVersionID.BasedID())
	baseDocumentUIDHash := sha3.Sum256(baseDocumentUID)

	docMeta, err := s.Infra.Index.GetDocLastByBaseID(ctx, ehrUUID, types.Composition, &baseDocumentUIDHash)
	if err != nil {
		return nil, fmt.Errorf("GetLastVersionDocIndexByBaseID error: %w userID %s objectVersionID %s", err, userID, objectVersionID)
	}

	if docMeta.Status == uint8(status.DELETED) {
		return nil, fmt.Errorf("GetLastByBaseID error: %w", errors.ErrAlreadyDeleted)
	}

	docDecrypted, err := s.GetDocFromStorageByID(ctx, userID, docMeta.Cid(), ehrUUID[:], docMeta.DocUIDEncrypted)
	if err != nil && errors.Is(err, errors.ErrIsInProcessing) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("GetDocFromStorageByID error: %w userID %s storageID %s", err, userID, docMeta.CID)
	}

	var composition *model.Composition
	if err = json.Unmarshal(docDecrypted, &composition); err != nil {
		return nil, fmt.Errorf("Composition unmarshal error: %w", err)
	}

	return composition, nil
}

func (s *Service) GetByID(ctx context.Context, userID string, ehrUUID *uuid.UUID, versionUID string, ehrSystemID base.EhrSystemID) (*model.Composition, error) {
	objectVersionID, err := base.NewObjectVersionID(versionUID, ehrSystemID)
	if err != nil {
		return nil, fmt.Errorf("NewObjectVersionID error: %w versionUID %s ehrSystemID %s", err, versionUID, ehrSystemID.String())
	}

	baseDocumentUID := []byte(objectVersionID.BasedID())
	baseDocumentUIDHash := sha3.Sum256(baseDocumentUID)

	docMeta, err := s.Infra.Index.GetDocByVersion(ctx, ehrUUID, types.Composition, &baseDocumentUIDHash, objectVersionID.VersionBytes())
	if err != nil && errors.Is(err, errors.ErrNotFound) {
		return nil, errors.ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("Index.GetDocByVersion error: %w ehrUUID %s objectVersionID %s", err, ehrUUID.String(), objectVersionID.String())
	}

	if docMeta.Status == uint8(status.DELETED) {
		return nil, fmt.Errorf("GetCompositionByID error: %w", errors.ErrAlreadyDeleted)
	}

	docDecrypted, err := s.GetDocFromStorageByID(ctx, userID, docMeta.Cid(), ehrUUID[:], docMeta.DocUIDEncrypted)
	if err != nil && errors.Is(err, errors.ErrIsInProcessing) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("GetDocFromStorageByID error: %w userID %s CID %x", err, userID, docMeta.Cid().String())
	}

	var composition model.Composition
	if err = json.Unmarshal(docDecrypted, &composition); err != nil {
		return nil, fmt.Errorf("Composition unmarshal error: %w", err)
	}

	return &composition, nil
}

func (s *Service) DeleteByID(ctx context.Context, userID string, ehrUUID *uuid.UUID, versionUID string, ehrSystemID base.EhrSystemID) (string, error) {
	objectVersionID, err := base.NewObjectVersionID(versionUID, ehrSystemID)
	if err != nil {
		return "", fmt.Errorf("NewObjectVersionID error: %w versionUID %s ehrSystemID %s", err, versionUID, ehrSystemID.String())
	}

	baseDocumentUID := []byte(objectVersionID.BasedID())
	baseDocumentUIDHash := sha3.Sum256(baseDocumentUID)

	// Start processing request
	reqID := ctx.(*gin.Context).GetString("reqId")

	procReq := &processing.Request{
		ReqID:       reqID,
		UserID:      userID,
		EhrUUID:     ehrUUID.String(),
		Status:      processing.StatusProcessing,
		Kind:        processing.RequestCompositionDelete,
		BaseUIDHash: hex.EncodeToString(baseDocumentUIDHash[:]),
		Version:     objectVersionID.VersionString(),
	}
	if err = s.Proc.AddRequest(procReq); err != nil {
		return "", fmt.Errorf("Proc.AddRequest error: %w", err)
	}

	docDeleteTx, err := s.Infra.Index.DeleteDoc(ctx, ehrUUID, types.Composition, &baseDocumentUIDHash, objectVersionID.VersionBytes())
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return "", err
		}
		return "", fmt.Errorf("Index.DeleteDoc error: %w", err)
	}

	err = s.Proc.AddTx(reqID, docDeleteTx, "", processing.TxDeleteDoc, processing.StatusPending)
	if err != nil {
		return "", fmt.Errorf("Proc.AddTx error: %w", err)
	}

	// Waiting for tx processed and pending nonce increased
	time.Sleep(common.BlockchainTxProcAwaitTime)

	if _, err = objectVersionID.IncreaseUIDVersion(); err != nil {
		return "", fmt.Errorf("IncreaseUIDVersion error: %w objectVersionID %s", err, objectVersionID.String())
	}

	return objectVersionID.String(), nil
}
