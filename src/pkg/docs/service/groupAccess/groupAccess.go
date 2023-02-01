package groupAccess

import (
	"log"

	"github.com/google/uuid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
)

type Service struct {
	*service.DefaultDocumentService
	defaultGroupAccessUUID *uuid.UUID
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

	/*
		ctx := context.Background()

		groupAccess, err := service.Get(ctx, defaultUserID, &groupUUID)
		if err != nil {
			if errors.Is(err, errors.ErrIsNotExist) {
				log.Println("Default access group is not registered.")
			} else {
				log.Fatal(err)
			}
		}

		service.defaultGroupAccess = groupAccess
	*/

	return service
}

func (s *Service) Default() *uuid.UUID {
	return s.defaultGroupAccessUUID
}

/*
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

	_, err = s.Infra.Index.SetGroupAccess(ctx, &h, groupAccessEncrypted, uint8(access.Owner), userPrivKey, nil)
	if err != nil {
		return fmt.Errorf("Index.SetGroupAccess error: %w", err)
	}

	return nil
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
*/
