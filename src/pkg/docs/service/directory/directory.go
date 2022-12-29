package directory

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/processing"
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

func (s *Service) NewProcRequest(reqID, userID, ehrUUID string, kind processing.RequestKind) (processing.RequestInterface, error) {
	return s.Proc.NewRequest(reqID, userID, ehrUUID, kind)
}

// TODO
func (s *Service) Create(ctx context.Context, req processing.RequestInterface, systemID string, ehrUUID *uuid.UUID, user *userModel.UserInfo, d *model.Directory) error {
	return errors.ErrNotImplemented
}

// TODO
func (s *Service) Update(ctx context.Context, req processing.RequestInterface, systemID string, ehrUUID *uuid.UUID, user *userModel.UserInfo, d *model.Directory) error {
	if err := s.increaseVersion(d); err != nil {
		return fmt.Errorf("Directory increaseVersion error: %w directory.UID %s", err, d.UID.Value)
	}

	// TODO need realization
	//err = s.save(ctx, multiCallTx, procRequest, userID, systemID, ehrUUID, groupAccess, d)

	return errors.ErrNotImplemented
}

// TODO
func (s *Service) Delete(ctx context.Context, req processing.RequestInterface, systemID string, ehrUUID *uuid.UUID, versionUID, userID string) error {
	return errors.ErrNotImplemented
}

// TODO
func (s *Service) GetByTime(ctx context.Context, systemID string, ehrUUID *uuid.UUID, userID string, versionTime time.Time) (*model.Directory, error) {
	return nil, errors.ErrNotImplemented
}

// TODO
func (s *Service) GetByVersion(ctx context.Context, systemID string, ehrUUID *uuid.UUID, versionUID, userID string) (*model.Directory, error) {
	return nil, errors.ErrNotImplemented
}

// TODO
func (s *Service) GetByID(ctx context.Context, userID string, versionUID string) (*model.Directory, error) {
	return nil, errors.ErrNotImplemented
}

func (s *Service) increaseVersion(d *model.Directory) error {
	if _, err := d.IncreaseUIDVersion(); err != nil {
		return fmt.Errorf("Directory %s IncreaseUIDVersion error: %w", d.UID.Value, err)
	}

	return nil
}
