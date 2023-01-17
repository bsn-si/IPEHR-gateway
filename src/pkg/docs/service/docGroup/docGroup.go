package docGroup

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/access"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/errors"
)

type Service struct {
	*service.DefaultDocumentService
}

func NewService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		docService,
	}
}

func (s *Service) GroupGetByID(ctx context.Context, userID, systemID string, groupID *uuid.UUID) (*model.DocumentGroup, error) {
	var err error

	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	docGroup, err := s.Infra.Index.DocGroupGetByID(ctx, groupID, userPubKey, userPrivKey)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("Index.DocGroupGetByID error: %w", err)
	}

	return docGroup, nil
}

func (s *Service) GroupGetList(ctx context.Context, userID, systemID string) ([]*model.DocumentGroup, error) {
	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("keystore.Get error: %w userID %s", err, userID)
	}

	IDHash := sha3.Sum256([]byte(userID + systemID))

	acl, err := s.Infra.Index.GetAccessList(ctx, &IDHash, access.DocGroup)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("Index.GetAccessList error: %w userID: %s", err, userID)
	}

	var docGroupList []*model.DocumentGroup

	for i, a := range acl {
		err = access.ExtractWithUserKey(a, userPubKey, userPrivKey)
		if err != nil {
			return nil, fmt.Errorf("index: %d access.Extract error: %w", i, err)
		}

		groupUUID, err := uuid.FromBytes(a.ID)
		if err != nil {
			return nil, fmt.Errorf("groupID %d uuid.ParseBytes error: %w idDecr: %x", i, err, a.ID)
		}

		docGroup, err := s.GroupGetByID(ctx, userID, systemID, &groupUUID)
		if err != nil {
			return nil, fmt.Errorf("GroupGetByID error: %w groupUUID: %s", err, groupUUID)
		}

		docGroupList = append(docGroupList, docGroup)
	}

	return docGroupList, nil
}
