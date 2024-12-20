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
	FetchListByUser(ctx context.Context, username string) ([]model.Tender, error)
	CheckResponsibleTender(ctx context.Context, username string, tenderId string) (bool, error)
	FetchTenderById(ctx context.Context, tenderId string) (model.Tender, error)
	CreateTender(ctx context.Context, saveModel model.Tender) (model.Tender, error)
}

func NewService(r repository) *Service {
	return &Service{tender: r}
}

func (s *Service) FetchList(ctx context.Context, serviceType string) ([]model.Tender, error) {
	return s.tender.FetchList(ctx, serviceType)
}

func (s *Service) FetchListByUser(ctx context.Context, username string) ([]model.Tender, error) {
	return s.tender.FetchListByUser(ctx, username)
}

func (s *Service) CheckResponsibleTender(ctx context.Context, username string, tenderId string) (bool, error) {
	return s.tender.CheckResponsibleTender(ctx, username, tenderId)
}

func (s *Service) FetchTenderById(ctx context.Context, tenderId string) (model.Tender, error) {
	return s.tender.FetchTenderById(ctx, tenderId)
}

func (s *Service) CreateTender(ctx context.Context, saveModel model.Tender) (model.Tender, error) {
	return s.tender.CreateTender(ctx, saveModel)
}
