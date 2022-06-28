package group_access

import (
	"log"

	"github.com/google/uuid"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/errors"
)

type GroupAccessService struct {
	Doc                    *service.DefaultDocumentService
	DefaultGroupAccessUUID *uuid.UUID
}

func NewGroupAccessService(docService *service.DefaultDocumentService, cfg *config.Config) *GroupAccessService {
	groupUUID, err := uuid.Parse(cfg.DefaultGroupAccessId)
	if err != nil {
		log.Fatal(err)
	}

	_, err = uuid.Parse(cfg.DefaultUserId)
	if err != nil {
		log.Fatal(err)
	}

	groupAccess, err := docService.GroupAccessIndex.Get(cfg.DefaultUserId, &groupUUID)
	if err != nil {
		if errors.Is(err, errors.IsNotExist) {
			groupAccess = &model.GroupAccess{
				GroupUUID:   &groupUUID,
				Description: "Default access group",
				Key:         chacha_poly.GenerateKey(),
			}

			err = docService.GroupAccessIndex.Add(cfg.DefaultUserId, groupAccess)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
	}

	return &GroupAccessService{
		Doc:                    docService,
		DefaultGroupAccessUUID: &groupUUID,
	}
}

func (g *GroupAccessService) Get(userId string, groupUUID *uuid.UUID) (groupAccess *model.GroupAccess, err error) {
	return g.Doc.GroupAccessIndex.Get(userId, groupUUID)
}

func (g *GroupAccessService) Create(userId string, c *model.GroupAccessCreateRequest) (groupAccess *model.GroupAccess, err error) {
	groupUUID, err := uuid.NewUUID()
	if err != nil {
		return
	}

	groupAccess = &model.GroupAccess{
		GroupUUID:   &groupUUID,
		Description: c.Description,
		Key:         chacha_poly.GenerateKey(),
	}

	err = g.Doc.GroupAccessIndex.Add(userId, groupAccess)

	return
}
