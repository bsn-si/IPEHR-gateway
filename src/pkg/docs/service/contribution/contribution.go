package contribution

import (
	"context"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/helper"
	userModel "github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
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
					ok, err := hComposition.IsExist(v.PrecedingVersionUID.Value)
					if !ok || err != nil {
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

func (s *Service) PrepareResponse(ctx context.Context, systemID string, c *model.Contribution) (*model.ContributionResponse, error) {
	cR := model.ContributionResponse{
		UID: c.UID,
	}

	for _, v := range c.Versions {
		t := model.ContributionVersionResponse{
			Type:         v.Type,
			Contribution: v.Contribution,
		}

		cR.Versions = append(cR.Versions, t)
	}

	cR.Audit.TimeCommitted = base.DvDateTime{
		Value: time.Now().Format(common.OpenEhrTimeFormat),
	}

	cR.Audit.SystemID = systemID

	return &cR, nil
}

//func (*Service) rollback(ctx context.Context) error {
//	return errors.ErrNotImplemented
//}
