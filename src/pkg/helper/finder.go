package helper

import (
	"context"

	"github.com/google/uuid"
)

type (
	Finder interface {
		IsExist(ctx context.Context, args ...string) (bool, error)
		GetEhrUUIDByUserID(ctx context.Context, userID, systemID string) (*uuid.UUID, error)
	}

	Searcher interface {
		UseService(s Finder) *Search
		IsExist(ID string) (bool, error)
		GetEhrUUIDByUserID() (*uuid.UUID, error)
		IsEhrBelongsToUser() bool
	}

	Search struct {
		service  Finder
		ctx      context.Context
		userID   string
		ehrUUID  string
		systemID string
	}
)

func (h *Search) IsExist(ID string) (bool, error) {
	return h.service.IsExist(h.ctx, h.userID, h.systemID, ID)
}

func (h *Search) GetEhrUUIDByUserID() (*uuid.UUID, error) {
	return h.service.GetEhrUUIDByUserID(h.ctx, h.userID, h.systemID)
}

func (h *Search) IsEhrBelongsToUser() bool {
	ehrUUID, err := h.service.GetEhrUUIDByUserID(h.ctx, h.userID, h.systemID)

	if err != nil {
		return false
	}

	return h.ehrUUID == ehrUUID.String()
}

func (h *Search) UseService(s Finder) *Search {
	h.service = s
	return h
}

func NewSearcher(ctx context.Context, userID, systemID, ehrUUID string) *Search {
	return &Search{
		ctx:      ctx,
		userID:   userID,
		ehrUUID:  ehrUUID,
		systemID: systemID,
	}
}
