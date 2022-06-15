package group_access

import (
	"github.com/google/uuid"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
)

type GroupAccessService struct {
	Doc *service.DefaultDocumentService
	Cfg *config.Config
}

func NewGroupAccessService(docService *service.DefaultDocumentService, cfg *config.Config) *GroupAccessService {
	return &GroupAccessService{
		Doc: docService,
		Cfg: cfg,
	}
}

func (g *GroupAccessService) Get(userId, groupId string) (groupAccess *model.GroupAccess, err error) {
	groupAccess, err = g.Doc.GroupAccessIndex.Get(userId, groupId)
	return
}

func (g *GroupAccessService) Create(userId string, c *model.GroupAccessCreateRequest) (groupAccess *model.GroupAccess, err error) {
	groupAccessId, err := uuid.NewUUID()
	if err != nil {
		return
	}

	groupAccess = &model.GroupAccess{
		GroupId:     groupAccessId.String(),
		Description: c.Description,
		Key:         chacha_poly.GenerateKey(),
	}

	err = g.Doc.GroupAccessIndex.Add(userId, groupAccess)

	return
}
