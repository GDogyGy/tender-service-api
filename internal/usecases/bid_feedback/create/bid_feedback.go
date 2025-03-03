package create

import (
	"TenderServiceApi/internal/model"
	"context"
)

type Service struct {
	bidFeedback              repository
	organization             repositoryOrganization
	useCaseBidFeedbackVerify useCaseBidFeedbackVerify
}

type repository interface {
	Create(ctx context.Context, saveModel model.BidFeedback) (model.BidFeedback, error)
}

type repositoryOrganization interface {
	FetchByUserName(ctx context.Context, username string) (model.OrganizationResponsible, error)
}
type useCaseBidFeedbackVerify interface {
	CheckResponsible(ctx context.Context, username string, bidID string) (bool, error)
}

func NewService(r repository, o repositoryOrganization, useCaseBidFeedbackVerify useCaseBidFeedbackVerify) *Service {
	return &Service{bidFeedback: r, organization: o, useCaseBidFeedbackVerify: useCaseBidFeedbackVerify}
}

func (s *Service) Create(ctx context.Context, username string, bidID string, saveModel model.BidFeedback) (model.BidFeedback, error) {
	_, err := s.useCaseBidFeedbackVerify.CheckResponsible(ctx, username, bidID)
	if err != nil {
		return model.BidFeedback{}, err
	}
	responsible, err := s.organization.FetchByUserName(ctx, username)
	if err != nil {
		return model.BidFeedback{}, err
	}

	saveModel.BidID = bidID
	saveModel.Responsible = responsible.Id

	return s.bidFeedback.Create(ctx, saveModel)
}
