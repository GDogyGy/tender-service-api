package tender

import (
	"context"

	"TenderServiceApi/internal/model"
)

type Service struct {
	tender repository
}

type repository interface {
	FetchList(ctx context.Context, serviceType string) ([]model.Tender, error)
}

func NewService(r repository) *Service {
	return &Service{tender: r}
}

func (s *Service) FetchList(ctx context.Context, serviceType string) ([]model.Tender, error) {
	return s.tender.FetchList(ctx, serviceType)
}
