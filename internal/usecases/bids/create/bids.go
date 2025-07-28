package create

import (
	"TenderServiceApi/internal/model"
	"context"
)

type Service struct {
	bids                      repository
	useCaseOrganizationVerify useCaseOrganizationVerify
}

//go:generate mockery --inpackage --name=repository --exported --testonly --inpackage-suffix
type repository interface {
	Create(ctx context.Context, saveModel model.Bids) (model.Bids, error)
}

//go:generate mockery --inpackage --name=useCaseOrganizationVerify --exported --testonly --inpackage-suffix
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
