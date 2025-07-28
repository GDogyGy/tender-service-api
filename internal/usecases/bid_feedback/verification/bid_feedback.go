package fetch

import (
	"context"
)

type Service struct {
	bidFeedback repository
}

type repository interface {
	CheckResponsible(ctx context.Context, username string, bidID string) (bool, error)
}

func NewService(r repository) *Service {
	return &Service{bidFeedback: r}
}

func (s *Service) CheckResponsible(ctx context.Context, username string, bidID string) (bool, error) {
	return s.bidFeedback.CheckResponsible(ctx, username, bidID)
}
