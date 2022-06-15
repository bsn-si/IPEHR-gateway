package group_access

import (
	"github.com/google/uuid"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/indexer/service/group_access"
)

type GroupAccessService struct {
	GroupAccessIndex *group_access.GroupAccessIndex
	Cfg              *config.Config
}

func NewGroupAccessService(cfg *config.Config) *GroupAccessService {
	return &GroupAccessService{
		Cfg:              cfg,
		GroupAccessIndex: group_access.New(),
	}
}

func (g *GroupAccessService) Get(userId, groupId string) (groupAccess *model.GroupAccess, err error) {
	groupAccess, err = g.GroupAccessIndex.Get(userId, groupId)
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

	err = g.GroupAccessIndex.Add(userId, groupAccess)

	return
}
