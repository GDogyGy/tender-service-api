package fetch

import (
	"context"

	"TenderServiceApi/internal/model"
)

type Service struct {
	tender              repository
	useCaseTenderVerify useCaseTenderVerify
}

type repository interface {
	FetchList(ctx context.Context, serviceType string) ([]model.Tender, error)
	FetchListByUser(ctx context.Context, username string) ([]model.Tender, error)
	FetchById(ctx context.Context, tenderId string) (model.Tender, error)
}

type useCaseTenderVerify interface {
	CheckResponsible(ctx context.Context, username string, tenderId string) (bool, error)
}

func NewService(r repository, useCaseTenderVerify useCaseTenderVerify) *Service {
	return &Service{tender: r, useCaseTenderVerify: useCaseTenderVerify}
}

func (s *Service) FetchList(ctx context.Context, serviceType string) ([]model.Tender, error) {
	return s.tender.FetchList(ctx, serviceType)
}

func (s *Service) FetchListByUser(ctx context.Context, username string) ([]model.Tender, error) {
	return s.tender.FetchListByUser(ctx, username)
}

func (s *Service) FetchStatus(ctx context.Context, username string, tenderId string) (model.Tender, error) {
	_, err := s.useCaseTenderVerify.CheckResponsible(ctx, username, tenderId)
	if err != nil {
		return model.Tender{}, err
	}

	tender, err := s.tender.FetchById(ctx, tenderId)
	if err != nil {
		return model.Tender{}, err
	}

	return tender, nil
}
