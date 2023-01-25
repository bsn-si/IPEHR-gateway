package docAccess

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
	proc "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer"
)

type Service struct {
	*service.DefaultDocumentService
}

func NewService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		docService,
	}
}

func (s *Service) List(ctx context.Context, userID, systemID string) (*model.DocAccessListResponse, error) {
	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("keystore.Get error: %w userID %s", err, userID)
	}

	result := model.DocAccessListResponse{
		Documents:      []*model.DocAccessDocument{},
		DocumentGroups: []*model.DocAccessDocumentGroup{},
	}

	result.Documents, err = s.getDocumentAccess(ctx, userID, systemID, userPubKey, userPrivKey)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("getDocumentAccess error: %w", err)
	}

	result.DocumentGroups, err = s.getDocumentGroupsAccess(ctx, userID, systemID, userPubKey, userPrivKey)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("getDocumentGroupsAccess error: %w", err)
	}

	userGroupsAccess, err := s.getDocumentUserGroupsAccess(ctx, userID, systemID, userPubKey, userPrivKey)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("getDocumentUserGroupsAccess error: %w", err)
	}

	result.DocumentGroups = append(result.DocumentGroups, userGroupsAccess...)

	return &result, nil
}

func (s *Service) Set(ctx context.Context, userID, systemID, toUserID, reqID string, CID *cid.Cid, accessLevel uint8) error {
	_, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return fmt.Errorf("keystore.Get error: %w userID %s", err, userID)
	}

	toUserPubKey, toUserPrivKey, err := s.Infra.Keystore.Get(toUserID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	var keyEncr, CIDEncr []byte
	{
		docAccessKey, err := s.GetDocAccessKey(ctx, userID, systemID, CID)
		if err != nil {
			return fmt.Errorf("Index.GetDocKeyEncrypted error: %w", err)
		}

		keyEncr, err = keybox.SealAnonymous(docAccessKey.Bytes(), toUserPubKey)
		if err != nil {
			return fmt.Errorf("keybox.SealAnonymous error: %w", err)
		}

		CIDEncr, err = keybox.SealAnonymous(CID.Bytes(), toUserPubKey)
		if err != nil {
			return fmt.Errorf("keybox.SealAnonymous error: %w", err)
		}
	}

	data, err := s.Infra.Index.DocAccessSet(ctx, CID.Bytes(), CIDEncr, keyEncr, accessLevel, userPrivKey, toUserPrivKey, nil)
	if err != nil {
		return fmt.Errorf("Index.DocAccessSet error: %w", err)
	}

	txHash, err := s.Infra.Index.SendSingle(ctx, data, indexer.MulticallEhr)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return errors.ErrNotFound
		}

		return fmt.Errorf("Index.SendSingle error: %w", err)
	}

	procRequest, err := s.Proc.NewRequest(reqID, userID, "", proc.RequestDocAccessSet)
	if err != nil {
		return fmt.Errorf("Proc.NewRequest error: %w", err)
	}

	procRequest.AddEthereumTx(proc.TxSetDocAccess, txHash)

	return nil
}

func (s *Service) getDocumentAccess(ctx context.Context, userID, systemID string, userPubKey, userPrivKey *[32]byte) ([]*model.DocAccessDocument, error) {
	IDHash := sha3.Sum256([]byte(userID + systemID))

	acl, err := s.Infra.Index.GetAccessList(ctx, &IDHash, access.Doc)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("Index.GetAccessList documents error: %w userID: %s", err, userID)
	}

	var documents []*model.DocAccessDocument

	for i, a := range acl {
		err = access.ExtractWithUserKey(a, userPubKey, userPrivKey)
		if err != nil {
			return nil, fmt.Errorf("index: %d access.Extract doc error: %w", i, err)
		}

		CID, err := cid.Parse(a.ID)
		if err != nil {
			return nil, fmt.Errorf("cid.Parse error: %w id: %x", err, a.ID)
		}

		//TODO doc description

		documents = append(documents, &model.DocAccessDocument{
			CID:   CID.String(),
			Level: a.Level,
		})
	}

	return documents, nil
}

func (s *Service) getDocumentGroupsAccess(ctx context.Context, userID, systemID string, userPubKey, userPrivKey *[32]byte) ([]*model.DocAccessDocumentGroup, error) {
	IDHash := sha3.Sum256([]byte(userID + systemID))

	acl, err := s.Infra.Index.GetAccessList(ctx, &IDHash, access.DocGroup)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("Index.GetAccessList document groups error: %w userID: %s", err, userID)
	}

	var docGroups []*model.DocAccessDocumentGroup

	for i, a := range acl {
		err = access.ExtractWithUserKey(a, userPubKey, userPrivKey)
		if err != nil {
			return nil, fmt.Errorf("index: %d access.Extract doc groups error: %w", i, err)
		}

		groupID, err := uuid.FromBytes(a.ID)
		if err != nil {
			return nil, fmt.Errorf("groupID UUID parse error: %w", err)
		}

		docGroups = append(docGroups, &model.DocAccessDocumentGroup{
			GroupID: groupID.String(),
			Level:   a.Level,
		})
	}

	return docGroups, nil
}

func (s *Service) getDocumentUserGroupsAccess(ctx context.Context, userID, systemID string, userPubKey, userPrivKey *[32]byte) ([]*model.DocAccessDocumentGroup, error) {
	IDHash := sha3.Sum256([]byte(userID + systemID))

	acl, err := s.Infra.Index.GetAccessList(ctx, &IDHash, access.UserGroup)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("Index.GetAccessList user groups error: %w userID: %s", err, userID)
	}

	var docGroups []*model.DocAccessDocumentGroup

	for i, a := range acl {
		err := access.ExtractWithUserKey(a, userPubKey, userPrivKey)
		if err != nil {
			if errors.Is(err, errors.ErrAccessDenied) {
				continue
			}

			return nil, fmt.Errorf("index: %d access.Extract user groups error: %w", i, err)
		}

		userGroupID, err := uuid.FromBytes(a.ID)
		if err != nil {
			return nil, fmt.Errorf("userGroupID uuid.FromBytes error: %w", err)
		}

		IDHash = sha3.Sum256(userGroupID[:])

		docGroupACL, err := s.Infra.Index.GetAccessList(ctx, &IDHash, access.DocGroup)
		if err != nil {
			if errors.Is(err, errors.ErrNotFound) {
				return nil, err
			}

			return nil, fmt.Errorf("Index.GetAccessList document groups error: %w userID: %s", err, userID)
		}

		for j, ga := range docGroupACL {
			// Getting docGroup IDs
			err = access.ExtractWithGroupKey(ga, a.Key)
			if err != nil {
				return nil, fmt.Errorf("index %d: access.ExtractWithGroupKey error: %w", j, err)
			}

			groupID, err := uuid.FromBytes(a.ID)
			if err != nil {
				return nil, fmt.Errorf("index %d: uuid.FromBytes error: %w", j, err)
			}

			docGroup := &model.DocAccessDocumentGroup{
				GroupID:       groupID.String(),
				Level:         a.Level,
				ParentGroupID: userGroupID.String(),
			}

			// Getting doc IDs from doGroups
			CIDsEncr, err := s.Infra.Index.DocGroupGetDocs(ctx, &groupID)
			if err != nil {
				return nil, fmt.Errorf("Index.DocGroupGetDocs error: %w", err)
			}

			for k, CIDEncr := range CIDsEncr {
				CIDBytes, err := a.Key.Decrypt(CIDEncr)
				if err != nil {
					return nil, fmt.Errorf("index %d CID decryption error: %w", k, err)
				}

				CID, err := cid.Parse(CIDBytes)
				if err != nil {
					return nil, fmt.Errorf("index %d cid.Parse error: %w CIDBytes: %x", k, err, CIDBytes)
				}

				docGroup.Documents = append(docGroup.Documents, &model.DocAccessDocument{CID: CID.String()})
			}

			docGroups = append(docGroups, docGroup)
		}
	}

	return docGroups, nil
}
