package fetch

import (
	"context"

	"TenderServiceApi/internal/model"
)

type Service struct {
	bids                repository
	useCaseTenderVerify useCaseTenderVerify
	useCaseBidsVerify   useCaseBidsVerify
}

type repository interface {
	FetchListByTender(ctx context.Context, tenderId string) ([]model.Bids, error)
	FetchListByUser(ctx context.Context, username string) ([]model.Bids, error)
	FetchById(ctx context.Context, bidsId string) (model.Bids, error)
}

type useCaseTenderVerify interface {
	CheckResponsible(ctx context.Context, username string, tenderId string) (bool, error)
}

type useCaseBidsVerify interface {
	CheckResponsible(ctx context.Context, username string, bidsId string) (bool, error)
}

func NewService(r repository, useCaseTenderVerify useCaseTenderVerify, useCaseBidsVerify useCaseBidsVerify) *Service {
	return &Service{bids: r, useCaseTenderVerify: useCaseTenderVerify, useCaseBidsVerify: useCaseBidsVerify}
}

func (s *Service) FetchListByTender(ctx context.Context, username string, tenderId string) ([]model.Bids, error) {
	_, err := s.useCaseTenderVerify.CheckResponsible(ctx, username, tenderId)
	if err != nil {
		return []model.Bids{}, err
	}

	return s.bids.FetchListByTender(ctx, tenderId)
}

func (s *Service) FetchListByUser(ctx context.Context, username string) ([]model.Bids, error) {
	return s.bids.FetchListByUser(ctx, username)
}

func (s *Service) FetchStatus(ctx context.Context, username string, bidsId string) (model.Bids, error) {
	_, err := s.useCaseBidsVerify.CheckResponsible(ctx, username, bidsId)
	if err != nil {
		return model.Bids{}, err
	}

	bids, err := s.bids.FetchById(ctx, bidsId)
	if err != nil {
		return model.Bids{}, err
	}

	return bids, nil
}
