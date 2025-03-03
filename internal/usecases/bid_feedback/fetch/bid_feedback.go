package fetch

import (
	"TenderServiceApi/internal/model"
	"context"
)

type Service struct {
	feedback                  repository
	useCaseTenderVerify       useCaseTenderVerify
	useCaseOrganizationVerify useCaseOrganizationVerify
}

type repository interface {
	FetchReviews(ctx context.Context, tenderID string, authorUsername string, organizationID string) ([]model.BidFeedback, error)
}

type useCaseTenderVerify interface {
	CheckResponsible(ctx context.Context, username string, tenderId string) (bool, error)
}

type useCaseOrganizationVerify interface {
	CheckResponsible(ctx context.Context, username string, organizationId string) (model.OrganizationResponsible, error)
}

func NewService(r repository, useCaseTenderVerify useCaseTenderVerify, useCaseOrganizationVerify useCaseOrganizationVerify) *Service {
	return &Service{feedback: r, useCaseTenderVerify: useCaseTenderVerify, useCaseOrganizationVerify: useCaseOrganizationVerify}
}

func (s *Service) FetchReviews(ctx context.Context, username string, tenderID string, authorUsername string, organizationID string) ([]model.BidFeedback, error) {
	_, err := s.useCaseTenderVerify.CheckResponsible(ctx, username, tenderID)
	if err != nil {
		return []model.BidFeedback{}, err
	}

	_, err = s.useCaseOrganizationVerify.CheckResponsible(ctx, username, organizationID)
	if err != nil {
		return []model.BidFeedback{}, err
	}

	bids, err := s.feedback.FetchReviews(ctx, tenderID, authorUsername, organizationID)
	if err != nil {
		return []model.BidFeedback{}, err
	}

	return bids, nil
}
