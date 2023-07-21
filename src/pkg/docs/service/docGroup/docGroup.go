package docGroup

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/internal/observability/tracer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
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
	ctx, span := tracer.Start(ctx, "doc_group_service.GroupGetByID", trace.WithAttributes(
		attribute.String("userID", userID),
		attribute.String("systemID", systemID),
		attribute.String("groupID", groupID.String()),
	))
	defer span.End()

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
	ctx, span := tracer.Start(ctx, "doc_group_service.GroupGetList", trace.WithAttributes(
		attribute.String("userID", userID),
		attribute.String("systemID", systemID),
	))
	defer span.End()

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

func (s *Service) GroupGetByName(ctx context.Context, groupName, userID, systemID string) (*model.DocumentGroup, error) {
	ctx, span := tracer.Start(ctx, "doc_group_service.GroupGetByName", trace.WithAttributes(
		attribute.String("groupName", groupName),
		attribute.String("userID", userID),
		attribute.String("systemID", systemID),
	))
	defer span.End()

	docGroups, err := s.GroupGetList(ctx, userID, systemID)
	if err != nil {
		return nil, fmt.Errorf("DocGroup.GroupGetList error: %w", err)
	}

	for _, dg := range docGroups {
		if dg.Name == groupName {
			return dg, nil
		}
	}

	return nil, fmt.Errorf("user '%s' group not found: %w", groupName, errors.ErrNotFound)

}
