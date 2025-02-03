package fetch

import (
	"context"

	"TenderServiceApi/internal/model"
)

type Service struct {
	organization repository
}

type repository interface {
	FetchById(ctx context.Context, id string) (model.Organization, error)
	FetchRelationsById(ctx context.Context, id string) (model.OrganizationResponsible, error)
}

func NewService(r repository) *Service {
	return &Service{organization: r}
}

func (s *Service) Edite(ctx context.Context, id string) (model.Organization, error) {
	return s.organization.FetchById(ctx, id)
}

func (s *Service) FetchRelationsById(ctx context.Context, id string) (model.OrganizationResponsible, error) {
	return s.organization.FetchRelationsById(ctx, id)
}
