package groupAccess

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/errors"
)

type Service struct {
	*service.DefaultDocumentService
	defaultGroupAccessUUID *uuid.UUID
	defaultGroupAccess     *model.GroupAccess
}

func NewService(docService *service.DefaultDocumentService, defaultGroupAccessID, defaultUserID string) *Service {
	groupUUID, err := uuid.Parse(defaultGroupAccessID)
	if err != nil {
		log.Fatal(err)
	}

	service := &Service{
		DefaultDocumentService: docService,
		defaultGroupAccessUUID: &groupUUID,
	}

	_, err = uuid.Parse(defaultUserID)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	groupAccess, err := service.Get(ctx, defaultUserID, &groupUUID)
	if err != nil {
		if errors.Is(err, errors.ErrIsNotExist) {
			groupAccess = &model.GroupAccess{
				GroupUUID:   &groupUUID,
				Description: "Default access group",
				Key:         chachaPoly.GenerateKey(),
				Nonce:       &[12]byte{},
			}

			if _, err := rand.Read(groupAccess.Nonce[:]); err != nil {
				log.Fatal(err)
			}

			if err = service.save(ctx, defaultUserID, groupAccess); err != nil {
				log.Fatal(err)
			}

			log.Println("defaultUserID:", defaultUserID)

		} else {
			log.Fatal(err)
		}
	}

	service.defaultGroupAccess = groupAccess

	return service
}

func (s *Service) Default() *model.GroupAccess {
	return s.defaultGroupAccess
}

func (s *Service) Create(ctx context.Context, userID string, c *model.GroupAccessCreateRequest) (*model.GroupAccess, error) {
	groupAccessUUID := uuid.New()
	groupAccess := &model.GroupAccess{
		GroupUUID:   &groupAccessUUID,
		Description: c.Description,
		Key:         chachaPoly.GenerateKey(),
		Nonce:       new([12]byte),
	}

	if _, err := rand.Read(groupAccess.Nonce[:]); err != nil {
		return nil, err
	}

	if err := s.save(ctx, userID, groupAccess); err != nil {
		return nil, fmt.Errorf("groupAccess save error: %w", err)
	}

	return groupAccess, nil
}

func (s *Service) save(ctx context.Context, userID string, groupAccess *model.GroupAccess) error {
	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	groupAccessByte, err := msgpack.Marshal(groupAccess)
	if err != nil {
		return fmt.Errorf("msgpack.Marshal error: %w", err)
	}

	groupAccessEncrypted, err := keybox.Seal(groupAccessByte, userPubKey, userPrivKey)
	if err != nil {
		return fmt.Errorf("keybox.SealAnonymous error: %w", err)
	}

	h := sha3.Sum256(append([]byte(userID), groupAccess.GroupUUID[:]...))

	txHash, err := s.Infra.Index.SetGroupAccess(ctx, &h, groupAccessEncrypted)
	if err != nil {
		return fmt.Errorf("Index.SetGroupAccess error: %w", err)
	}

	/* DEBUG
	log.Printf("SetGroupAccess txHash: %s", txHash)
	log.Printf("groupAccessEncrypted: %x", groupAccessEncrypted)
	log.Printf("userPubKey: %x", userPubKey)
	log.Printf("userPrivKey: %x", userPrivKey)
	*/

	txStatus, err := s.Infra.Index.TxWait(ctx, txHash)
	if err != nil {
		return fmt.Errorf("index.TxWait error: %w txHash %s", err, txHash)
	}

	if txStatus == 1 {
		return nil
	}

	return fmt.Errorf("%w: tx %s Failed", errors.ErrCustom, txHash)
}

func (s *Service) Get(ctx context.Context, userID string, groupAccessUUID *uuid.UUID) (*model.GroupAccess, error) {
	groupAccessBytes, err := s.Infra.Index.GetGroupAccess(ctx, userID, groupAccessUUID)
	if err != nil {
		return nil, fmt.Errorf("Index.GetGroupAccess error: %w", err)
	}

	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	/* DEBUG
	log.Printf("userID: %s", userID)
	log.Printf("groupAccessBytes: %x", groupAccessBytes)
	log.Printf("userPubKey: %x", *userPubKey)
	log.Printf("userPrivKey: %x", *userPrivKey)
	*/

	groupAccessBytes, err = keybox.Open(groupAccessBytes, userPubKey, userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("keybox.OpenAnonymous error: %w", err)
	}

	var groupAccess model.GroupAccess
	if err = msgpack.Unmarshal(groupAccessBytes, &groupAccess); err != nil {
		return nil, fmt.Errorf("msgpack.Unmarshal error: %w", err)
	}

	return &groupAccess, nil
}
