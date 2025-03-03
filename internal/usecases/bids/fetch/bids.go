package fetch

import (
	"context"

	"TenderServiceApi/internal/model"
)

type Service struct {
	bids                repository
	bidFeedback         bidFeedback
	useCaseTenderVerify useCaseTenderVerify
	useCaseBidsVerify   useCaseBidsVerify
}

type repository interface {
	FetchListByTender(ctx context.Context, tenderId string) ([]model.Bids, error)
	FetchListByUser(ctx context.Context, username string) ([]model.Bids, error)
	FetchById(ctx context.Context, bidsId string) (model.Bids, error)
}

type bidFeedback interface {
	FetchReviews(ctx context.Context, tenderID string, authorUsername string, organizationID string) ([]model.BidFeedback, error)
}

type useCaseTenderVerify interface {
	CheckResponsible(ctx context.Context, username string, tenderId string) (bool, error)
}

type useCaseBidsVerify interface {
	CheckResponsible(ctx context.Context, username string, bidsId string) (bool, error)
}

func NewService(r repository, b bidFeedback, useCaseTenderVerify useCaseTenderVerify, useCaseBidsVerify useCaseBidsVerify) *Service {
	return &Service{bids: r, bidFeedback: b, useCaseTenderVerify: useCaseTenderVerify, useCaseBidsVerify: useCaseBidsVerify}
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

func (s *Service) FetchReviews(ctx context.Context, username string, tenderID string, authorUsername string, organizationID string) ([]model.BidFeedback, error) {
	_, err := s.useCaseTenderVerify.CheckResponsible(ctx, username, tenderID)
	if err != nil {
		return []model.BidFeedback{}, err
	}

	bids, err := s.bidFeedback.FetchReviews(ctx, tenderID, authorUsername, organizationID)
	if err != nil {
		return []model.BidFeedback{}, err
	}

	return bids, nil
}
