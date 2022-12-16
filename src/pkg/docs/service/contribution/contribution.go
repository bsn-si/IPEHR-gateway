package contribution

import (
	"context"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/helper"
	userModel "hms/gateway/pkg/user/model"
)

type Service struct {
	*service.DefaultDocumentService
}

func NewService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		docService,
	}
}

func (*Service) GetByID(ctx context.Context, userID string, cID string) (*model.ContributionResponse, error) {
	return nil, errors.ErrNotImplemented
}

func (*Service) Store(ctx context.Context, reqID, systemID string, user *userModel.UserInfo, c *model.Contribution) error {
	return errors.ErrNotImplemented
}

func (s *Service) Validate(_ context.Context, c *model.Contribution, template helper.Searcher) (bool, error) {
	return c.Validate(template)
}

// TODO need create loop, in which create actions should be in single PARENT transaction
// TODO commit only if lifecicle is complete or incomplete
// TODO when prev state was incomplete, and we active complete, we need change version +1
func (*Service) Commit(ctx context.Context, reqID, systemID string, user *userModel.UserInfo, c *model.Contribution) error {
	return errors.ErrNotImplemented
}

//func (h *Service) addCommiter(ctx context.Context, c *model.Contribution, u userModel.UserInfo) error {
//Composer: base.NewPartyProxy(
//	&base.PartyIdentified{
//		Name: "Silvia Blake",
//		PartyProxyBase: base.PartyProxyBase{
//			Type: base.PartyIdentifiedItemType,
//		},
//	},
//),

//"committer": {
//	"_type": "PARTY_IDENTIFIED",
//		"external_ref": {
//		"id": {
//			"_type": "HIER_OBJECT_ID",
//				"value": "f7e48c23-21b2-4b58-b9e0-a3ccece1bcf1"
//		},
//		"namespace": "DEMOGRAPHIC", // TODO ???
//			"type": "PERSON"
//	},
//	"name": "Dr. Yamamoto"
//}
//return nil
//}

//func (*Service) rollback(ctx context.Context) error {
//	return errors.ErrNotImplemented
//}
