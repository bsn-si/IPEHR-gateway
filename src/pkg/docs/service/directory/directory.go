package directory

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/types"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/infrastructure"
	userModel "github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
)

type Service struct {
	Infra *infrastructure.Infra
	Proc  *processing.Proc
}

func NewService(infra *infrastructure.Infra, p *processing.Proc) *Service {
	return &Service{
		Infra: infra,
		Proc:  p,
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

func (s *Service) Delete(ctx context.Context, req processing.RequestInterface, systemID string, ehrUUID *uuid.UUID, versionUID, userID string) (string, error) {
	objectVersionID, err := base.NewObjectVersionID(versionUID, systemID)
	if err != nil {
		return "", fmt.Errorf("NewObjectVersionID error: %w versionUID %s ehrSystemID %s", err, versionUID, systemID)
	}

	_, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return "", fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	baseDocumentUID := []byte(objectVersionID.BasedID())
	baseDocumentUIDHash := sha3.Sum256(baseDocumentUID)

	txHash, err := s.Infra.Index.DeleteDoc(ctx, ehrUUID, types.Directory, &baseDocumentUIDHash, objectVersionID.VersionBytes(), userPrivKey, nil)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return "", err
		}
		return "", fmt.Errorf("Index.DeleteDoc error: %w", err)
	}

	req.AddEthereumTx(processing.TxDeleteDoc, txHash)

	if _, err = objectVersionID.IncreaseUIDVersion(); err != nil {
		return "", fmt.Errorf("IncreaseUIDVersion error: %w objectVersionID %s", err, objectVersionID.String())
	}

	return objectVersionID.String(), nil
}

// TODO
func (s *Service) GetByTime(ctx context.Context, systemID string, ehrUUID *uuid.UUID, userID string, versionTime time.Time) (*model.Directory, error) {
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
