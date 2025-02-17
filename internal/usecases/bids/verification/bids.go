package fetch

import (
	"context"
)

type Service struct {
	bids repository
}

type repository interface {
	CheckResponsible(ctx context.Context, username string, bidsId string) (bool, error)
}

func NewService(r repository) *Service {
	return &Service{bids: r}
}

func (s *Service) CheckResponsible(ctx context.Context, username string, bidsId string) (bool, error) {
	return s.bids.CheckResponsible(ctx, username, bidsId)
}
