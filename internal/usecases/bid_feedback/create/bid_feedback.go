package create

import (
	"TenderServiceApi/internal/model"
	"context"
)

type Service struct {
	bidFeedback              repositoryFeedback
	organization             repositoryOrganization
	useCaseBidFeedbackVerify useCaseBidFeedbackVerify
}

//go:generate mockery --inpackage --name=repositoryFeedback --exported --testonly --inpackage-suffix
type repositoryFeedback interface {
	Create(ctx context.Context, saveModel model.BidFeedback) (model.BidFeedback, error)
}

//go:generate mockery --inpackage --name=repositoryOrganization --exported --testonly --inpackage-suffix
type repositoryOrganization interface {
	FetchByUserName(ctx context.Context, username string) (model.OrganizationResponsible, error)
}

//go:generate mockery --inpackage --name=useCaseBidFeedbackVerify --exported --testonly --inpackage-suffix
type useCaseBidFeedbackVerify interface {
	CheckResponsible(ctx context.Context, username string, bidID string) (bool, error)
}

func NewService(r repositoryFeedback, o repositoryOrganization, useCaseBidFeedbackVerify useCaseBidFeedbackVerify) *Service {
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
