package fetch

import (
	"context"

	"TenderServiceApi/internal/model"
)

type Service struct {
	organization repository
}

type repository interface {
	CheckResponsible(ctx context.Context, username string, organizationId string) (model.OrganizationResponsible, error)
}

func NewService(r repository) *Service {
	return &Service{organization: r}
}

func (s *Service) CheckResponsible(ctx context.Context, username string, organizationId string) (model.OrganizationResponsible, error) {
	return s.organization.CheckResponsible(ctx, username, organizationId)
}
