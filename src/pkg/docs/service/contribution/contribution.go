package contribution

import (
	"context"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/processing"
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

func (s *Service) NewProcRequest(reqID, userID, ehrUUID string, kind processing.RequestKind) (processing.RequestInterface, error) {
	return s.Proc.NewRequest(reqID, userID, ehrUUID, kind)
}

func (*Service) GetByID(ctx context.Context, userID string, cID string) (*model.ContributionResponse, error) {
	return nil, errors.ErrNotImplemented
}

func (*Service) Store(ctx context.Context, req processing.RequestInterface, systemID string, user *userModel.UserInfo, c *model.Contribution) error {
	// TODO what about increase version?
	return errors.ErrNotImplemented
}

func (s *Service) Validate(_ context.Context, c *model.Contribution, template helper.Searcher) (bool, error) {
	return c.Validate(template)
}

// TODO need create loop, in which CREATE actions should be in single PARENT transaction
// TODO when prev state was incomplete, and we active complete, we need change version +1
func (s *Service) Execute(ctx context.Context, req processing.RequestInterface, userID, ehrUUID string, c *model.Contribution, hComposition helper.Searcher) error {
	listReqForCreation := make(map[string]bool, 0)

	for _, v := range c.Versions {
		switch v.Data.GetType() {
		case base.CompositionItemType:
			if !(v.LifecycleState.Value == "complete") || v.LifecycleState.Value != "incomplete" {
				// TODO we do not work with drafts, but "incomplete" - used in emergency situations, and may contain invalid data
				continue
			}

			switch v.CommitAudit.ChangeType.DefiningCode.CodeString {
			case "249":
				if v.UID.Value == "" {
					return errors.ErrFieldIsEmpty("uid")
				}
				// TODO create composition and put in dictionary list
				//if (err:=composition.Create(); err!=nil) {return err}
				listReqForCreation[v.UID.Value] = true

			case "251":
				if v.PrecedingVersionUID.Value == "" {
					return errors.ErrFieldIsEmpty("preceding_version_uid")
				}

				if !listReqForCreation[v.PrecedingVersionUID.Value] {
					if !hComposition.IsExist(v.PrecedingVersionUID.Value) {
						return errors.New("Can not modify composition because its not created")
					}
				}
				// TODO put modification in composition

			case "523":
				if v.PrecedingVersionUID.Value == "" {
					return errors.ErrFieldIsEmpty("preceding_version_uid")
				}
				// TODO delete composition
				continue
			}

		case base.ContributionItemType:
			c := v.Data.(model.Contribution)

			err := s.Execute(ctx, req, userID, ehrUUID, &c, hComposition)
			if err != nil {
				return errors.Wrap(err, "CONTRIBUTION commit error")
			}
		}
	}
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
