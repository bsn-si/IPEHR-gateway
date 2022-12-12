package contribution

import (
	"context"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/errors"
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

func (*Service) GetByID(ctx context.Context, userID string, cID string) (*model.Contribution, error) {
	return nil, errors.ErrNotImplemented
}

func (*Service) Store(ctx context.Context) error {
	return errors.ErrNotImplemented
}

func (*Service) Validate(ctx context.Context, data interface{}) (bool, error) {
	return false, errors.ErrNotImplemented
}

func (*Service) Execute(ctx context.Context) error {
	// TODO need create loop, in which create actions should be in single PARENT transaction
	return errors.ErrNotImplemented
}

func (h *Service) AddCommiter(ctx context.Context, c *model.Contribution, u userModel.UserInfo) error {
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
}

//func (*Service) Rollback(ctx context.Context) error {
//	return errors.ErrNotImplemented
//}
