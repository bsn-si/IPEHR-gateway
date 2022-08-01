package groupAccess

import (
	"crypto/rand"
	"log"

	"github.com/google/uuid"

	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/errors"
)

type Service struct {
	Doc                    *service.DefaultDocumentService
	DefaultGroupAccessUUID *uuid.UUID
}

func NewGroupAccessService(docService *service.DefaultDocumentService, defaultGroupAccessID, defaultUserID string) *Service {
	groupUUID, err := uuid.Parse(defaultGroupAccessID)
	if err != nil {
		log.Fatal(err)
	}

	_, err = uuid.Parse(defaultUserID)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = docService.GroupAccessIndex.Get(defaultUserID, &groupUUID); err != nil {
		if errors.Is(err, errors.ErrIsNotExist) {
			groupAccess := &model.GroupAccess{
				GroupUUID:   &groupUUID,
				Description: "Default access group",
				Key:         chachaPoly.GenerateKey(),
			}

			err = docService.GroupAccessIndex.Add(defaultUserID, groupAccess)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	return &Service{
		Doc:                    docService,
		DefaultGroupAccessUUID: &groupUUID,
	}
}

func (g *Service) Get(userID string, groupUUID *uuid.UUID) (groupAccess *model.GroupAccess, err error) {
	return g.Doc.GroupAccessIndex.Get(userID, groupUUID)
}

func (g *Service) Create(userID string, c *model.GroupAccessCreateRequest) (groupAccess *model.GroupAccess, err error) {
	groupUUID, err := uuid.NewUUID()
	if err != nil {
		return
	}

	groupAccess = &model.GroupAccess{
		GroupUUID:   &groupUUID,
		Description: c.Description,
		Key:         chachaPoly.GenerateKey(),
		Nonce:       &[12]byte{},
	}
	if _, err := rand.Read(groupAccess.Nonce[:]); err != nil {
		return nil, err
	}

	err = g.Doc.GroupAccessIndex.Add(userID, groupAccess)

	return
}
