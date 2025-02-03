package fetch

import (
	"context"
)

type Service struct {
	tender repository
}

type repository interface {
	CheckResponsible(ctx context.Context, username string, tenderId string) (bool, error)
}

func NewService(r repository) *Service {
	return &Service{tender: r}
}

func (s *Service) CheckResponsible(ctx context.Context, username string, tenderId string) (bool, error) {
	return s.tender.CheckResponsible(ctx, username, tenderId)
}
