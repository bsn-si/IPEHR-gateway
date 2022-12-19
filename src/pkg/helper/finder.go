package helper

import (
	"context"
)

type Finder interface {
	IsExist(ctx context.Context, userID, systemID, ehrUUID, ID string) bool
}

type Searcher interface {
	IsExist(ID string) bool
}

type Search struct {
	service  Finder
	ctx      context.Context
	userID   string
	ehrUUID  string
	systemID string
}

func (h *Search) IsExist(ID string) bool {
	return h.service.IsExist(h.ctx, h.userID, h.systemID, h.ehrUUID, ID)
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
