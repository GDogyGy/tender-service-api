package create

import (
	"context"

	"TenderServiceApi/internal/model"
)

type Service struct {
	tender                    repository
	useCaseOrganizationVerify useCaseOrganizationVerify
}

type repository interface {
	Create(ctx context.Context, saveModel model.Tender) (model.Tender, error)
}

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

	saveModel.Responsible = organizationResponsible.Id
	return s.tender.Create(ctx, saveModel)
}
