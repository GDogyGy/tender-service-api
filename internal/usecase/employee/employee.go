package employee

import (
	"context"

	"TenderServiceApi/internal/model"
)

type Service struct {
	employee Repository
}

type Repository interface {
	FetchById(ctx context.Context, id string) (model.Employee, error)
}

func NewService(r Repository) *Service {
	return &Service{employee: r}
}

func (s *Service) FetchById(ctx context.Context, id string) (model.Employee, error) {
	return s.employee.FetchById(ctx, id)
}
