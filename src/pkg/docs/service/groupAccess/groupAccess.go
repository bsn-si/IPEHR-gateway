package groupAccess

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type Service struct {
	*service.DefaultDocumentService
	defaultGroupAccess *model.GroupAccess
}

func NewService(docService *service.DefaultDocumentService, defaultGroupAccessID, defaultUserID string) *Service {
	groupUUID, err := uuid.Parse(defaultGroupAccessID)
	if err != nil {
		log.Fatal(err)
	}

	_, err = uuid.Parse(defaultUserID)
	if err != nil {
		log.Fatal(err)
	}

	service := &Service{
		DefaultDocumentService: docService,
	}

	service.defaultGroupAccess, err = service.Get(context.Background(), defaultUserID, common.EhrSystemID, &groupUUID)
	if err != nil {
		if errors.Is(err, errors.ErrIsNotExist) {
			log.Println("Default access group is not registered.")
		} else {
			log.Fatal(err)
		}
	}

	return service
}

func (s *Service) Default() *model.GroupAccess {
	return s.defaultGroupAccess
}

func (s *Service) Create(ctx context.Context, userID, systemID string, c *model.GroupAccessCreateRequest) (*model.GroupAccess, error) {
	groupAccessUUID := uuid.New()
	groupAccess := &model.GroupAccess{
		UUID:  &groupAccessUUID,
		Key:   chachaPoly.GenerateKey(),
		Nonce: chachaPoly.GenerateNonce(),
	}

	txHash, err := s.save(ctx, userID, systemID, groupAccess)
	if err != nil {
		return nil, fmt.Errorf("groupAccess save error: %w", err)
	}

	log.Printf("GroupAccess %s creation txHash: %s", groupAccess.UUID.String(), txHash)

	return groupAccess, nil
}

func (s *Service) save(ctx context.Context, userID, systemID string, groupAccess *model.GroupAccess) (string, error) {
	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return "", fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	subjectID := sha3.Sum256([]byte(userID + systemID))

	IDEncr, err := keybox.Seal(groupAccess.UUID[:], userPubKey, userPrivKey)
	if err != nil {
		return "", fmt.Errorf("keybox.Seal GroupUUID error: %w", err)
	}

	keyEncr, err := keybox.SealAnonymous(append(groupAccess.Key[:], groupAccess.Nonce[:]...), userPubKey)
	if err != nil {
		return "", fmt.Errorf("keybox.Seal Key error: %w", err)
	}

	txHash, err := s.Infra.Index.SetAccess(ctx, groupAccess.UUID[:], &subjectID, IDEncr, keyEncr, access.GroupAccess, access.Owner)
	if err != nil {
		return "", fmt.Errorf("Index.SetAccess error: %w", err)
	}

	return txHash, nil
}

func (s *Service) Get(ctx context.Context, userID, systemID string, groupAccessUUID *uuid.UUID) (*model.GroupAccess, error) {
	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	keyEncrypted, err := s.Infra.Index.GetKeyEncrypted(ctx, userID, systemID, groupAccessUUID[:], access.GroupAccess)
	if err != nil {
		return nil, fmt.Errorf("Index.GetGroupAccess error: %w", err)
	}

	keyBytes, err := keybox.OpenAnonymous(keyEncrypted, userPubKey, userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("keybox.OpenAnonymous error: %w", err)
	}

	if len(keyBytes) != chachaPoly.KeyLength+chachaPoly.NonceLength {
		return nil, errors.Errorf("Incorrect key length. Expected %d received %d", chachaPoly.KeyLength+chachaPoly.NonceLength, len(keyBytes))
	}

	groupAccess := &model.GroupAccess{
		UUID: groupAccessUUID,
	}

	groupAccess.Key, err = chachaPoly.NewKeyFromBytes(keyBytes[:chachaPoly.KeyLength])
	if err != nil {
		return nil, fmt.Errorf("chachaPoly.NewKeyFromBytes error: %w", err)
	}

	groupAccess.Nonce, err = chachaPoly.NewNonceFromBytes(keyBytes[chachaPoly.KeyLength:])
	if err != nil {
		return nil, fmt.Errorf("chachaPoly.NewKeyFromBytes error: %w", err)
	}

	return groupAccess, nil
}
