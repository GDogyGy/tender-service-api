package create

import (
	"TenderServiceApi/internal/model"
	"context"
)

type Service struct {
	bids                      repository
	useCaseOrganizationVerify useCaseOrganizationVerify
}

type repository interface {
	Create(ctx context.Context, saveModel model.Bids) (model.Bids, error)
}

type useCaseOrganizationVerify interface {
	CheckResponsible(ctx context.Context, username string, organizationId string) (model.OrganizationResponsible, error)
}

func NewService(r repository, useCaseOrganizationVerify useCaseOrganizationVerify) *Service {
	return &Service{bids: r, useCaseOrganizationVerify: useCaseOrganizationVerify}
}

func (s *Service) Create(ctx context.Context, username string, organizationId string, saveModel model.Bids) (model.Bids, error) {
	organizationResponsible, err := s.useCaseOrganizationVerify.CheckResponsible(ctx, username, organizationId)
	if err != nil {
		return model.Bids{}, err
	}
	saveModel.Responsible = organizationResponsible.Id

	return s.bids.Create(ctx, saveModel)
}
