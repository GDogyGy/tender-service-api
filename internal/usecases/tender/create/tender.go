package create

import (
	"context"

	"TenderServiceApi/internal/model"
)

type Service struct {
	tender                    repository
	useCaseOrganizationVerify useCaseOrganizationVerify
}

//go:generate mockery --inpackage --name=repository --exported --testonly --inpackage-suffix
type repository interface {
	Create(ctx context.Context, saveModel model.Tender) (model.Tender, error)
}

//go:generate mockery --inpackage --name=useCaseOrganizationVerify --exported --testonly --inpackage-suffix
type useCaseOrganizationVerify interface {
	CheckResponsible(ctx context.Context, username string, organizationId string) (model.OrganizationResponsible, error)
}

func NewService(r repository, useCaseOrganizationVerify useCaseOrganizationVerify) *Service {
	return &Service{tender: r, useCaseOrganizationVerify: useCaseOrganizationVerify}
}

func (s *Service) Create(ctx context.Context, username string, organizationId string, saveModel model.Tender) (model.Tender, error) {
	organizationResponsible, err := s.useCaseOrganizationVerify.CheckResponsible(ctx, username, organizationId)
	if err != nil {
		return model.Tender{}, err
	}

	if saveModel.Version == 0 {
		saveModel.Version = 1
	}
	saveModel.Responsible = organizationResponsible.OrganizationId
	return s.tender.Create(ctx, saveModel)
}
