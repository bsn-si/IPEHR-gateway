package helper

import (
	"context"
)

type Finder interface {
	IsExist(ctx context.Context, userID string, systemID string, ID string) bool
}

type Searcher interface {
	IsExist(ID string) bool
}

type Search struct {
	service  Finder
	ctx      context.Context
	userID   string
	systemID string
}

func (h *Search) IsExist(ID string) bool {
	return h.service.IsExist(h.ctx, h.userID, h.systemID, ID)
}

func NewSearcher(ctx context.Context, userID string, systemID string, s Finder) *Search {
	return &Search{
		service:  s,
		ctx:      ctx,
		userID:   userID,
		systemID: systemID,
	}
}
